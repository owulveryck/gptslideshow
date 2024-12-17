package ollama

import (
	"context"
	"image"
	"log"

	"github.com/ollama/ollama/api"
	"github.com/owulveryck/gptslideshow/internal/structure"
)

type Ollama struct{}

// ExtractTextFromAudio extracts text from a given audio file.
// The filePath parameter specifies the path to the audio file.
// It returns a string containing the extracted text and an error if any issue occurs.
func (ollama *Ollama) ExtractTextFromAudio(ctx context.Context, filePath string) (string, error) {
	panic("not implemented") // TODO: Implement
}

// GenerateSlide generates a slide from a pre-prompt and specific content.
// The preprompt parameter allows customization of the slide context.
// The content parameter is a byte array representing the content to include in the slide.
// It returns a Slide object representing the generated slide, and any error encountered.
func (ollama *Ollama) GenerateSlide(ctx context.Context, preprompt string, content []byte) (*structure.Slide, error) {
	panic("not implemented") // TODO: Implement
}

// GeneratePresentationFromText generates a presentation from a pre-prompt and given text content.
// The preprompt parameter allows customization of the presentation context.
// The content parameter is a byte array representing the content of the presentation.
// It returns a Presentation object representing the generated presentation, and any error encountered.
func (ollama *Ollama) GeneratePresentationFromText(ctx context.Context, preprompt string, content []byte) (*structure.Presentation, error) {
	panic("not implemented") // TODO: Implement
}

// GenerateImageFromText generates an image from a given text prompt.
// The prompt parameter provides the description of the image to generate.
// It returns the generated image as an image.Image object, and any error encountered.
func (ollama *Ollama) GenerateImageFromText(ctx context.Context, prompt string) (image.Image, error) {
	panic("not implemented") // TODO: Implement
}

type AI struct {
	Client *api.Client
	Model  string
}

func NewAI() *AI {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	return &AI{
		Client: client,
		Model:  "llama3.2",
	}
}

func (ai *AI) SimpleQuery(ctx context.Context, prompt string) (string, error) {
	// Define the model and prompt
	req := &api.GenerateRequest{
		Model:  ai.Model,
		Prompt: prompt,

		// set streaming to false
		Stream: new(bool),
	}
	var reply string
	respFunc := func(resp api.GenerateResponse) error {
		// Only print the response here; GenerateResponse has a number of other
		// interesting fields you want to examine.
		reply = resp.Response
		return nil
	}

	err := ai.Client.Generate(ctx, req, respFunc)
	if err != nil {
		return "", nil
	}
	return reply, nil
}
