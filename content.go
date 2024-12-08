package main

import (
	"context"
	"log"
	"os"

	"github.com/owulveryck/gptslideshow/internal/ai"
)

func readContent(ctx context.Context, openaiClient *ai.AI, textfile, audiofile *string) []byte {
	var content []byte
	var err error

	if *textfile != "" {
		content, err = os.ReadFile(*textfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *audiofile != "" {
		b, err := openaiClient.ExtractTextFromAudio(ctx, *audiofile)
		if err != nil {
			log.Fatal(err)
		}
		content = []byte(b)
		// Create a temporary directory
		tempDir, err := os.MkdirTemp("", "gptslideshow")
		if err != nil {
			log.Fatalf("Failed to create temp directory: %v", err)
		}

		// Create a temporary file within the directory
		tempFile, err := os.CreateTemp(tempDir, "transcript-*.txt")
		if err != nil {
			log.Fatalf("Failed to create temp file: %v", err)
		}
		defer tempFile.Close()

		// Write the string to the file
		_, err = tempFile.WriteString(b)
		if err != nil {
			log.Fatalf("Failed to write to temp file: %v", err)
		}

		// Output the temporary file path
		log.Printf("Temporary file created: %s", tempFile.Name())
	}

	return content
}
