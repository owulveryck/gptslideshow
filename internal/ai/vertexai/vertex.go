package vertex

import (
	"context"
	"fmt"
	"image"
	"log"

	"cloud.google.com/go/vertexai/genai"
	"github.com/owulveryck/gptslideshow/internal/structure"
)

type AI struct {
	Client *genai.Client
	Gemini *genai.GenerativeModel
}

// ExtractTextFromAudio extracts text from a given audio file.
// The filePath parameter specifies the path to the audio file.
// It returns a string containing the extracted text and an error if any issue occurs.
func (ai *AI) ExtractTextFromAudio(ctx context.Context, filePath string) (string, error) {
	panic("not implemented") // TODO: Implement
}

// GenerateSlide generates a slide from a pre-prompt and specific content.
// The preprompt parameter allows customization of the slide context.
// The content parameter is a byte array representing the content to include in the slide.
// It returns a Slide object representing the generated slide, and any error encountered.
func (ai *AI) GenerateSlide(ctx context.Context, preprompt string, content []byte) (*structure.Slide, error) {
	panic("not implemented") // TODO: Implement
}

// GeneratePresentationFromText generates a presentation from a pre-prompt and given text content.
// The preprompt parameter allows customization of the presentation context.
// The content parameter is a byte array representing the content of the presentation.
// It returns a Presentation object representing the generated presentation, and any error encountered.
func (ai *AI) GeneratePresentationFromText(ctx context.Context, preprompt string, content []byte) (*structure.Presentation, error) {
	panic("not implemented") // TODO: Implement
}

// GenerateImageFromText generates an image from a given text prompt.
// The prompt parameter provides the description of the image to generate.
// It returns the generated image as an image.Image object, and any error encountered.
func (ai *AI) GenerateImageFromText(ctx context.Context, prompt string) (image.Image, error) {
	panic("not implemented") // TODO: Implement
}

// SimpleQuery to the AI model
func (ai *AI) SimpleQuery(ctx context.Context, prompt string) (string, error) {
	resp, err := ai.Gemini.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("error generating content: %w", err)
	}
	result := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if part, ok := part.(genai.Text); ok {
			result += string(part)
		}
	}

	return result, nil
}

func NewAI(ctx context.Context, projectID, locationID, modelID string) *AI {
	client, err := genai.NewClient(ctx, projectID, locationID)
	if err != nil {
		log.Fatalf("Failed to create the client: %v", err)
	}
	return &AI{
		Client: client,
		Gemini: client.GenerativeModel(modelID),
	}
}
