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

	// Create a slide with title, subtitle, and body
	err = builder.CreateSlideTitleSubtitleBody(ctx, slide)
	if err != nil {
		log.Fatal(err)
	}
	// Create a chapter slide
	err = builder.CreateChapter(ctx, slide)
	if err != nil {
		log.Fatal(err)
	}
	// Define the image properties
	imageUrl := "https://upload.wikimedia.org/wikipedia/commons/8/8d/Sinclair_QL_256x256_mode_example_image.png"

	var imageWidth, imageHeight, slideWidth, slideHeight, translateX, translateY float64
	// Dimensions in EMUs (e.g., 3x3 inches)
	imageWidth = 2743200  // 3 inches
	imageHeight = 2743200 // 3 inches

	// Slide dimensions in EMUs
	slideWidth = 9144000  // 10 inches
	slideHeight = 6858000 // 7.5 inches

	// Calculate position to center the image
	translateX = (slideWidth - imageWidth) / 2
	translateY = (slideHeight - imageHeight) / 2
	translateX = 1213950
	translateY = 1659800

	// Insert the image centered on the current slide
	err = builder.InsertImage(ctx, imageUrl, imageWidth, imageHeight, translateX, translateY)
	if err != nil {
		log.Fatalf("Error inserting image: %v", err)
	}

	log.Println("Image inserted successfully.")

	fmt.Println("New presentation created and modified successfully.")
}
