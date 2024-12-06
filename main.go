package main

import (
	"context"

	"github.com/owulveryck/gptslideshow/internal/ai"
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

	// Read content from file or audio
	content := readContent(ctx, openaiClient, textfile, audiofile)

	// Handle template copy if specified
	if *fromTemplate != "" {
		driveSrv := initDriveService(client)
		p := handleTemplateCopy(ctx, driveSrv, *fromTemplate)
		presentationId = &p
	}

	// Generate slides from content
	presentationData := generateSlides(ctx, openaiClient, *prompt, content)

	// Create presentation slides
	createPresentationSlides(ctx, slidesSrv, *presentationId, presentationData)
}
