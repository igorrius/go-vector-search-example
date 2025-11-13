package ai

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// GoogleSummarizer summarizes text using the Google AI API.
type GoogleSummarizer struct {
	client *genai.GenerativeModel
}

// NewGoogleSummarizer creates a new GoogleSummarizer.
func NewGoogleSummarizer(ctx context.Context, apiKey string, opts ...option.ClientOption) (*GoogleSummarizer, error) {
	opts = append(opts, option.WithAPIKey(apiKey))
	client, err := genai.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create new genai client: %w", err)
	}

	return &GoogleSummarizer{
		client: client.GenerativeModel("gemini-pro"),
	}, nil
}

// Summarize summarizes the given content.
func (s *GoogleSummarizer) Summarize(ctx context.Context, content []string) (string, error) {
	prompt := fmt.Sprintf("Provide a concise summary of the following documents:\n\n%s", strings.Join(content, "\n---\n"))
	resp, err := s.client.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("received an empty response from the API")
	}

	var summary strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					summary.WriteString(string(txt))
				}
			}
		}
	}

	if summary.Len() == 0 {
		return "", fmt.Errorf("unexpected response format from the API, no text part found")
	}

	return summary.String(), nil
}
