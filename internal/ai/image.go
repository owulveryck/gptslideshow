package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/png"

	"github.com/openai/openai-go"
)

// GenerateImageFromText generates a 256x256 image based on the provided text prompt.
// It communicates with OpenAI's API to create the image and returns it as an image.Image object.
//
// Parameters:
//   - ctx: The context for managing request deadlines and cancellation signals.
//   - prompt: A string containing the description of the image to generate.
//
// Returns:
//   - An image.Image object representing the generated image.
//   - An error if the image generation or processing fails.
func (ai *AI) GenerateImageFromText(ctx context.Context, prompt string) (image.Image, error) {
	// Request image generation from OpenAI's API with specified parameters.
	response, err := ai.Client.Images.Generate(ctx, openai.ImageGenerateParams{
		Prompt:         openai.String("generate an illustration based on those elements, the illustration should not contain any text: \n\n" + prompt),
		Model:          openai.F(openai.ImageModelDallE3),
		ResponseFormat: openai.F(openai.ImageGenerateParamsResponseFormatB64JSON),
		N:              openai.Int(1),
		Size:           openai.F(openai.ImageGenerateParamsSize1024x1024),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// Decode the base64-encoded image data from the API response.
	imageData, err := base64.StdEncoding.DecodeString(response.Data[0].B64JSON)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image data: %w", err)
	}

	// Decode the raw bytes into a PNG image.
	img, err := png.Decode(bytes.NewReader(imageData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG image: %w", err)
	}

	// Return the decoded image.
	return img, nil
}
