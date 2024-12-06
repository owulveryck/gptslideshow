package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/owulveryck/gptslideshow/internal/gcputils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils"
	"github.com/owulveryck/gptslideshow/internal/structure"
	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	slides "google.golang.org/api/slides/v1"
)

func main() {
	var presentationId string

	flag.StringVar(&presentationId, "id", "", "ID of the slide to update, empty means create a new one")
	flag.Parse()
	ctx := context.Background()

	// Load client secret file
	b, err := os.ReadFile("../../../credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Initialize Google OAuth2 config
	config, err := google.ConfigFromJSON(b, drive.DriveScope, slides.PresentationsScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Get authenticated client
	client := gcputils.GetClient(config)

	// Initialize Google Slides service
	slidesSrv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	// Initialize the Builder
	builder, err := slidesutils.NewBuilder(ctx, slidesSrv, presentationId)
	if err != nil {
		log.Fatalf("Unable to create Builder: %v", err)
	}

	log.Println("Generating slides...")
	slide := structure.Slide{
		Title:    "Title of the slide",
		Subtitle: "Subtitle of the slide",
		Body:     "The content of an awesome slide",
		Chapter:  false,
	}

	// Create a chapter slide
	err = builder.CreateChapter(ctx, slide)
	if err != nil {
		log.Fatal(err)
	}

	// Create a slide with title, subtitle, and body
	err = builder.CreateSlideTitleSubtitleBody(ctx, slide)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New presentation created and modified successfully.")
}
