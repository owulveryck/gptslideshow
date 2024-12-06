package main

import (
	"context"
	"log"
	"os"
)

func readContent(ctx context.Context, textfile, audiofile *string) []byte {
	var content []byte
	var err error

	if *textfile != "" {
		content, err = os.ReadFile(*textfile)
		if err != nil {
			log.Fatal(err)
		}
	}

	if *audiofile != "" {
		b, err := getAudioFromFile(ctx, *audiofile)
		if err != nil {
			log.Fatal(err)
		}
		content = []byte(b)
	}

	return content
}
