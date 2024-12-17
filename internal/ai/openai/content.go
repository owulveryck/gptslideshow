package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/openai/openai-go"

	"github.com/owulveryck/gptslideshow/config"
	"github.com/owulveryck/gptslideshow/internal/structure"
)

// sanitize remove unwanted characters (example: control characters)
func sanitize(input string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 { // Remove control characters
			return -1
		}
		return r
	}, input)
}

func (ai *AI) SimpleQuery(ctx context.Context, prompt string) (string, error) {
	completion, err := ai.Client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(sanitize(prompt)),
		}),
		Seed:  openai.Int(1),
		Model: openai.F(config.ConfigInstance.OpenAIModel),
	})
	if err != nil {
		return "", fmt.Errorf("Error in the chat request: %v\nprompt was:\n%v", err, prompt)
	}

	return completion.Choices[0].Message.Content, nil
}

// GenerateSlide
func (ai *AI) GenerateSlide(ctx context.Context, preprompt string, content []byte) (*structure.Slide, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("presentation"),
		Description: openai.F("A structured slide from content"),
		Schema:      openai.F(structure.SlideResponseSchema),
		Strict:      openai.Bool(true),
	}

	// Prompt to guide the model
	prompt := fmt.Sprintf(preprompt+`

		%s`, string(content))
	log.Printf("\n\nPrompting with: %s ...\n\n", prompt[:50])

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
	var slide structure.Slide
	err = json.Unmarshal([]byte(chat.Choices[0].Message.Content), &slide)
	if err != nil {
		return nil, err
	}
	return &slide, err
}

// GeneratePresentationFromText generates a presentation from Markdown content
func (ai *AI) GeneratePresentationFromText(ctx context.Context, preprompt string, content []byte) (*structure.Presentation, error) {
	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:        openai.F("presentation"),
		Description: openai.F("A structured presentation from content"),
		Schema:      openai.F(structure.PresentationResponseSchema),
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
