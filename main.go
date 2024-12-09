package main

import (
	"context"
	"log"

	drive "google.golang.org/api/drive/v3"

	"github.com/owulveryck/gptslideshow/config"
	"github.com/owulveryck/gptslideshow/internal/ai"
	"github.com/owulveryck/gptslideshow/internal/driveutils"
	"github.com/owulveryck/gptslideshow/internal/slidesutils/mytemplate"
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
	driveSrv = initDriveService(client)
	if *fromTemplate != "" {
		p := handleTemplateCopy(ctx, driveSrv, *fromTemplate)
		presentationId = &p
	}

	// Generate slides from content
	presentationData := generateSlides(ctx, openaiClient, *prompt, content)
	// Using mytemplate change to use yours
	builder, err := mytemplate.NewBuilder(ctx, slidesSrv, *presentationId)

	// Create presentation slides
	err = createPresentationSlides(ctx, builder, driveSrv, openaiClient, config.ConfigInstance.WithImage, presentationData)
	if err != nil {
		log.Fatal(err)
	}
	b, err := driveutils.ExtractPDF(ctx, driveSrv, *presentationId)
	if err != nil {
		log.Fatal(err)
	}
	err = saveContent("output-*.pdf", b)
	if err != nil {
		log.Fatal(err)
	}
}
