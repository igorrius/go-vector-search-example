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

func TestGoogleEmbeddingGenerator_Generate(t *testing.T) {
	t.Run("should return embedding when api call is successful", func(t *testing.T) {
		// Arrange
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(`{"embedding":{"values":[0.1, 0.2, 0.3]}}`)),
		}
		httpClient := &http.Client{
			Transport: &mockTransport{response: mockResp},
		}
		opts := option.WithHTTPClient(httpClient)
		generator, err := NewGoogleEmbeddingGenerator(context.Background(), "fake-api-key", opts)
		assert.NoError(t, err)

		// Act
		embedding, err := generator.Generate(context.Background(), "test content")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []float32{0.1, 0.2, 0.3}, embedding)
	})
}
