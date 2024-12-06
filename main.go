package main

import (
	"context"
	"log"

	"github.com/owulveryck/gptslideshow/config"
	"github.com/owulveryck/gptslideshow/internal/ai"
	drive "google.golang.org/api/drive/v3"
)

func main() {
	// Parse command-line flags
	presentationId, fromTemplate, prompt, textfile, audiofile, helpFlag := parseFlags()

	if *helpFlag {
		printHelp()
		return
	}

	ctx := context.Background()
	openaiClient := ai.NewAI()

	// Initialize Google services
	client := initGoogleClient()
	slidesSrv := initSlidesService(client)
	var driveSrv *drive.Service

	// Read content from file or audio
	content := readContent(ctx, openaiClient, textfile, audiofile)

	// Handle template copy if specified
	if *fromTemplate != "" {
		driveSrv = initDriveService(client)
		p := handleTemplateCopy(ctx, driveSrv, *fromTemplate)
		presentationId = &p
	}
	if config.ConfigInstance.WithImage && driveSrv == nil {
		driveSrv = initDriveService(client)
	}

	// Generate slides from content
	presentationData := generateSlides(ctx, openaiClient, *prompt, content)

	// Create presentation slides
	err := createPresentationSlides(ctx, slidesSrv, driveSrv, openaiClient, config.ConfigInstance.WithImage, *presentationId, presentationData)
	if err != nil {
		log.Fatal(err)
	}
}
