package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/gptslideshow/internal/ai/openai"
	vertex "github.com/owulveryck/gptslideshow/internal/ai/vertexai"
	"github.com/owulveryck/gptslideshow/internal/driveutils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"google.golang.org/api/slides/v1"
)

type configuration struct {
	GCPPRoject  string `envconfig:"GCP_PROJECT" required:"true"`
	GeminiModel string `envconfig:"GEMINI_MODEL" default:"gemini-1.5-pro"`
	GCPRegion   string `envconfig:"GCP_REGION" default:"us-central1"`
}

func main() {
	//	go startServer()
	// Parse command-line flags
	presentationID := flag.String("id", "", "ID of the slide to update, empty means create a new one")
	var config configuration
	err := envconfig.Process("", &config)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()
	ctx := context.Background()
	openaiClient := openai.NewAI()
	//	ollamaClient := ollama.NewAI()
	vertexAIClient := vertex.NewAI(ctx, config.GCPPRoject, config.GCPRegion, config.GeminiModel)

	// Initialize Google services
	client := initGoogleClient("../credentials.json")
	clientPerso := initGoogleClientPerso("../credentials_perso.json")
	srv := initSlidesService(ctx, client)
	driveSrv := initDriveService(clientPerso)

	fmt.Println("Scanning document")
	h := &helper{
		presentationID:          *presentationID,
		presentationService:     slides.NewPresentationsService(srv),
		presentationPageService: slides.NewPresentationsPagesService(srv),
	}
	for {
		must(h.updatePresentationPointer())
		if !h.presentationHasChanged() {
			time.Sleep(1 * time.Second)
			continue
		}
		for slide := range h.Slides() {
			// BUG for an unknown reason (yet) sometimes the slides seems to be garbage collected
			if slide == nil {
				continue
			}
			for _, element := range slide.PageElements {
				if element.Shape != nil && element.Shape.Text != nil {
					// Extract text content from shapes
					textContent := extractTextFromShape(element.Shape.Text)
					switch {
					case strings.Contains(textContent, "@dalle"):
						textContent = strings.ReplaceAll(textContent, "@dalle", "")
						img, err := openaiClient.GenerateImageURLFromText(ctx, textContent)
						if err != nil {
							log.Println(err)
						}
						if img == "" {
							log.Println("No image generated")
							continue
						}
						requests := processText(element.ObjectId, textContent)
						b, _ := element.Size.MarshalJSON()
						fmt.Printf("%s\n", b)
						b, _ = element.Transform.MarshalJSON()
						fmt.Printf("%s\n", b)
						imgRequest := slides.CreateImageRequest{
							ElementProperties: &slides.PageElementProperties{
								PageObjectId: slide.ObjectId,
								Size:         element.Size,
								Transform:    element.Transform,
							},
							Url:             img,
							ForceSendFields: []string{},
							NullFields:      []string{},
						}
						_, err = h.presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: append(requests, &slides.Request{CreateImage: &imgRequest}),
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}
					case strings.Contains(textContent, "@image"):
						textContent = strings.ReplaceAll(textContent, "@image", "")
						img, err := vertexAIClient.GenerateImageFromText(ctx, textContent)
						if err != nil {
							log.Println(err)
							continue
						}
						name := uuid.New().String() + ".png"
						imageName, err := driveutils.UploadImage(ctx, driveSrv, img, name)
						if err != nil {
							log.Println(err)
							continue
						}
						requests := processText(element.ObjectId, "")
						// Execute the batch update

						imgRequest := slides.CreateImageRequest{
							ElementProperties: &slides.PageElementProperties{
								PageObjectId: slide.ObjectId,
								Size:         element.Size,
								Transform:    element.Transform,
							},
							Url:             imageName,
							ForceSendFields: []string{},
							NullFields:      []string{},
						}
						_, err = h.presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: append(requests, &slides.Request{CreateImage: &imgRequest}),
						}).Do()
						if err != nil {
							log.Printf("unable to update text: %v", err)
						}

					case strings.Contains(textContent, "@gemini"):
						err := h.updateFromAI(ctx, vertexAIClient, element.ObjectId, textContent, "@gemini")
						if err != nil {
							log.Println(err)
						}
					case strings.Contains(textContent, "@format"):
						textContent = strings.ReplaceAll(textContent, "@format", "")
						requests := processText(element.ObjectId, textContent)
						// Execute the batch update
						_, err = h.presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: requests,
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}
					case strings.Contains(textContent, "@chatgpt"):
						err := h.updateFromAI(ctx, openaiClient, element.ObjectId, textContent, "@chatgpt")
						if err != nil {
							log.Println(err)
						}
						/*
							case strings.Contains(textContent, "@ollama"):
								err := h.updateFromAI(ctx, ollamaClient, element.ObjectId, textContent, "@ollama")
								if err != nil {
									log.Println(err)
								}
						*/
					}
				}
			}
		}
	}
}

// extractTextFromShape extracts text from a shape's text content.
func extractTextFromShape(text *slides.TextContent) string {
	var result string
	for _, textElement := range text.TextElements {
		if textElement.TextRun != nil {
			result += textElement.TextRun.Content
		}
	}
	return result
}

func processText(objectID, input string) []*slides.Request {
	content := []*slides.Request{
		{
			DeleteText: &slides.DeleteTextRequest{
				ObjectId: objectID,
				TextRange: &slides.Range{
					Type: "ALL",
				},
			},
		},
	}

	content = append(content, slidesutils.InsertMarkdownContent(input, objectID)...)
	slidesutils.SortRequests(content)
	// Create a batch update request to replace the text
	return content
}
