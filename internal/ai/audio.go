package ai

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go"
	"github.com/owulveryck/gptslideshow/config"
)

// ExtractTextFromAudio extracts text from an audio file using OpenAI's Whisper model.
// It takes an audio file path as input and returns the transcribed text.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - filePath: The path to the audio file to be transcribed.
//
// Returns:
//   - A string containing the transcribed text.
//   - An error if the transcription fails or if there is an issue with file handling.
func (ai *AI) ExtractTextFromAudio(ctx context.Context, filePath string) (string, error) {
	// Open the audio file for reading.
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open audio file: %w", err)
	}
	defer file.Close()

	// Request transcription from OpenAI's API using the Whisper model.
	transcription, err := ai.Client.Audio.Transcriptions.New(ctx, openai.AudioTranscriptionNewParams{
		Model:    openai.F(openai.AudioModelWhisper1),
		File:     openai.F[io.Reader](file),
		Language: openai.F(config.ConfigInstance.AudioLanguage),
	})
	if err != nil {
		return "", fmt.Errorf("failed to transcribe audio: %w", err)
	}

	// Return the transcribed text.
	return transcription.Text, nil
}
