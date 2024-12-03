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
	prompt := flag.String("prompt", "Convert the following Markdown text into an array of structured slides. Each slide should have a title, a subtitle, and a body:", "The prompt")
	textfile := flag.String("content", "", "The content file")
	audiofile := flag.String("audio", "", "The audio file in mp3")

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

	var content []byte
	// Read content file
	if *textfile != "" {
		content, err = os.ReadFile(*textfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *audiofile != "" {
		cont, err := getAudioFromFile(context.Background(), *audiofile)
		if err != nil {
			log.Fatal(err)
		}
		content = []byte(cont)
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
	presentationData, err := GenerateSlides(ctx, *prompt, content)
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
