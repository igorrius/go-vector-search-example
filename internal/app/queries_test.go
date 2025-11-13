package app

import (
	"context"
	"testing"

	"github.com/igorrius/go-vector-search/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks
type MockEmbeddingGenerator struct {
	mock.Mock
}

func (m *MockEmbeddingGenerator) Generate(ctx context.Context, content string) ([]float32, error) {
	args := m.Called(ctx, content)
	return args.Get(0).([]float32), args.Error(1)
}

type MockVectorStore struct {
	mock.Mock
}

func (m *MockVectorStore) Search(ctx context.Context, embedding []float32) ([]domain.Document, error) {
	args := m.Called(ctx, embedding)
	return args.Get(0).([]domain.Document), args.Error(1)
}

type MockSummarizer struct {
	mock.Mock
}

func (m *MockSummarizer) Summarize(ctx context.Context, content []string) (string, error) {
	args := m.Called(ctx, content)
	return args.String(0), args.Error(1)
}

func TestSearchDocumentsHandler_Handle(t *testing.T) {
	ctx := context.Background()
	query := SearchDocumentsQuery{Query: "test query"}
	embedding := []float32{1.0, 2.0, 3.0}
	docs := []domain.Document{
		{ID: "doc1", Content: "This is a test document."},
		{ID: "doc2", Content: "This is another test document."},
	}
	summary := "This is a summary."

	t.Run("Successful search", func(t *testing.T) {
		embedder := new(MockEmbeddingGenerator)
		store := new(MockVectorStore)
		summarizer := new(MockSummarizer)

		embedder.On("Generate", ctx, query.Query).Return(embedding, nil)
		store.On("Search", ctx, embedding).Return(docs, nil)
		var docContents []string
		for _, doc := range docs {
			docContents = append(docContents, doc.Content)
		}
		summarizer.On("Summarize", ctx, docContents).Return(summary, nil)

		handler := NewSearchDocumentsHandler(embedder, store, summarizer)
		result, err := handler.Handle(ctx, query)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, summary, result.Summary)
		assert.Len(t, result.Sources, 2)
		for i, source := range result.Sources {
			assert.Equal(t, docs[i].ID, source.DocumentID)
			assert.Equal(t, docs[i].Content, source.Snippet)
		}

		embedder.AssertExpectations(t)
		store.AssertExpectations(t)
		summarizer.AssertExpectations(t)
	})
}
