package main

import (
	"context"
)

func main() {
	// Parse command-line flags
	presentationId, fromTemplate, prompt, textfile, audiofile, helpFlag := parseFlags()

	if *helpFlag {
		printHelp()
		return
	}

	ctx := context.Background()

	// Initialize Google services
	client := initGoogleClient()
	slidesSrv := initSlidesService(client)

	// Read content from file or audio
	content := readContent(ctx, textfile, audiofile)

	// Handle template copy if specified
	if *fromTemplate != "" {
		driveSrv := initDriveService(client)
		p := handleTemplateCopy(ctx, driveSrv, *fromTemplate)
		presentationId = &p
	}

	// Generate slides from content
	presentationData := generateSlides(ctx, *prompt, content)

	// Create presentation slides
	createPresentationSlides(ctx, slidesSrv, *presentationId, presentationData)
}
