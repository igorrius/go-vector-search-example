package app

import (
	"context"

	"github.com/igorrius/go-vector-search/internal/domain"
)

// EmbeddingGenerator generates a vector embedding for a given content.
type EmbeddingGenerator interface {
	Generate(ctx context.Context, content string) ([]float32, error)
}

// VectorStore defines the interface for a vector store.
type VectorStore interface {
	Search(ctx context.Context, embedding []float32) ([]domain.Document, error)
}

// Summarizer defines the interface for a text summarizer.
type Summarizer interface {
	Summarize(ctx context.Context, content []string) (string, error)
}
