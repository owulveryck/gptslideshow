package ai

import (
	"context"
	"image"

	"github.com/owulveryck/gptslideshow/internal/structure"
)

// AIInterface defines the interface for interacting with an AI model
// that is capable of extracting text from an audio file, generating slides and presentations,
// as well as generating images from text.
type AIInterface interface {
	// ExtractTextFromAudio extracts text from a given audio file.
	// The filePath parameter specifies the path to the audio file.
	// It returns a string containing the extracted text and an error if any issue occurs.
	ExtractTextFromAudio(ctx context.Context, filePath string) (string, error)

	// GenerateSlide generates a slide from a pre-prompt and specific content.
	// The preprompt parameter allows customization of the slide context.
	// The content parameter is a byte array representing the content to include in the slide.
	// It returns a Slide object representing the generated slide, and any error encountered.
	GenerateSlide(ctx context.Context, preprompt string, content []byte) (*structure.Slide, error)

	// GeneratePresentationFromText generates a presentation from a pre-prompt and given text content.
	// The preprompt parameter allows customization of the presentation context.
	// The content parameter is a byte array representing the content of the presentation.
	// It returns a Presentation object representing the generated presentation, and any error encountered.
	GeneratePresentationFromText(ctx context.Context, preprompt string, content []byte) (*structure.Presentation, error)

	// GenerateImageFromText generates an image from a given text prompt.
	// The prompt parameter provides the description of the image to generate.
	// It returns the generated image as an image.Image object, and any error encountered.
	GenerateImageFromText(ctx context.Context, prompt string) (image.Image, error)
	// SimpleQuery to the AI model
	SimpleQuery(ctx context.Context, prompt string) (string, error)
}
