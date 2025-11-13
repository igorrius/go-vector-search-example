package ai

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestGoogleSummarizer_Summarize(t *testing.T) {
	t.Run("should return summary when api call is successful", func(t *testing.T) {
		// Arrange
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"candidates":[{"content":{"parts":[{"text":"This is a summary."}]}}]}`)),
		}
		httpClient := &http.Client{
			Transport: &mockTransport{response: mockResp},
		}
		opts := option.WithHTTPClient(httpClient)
		summarizer, err := NewGoogleSummarizer(context.Background(), "fake-api-key", opts)
		assert.NoError(t, err)

		// Act
		summary, err := summarizer.Summarize(context.Background(), []string{"doc1", "doc2"})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "This is a summary.", summary)
	})
}
