package main

import (
	"context"
	"log"
	"os"

	"github.com/owulveryck/gptslideshow/internal/ai"
)

func readContent(ctx context.Context, aiClient ai.AIInterface, textfile, audiofile *string) []byte {
	var content []byte
	var err error

	if *textfile != "" {
		content, err = os.ReadFile(*textfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *audiofile != "" {
		b, err := aiClient.ExtractTextFromAudio(ctx, *audiofile)
		if err != nil {
			log.Fatal(err)
		}
		content = []byte(b)
	}
	saveContent("input-*.txt", content)
	return content
}
