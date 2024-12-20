package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"google.golang.org/api/slides/v1"
)

func main() {
	presentationID := flag.String("id", "", "ID of the slide to update, empty means create a new one")
	flag.Parse()
	ctx := context.Background()

	// Initialize Google services
	cacheDir, err := gcputils.GetCredentialCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	client := gcputils.InitGoogleClient("../../../credentials.json", filepath.Join(cacheDir, "slides.googleapis.com-go-quickstart.json"))
	srv := gcputils.InitSlidesService(ctx, client)
	p := &slidesutils.Presentation{
		PresentationID:          *presentationID,
		PresentationService:     slides.NewPresentationsService(srv),
		PresentationPageService: slides.NewPresentationsPagesService(srv),
	}
	p.UpdatePresentationPointer()

	slideNumber := 0
	for slide := range p.Slides() {
		content := make([]string, 3)
		// BUG for an unknown reason (yet) sometimes the slides seems to be garbage collected
		if slide == nil {
			continue
		}
		for _, element := range slide.PageElements {
			if element.Shape != nil && element.Shape.Text != nil && element.Shape.Placeholder != nil {
				textContent := extractTextFromShape(element.Shape.Text)
				switch element.Shape.Placeholder.Type {
				case "TITLE":
					content[0] = textContent
				case "SUBTITLE":
					content[1] = textContent
				case "BODY":
					content[2] = textContent
				default:
					content = append(content, textContent)
				}
			}
		}
		if content[0] == "" && content[1] == "" && len(content[2]) < 50 {
			fmt.Printf("# %v\n", content[2])
		}
		if content[0] == "" {
			content[0] = "Slide " + strconv.Itoa(slideNumber)
		}
		fmt.Printf("## %v\n", content[0])

		if content[1] != "" {
			fmt.Printf("### %v\n", content[0])
		}
		for i := 2; i < len(content); i++ {
			fmt.Printf("%v\n", content[i])
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
