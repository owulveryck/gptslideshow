package ai

import (
	"net/http"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

// AI represents a client for interacting with OpenAI's API.
type AI struct {
	Client *openai.Client
}

func NewAI() *AI {
	// Create a custom HTTP client with a 5-minute timeout.
	httpClient := &http.Client{
		Timeout: 5 * time.Minute,
	}

	// Create a new OpenAI client using the custom HTTP client.
	client := openai.NewClient(
		option.WithHTTPClient(httpClient),
	)

	return &AI{client}
}
