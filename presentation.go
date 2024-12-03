package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/invopop/jsonschema"
	"github.com/openai/openai-go"
)

// Presentation represents the entire presentation structure
type Presentation struct {
	Title  string  `json:"presentation_title" jsonschema_description:"The title of the presentation"`
	Slides []Slide `json:"slides" jsonschema_description:"The content of the presentation"`
}

// Slide represents a single slide in the presentation
type Slide struct {
	Title    string `json:"title" jsonschema_description:"The title of the slide"`
	Subtitle string `json:"subtitle" jsonschema_description:"The subtitle of the slide"`
	Body     string `json:"body" jsonschema_description:"The main content of the slide"`
}

// GenerateSchema generates the JSON schema for a given type
func GenerateSchema[T any]() interface{} {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	return reflector.Reflect(v)
}

// Generate the JSON schema for Slide
var SlideResponseSchema = GenerateSchema[Presentation]()

// GenerateSlides generates a presentation from Markdown content
func GenerateSlides(ctx context.Context, content []byte) (*Presentation, error) {
	client := openai.NewClient()
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("presentation"),
		Description: openai.F("A structured presentation from Markdown content"),
		Schema:      openai.F(SlideResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Prompt to guide the model
	prompt := fmt.Sprintf(`Convert the following Markdown text into an array of structured slides.
Each slide should have a title, a subtitle, and a body:

%s`, string(content))
	log.Printf("\n\nPrompting with: %s ...\n\n", prompt[:300])

	// Query OpenAI API for validation or enhancement (optional)
	chat, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		return nil, err
	}

	// Parse the model's response
	var presentation Presentation
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &presentation)
	if err != nil {
		return nil, err
	}
	log.Printf("Generated %d slides", len(presentation.Slides))
	return &presentation, err
}
