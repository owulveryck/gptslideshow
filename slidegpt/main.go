package main

import (
	"context"
	"flag"
	"log"
	"strings"

	"github.com/owulveryck/gptslideshow/internal/ai/ollama"
	"github.com/owulveryck/gptslideshow/internal/ai/openai"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"google.golang.org/api/slides/v1"
)

func main() {
	// Parse command-line flags
	presentationID := flag.String("id", "", "ID of the slide to update, empty means create a new one")

	flag.Parse()
	ctx := context.Background()
	openaiClient := openai.NewAI()
	ollamaClient := ollama.NewAI()

	// Initialize Google services
	client := initGoogleClient()
	srv := initSlidesService(ctx, client)
	presentationService := slides.NewPresentationsService(srv)
	presentationPageService := slides.NewPresentationsPagesService(srv)
	// Retrieve the presentation
	presentation, err := presentationService.Get(*presentationID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve presentation: %v", err)
	}
	// Iterate through each slide (page) in the presentation
	for {
		for _, slide := range presentation.Slides {
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
