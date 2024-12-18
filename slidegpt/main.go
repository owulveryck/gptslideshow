package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/owulveryck/gptslideshow/internal/ai/ollama"
	"github.com/owulveryck/gptslideshow/internal/ai/openai"
	vertex "github.com/owulveryck/gptslideshow/internal/ai/vertexai"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"google.golang.org/api/slides/v1"
)

type configuration struct {
	GCPPRoject  string `envconfig:"GCP_PROJECT" required:"true"`
	GeminiModel string `envconfig:"GEMINI_MODEL" default:"gemini-1.5-pro"`
	GCPRegion   string `envconfig:"GCP_REGION" default:"us-central1"`
}

func main() {
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
	ollamaClient := ollama.NewAI()
	vertexAIClient := vertex.NewAI(ctx, config.GCPPRoject, config.GCPRegion, config.GeminiModel)

	// Initialize Google services
	client := initGoogleClient()
	srv := initSlidesService(ctx, client)
	presentationService := slides.NewPresentationsService(srv)
	presentationPageService := slides.NewPresentationsPagesService(srv)
	// Retrieve the presentation
	// Iterate through each slide (page) in the presentation

	fmt.Println("Scanning document")
	h := &helper{
		presentationID:          *presentationID,
		presentationService:     presentationService,
		presentationPageService: presentationPageService,
	}
	for {
		must(h.updatePresentationPointer())
		if !h.presentationHasChanged() {
			time.Sleep(1 * time.Second)
			continue
		}
		// Get the total content
		for _, slide := range h.presentation.Slides {
			// printProgress(i, len(presentation.Slides))
			slide, err := presentationPageService.Get(*presentationID, slide.ObjectId).Do()
			if err != nil {
				log.Fatalf("Unable to retrieve current slide: %v", err)
			}
			// Iterate through all page elements (shapes, images, etc.)
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
						_, err = presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: append(requests, &slides.Request{CreateImage: &imgRequest}),
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}
					case strings.Contains(textContent, "@gemini"):
						textContent = strings.ReplaceAll(textContent, "@gemini", "")
						result, err := vertexAIClient.SimpleQuery(ctx, textContent)
						if err != nil {
							log.Fatal(err)
						}
						requests := processText(element.ObjectId, result)
						// Execute the batch update
						_, err = presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: requests,
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}

					case strings.Contains(textContent, "@format"):
						textContent = strings.ReplaceAll(textContent, "@format", "")
						requests := processText(element.ObjectId, textContent)
						// Execute the batch update
						_, err = presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: requests,
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}
					case strings.Contains(textContent, "@chatgpt"):
						textContent = strings.ReplaceAll(textContent, "@chatgpt", "")
						if strings.Contains(textContent, "@withContext") {
							textContent = strings.ReplaceAll(textContent, "@withContent", "")
						}
						result, err := openaiClient.SimpleQuery(ctx, textContent)
						if err != nil {
							log.Fatal(err)
						}
						requests := processText(element.ObjectId, result)
						// Execute the batch update
						_, err = presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: requests,
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}
					case strings.Contains(textContent, "@ollama"):
						textContent = strings.ReplaceAll(textContent, "@ollama", "")
						result, err := ollamaClient.SimpleQuery(ctx, textContent)
						if err != nil {
							log.Fatal(err)
						}
						requests := processText(element.ObjectId, result)
						// Execute the batch update
						_, err = presentationService.BatchUpdate(*presentationID, &slides.BatchUpdatePresentationRequest{
							Requests: requests,
						}).Do()
						if err != nil {
							log.Fatalf("unable to update text: %v", err)
						}
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
