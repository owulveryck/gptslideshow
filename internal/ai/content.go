package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/openai/openai-go"

	"github.com/owulveryck/gptslideshow/config"
	"github.com/owulveryck/gptslideshow/internal/structure"
)

// GenerateContentFromText generates a presentation from Markdown content
func (ai *AI) GenerateContentFromText(ctx context.Context, preprompt string, content []byte) (*structure.Presentation, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("presentation"),
		Description: openai.F("A structured presentation from content"),
		Schema:      openai.F(structure.SlideResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Prompt to guide the model
	prompt := fmt.Sprintf(preprompt+`

		%s`, string(content))
	log.Printf("\n\nPrompting with: %s ...\n\n", prompt[:500])

	// Query OpenAI API for validation or enhancement (optional)
	chat, err := ai.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(prompt),
		}),
		ResponseFormat: openai.F[openai.ChatCompletionNewParamsResponseFormatUnion](
			openai.ResponseFormatJSONSchemaParam{
				Type:       openai.F(openai.ResponseFormatJSONSchemaTypeJSONSchema),
				JSONSchema: openai.F(schemaParam),
			},
		),
		Model: openai.F(config.ConfigInstance.OpenAIModel),
		// Model: openai.F(openai.ChatModelGPT4o2024_08_06),
	})
	if err != nil {
		return nil, err
	}

	// Parse the model's response
	var presentation structure.Presentation
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &presentation)
	if err != nil {
		return nil, err
	}
	log.Printf("Generated %d slides", len(presentation.Slides))
	return &presentation, err
}
