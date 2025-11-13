package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GoogleEmbeddingGenerator generates vector embeddings using the Google AI API.
type GoogleEmbeddingGenerator struct {
	client *genai.EmbeddingModel
}

// NewGoogleEmbeddingGenerator creates a new GoogleEmbeddingGenerator.
func NewGoogleEmbeddingGenerator(ctx context.Context, apiKey string, opts ...option.ClientOption) (*GoogleEmbeddingGenerator, error) {
	opts = append(opts, option.WithAPIKey(apiKey))
	client, err := genai.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new genai client: %w", err)
	}

	return &GoogleEmbeddingGenerator{
		client: client.EmbeddingModel("embedding-001"),
	}, nil
}

// Generate generates a vector embedding for the given content.
func (g *GoogleEmbeddingGenerator) Generate(ctx context.Context, content string) ([]float32, error) {
	res, err := g.client.EmbedContent(ctx, genai.Text(content))
	if err != nil {
		return nil, fmt.Errorf("failed to embed content: %w", err)
	}

	if res == nil || res.Embedding == nil {
		return nil, fmt.Errorf("received an empty embedding from the API")
	}

	return res.Embedding.Values, nil
}
