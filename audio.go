package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	openai "github.com/rakyll/openai-go"
	"github.com/rakyll/openai-go/audio"
)

func getAudioFromFile(ctx context.Context, audiofile string) (string, error) {
	s := openai.NewSession(os.Getenv("OPENAI_API_KEY"))
	s.HTTPClient = &http.Client{
		Timeout: 300 * time.Second, // Set the timeout to 30 seconds
	}
	client := audio.NewClient(s, "")
	f, err := os.Open(audiofile)
	if err != nil {
		log.Fatalf("error opening audio file: %v", err)
	}
	defer f.Close()
	resp, err := client.CreateTranscription(ctx, &audio.CreateTranscriptionParams{
		Language:    "en",
		Audio:       f,
		AudioFormat: "mp3",
	})
	if err != nil {
		log.Fatalf("error transcribing file: %v", err)
	}
	return resp.Text, nil
}
