package main

import (
	"flag"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	drive "google.golang.org/api/drive/v3"
	slides "google.golang.org/api/slides/v1"
)

func main() {
	var presentationId string
	flag.StringVar(&presentationId, "id", "", "ID of the slide to update, empty means create a new one")

	flag.Parse()

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
	client := GetClient(config)

	// Initialize Google Slides service
	slidesSrv, err := slides.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Slides client: %v", err)
	}

	err = ListSlidesAndElements(slidesSrv, presentationId)
	if err != nil {
		log.Fatal(err)
	}
}
