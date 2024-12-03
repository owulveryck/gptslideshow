package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	slides "google.golang.org/api/slides/v1"
)

func main() {
	var presentationId string
	fromTemplate := flag.String("t", "", "ID of a template file")
	flag.StringVar(&presentationId, "id", "", "ID of the slide to update, empty means create a new one")
	filename := flag.String("content", "./testdata/article.md", "The content file")

	flag.Parse()
	ctx := context.Background()

	// Load client secret file
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Initialize Google OAuth2 config
	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Get authenticated client
	client := GetClient(config)

	// Initialize Google Slides service
	slidesSrv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	// Read content file
	content, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}

	// Handle template copy if specified
	if *fromTemplate != "" {
		driveSrv, err := drive.New(client)
		if err != nil {
			log.Fatalf("Unable to retrieve Drive client: %v", err)
		}

		presentationId, err = CopyTemplate(ctx, driveSrv, *fromTemplate)
		if err != nil {
			log.Fatalf("Unable to copy presentation: %v", err)
		}
	}

	// Generate slides from content
	presentationData, err := GenerateSlides(ctx, content)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Generating slides...")
	for i, slide := range presentationData.Slides {
		log.Printf("Slide %v: %v", i, slide.Title)
		err = CreateSlide(ctx, slidesSrv, presentationId, slide)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("New presentation created and modified successfully.")
}
