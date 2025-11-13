package app_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/igorrius/go-vector-search/internal/app"
	"github.com/igorrius/go-vector-search/internal/domain"
)

// MockDocumentRepository is a mock for the DocumentRepository interface.
type MockDocumentRepository struct {
	mock.Mock
}

func (m *MockDocumentRepository) Save(ctx context.Context, doc *domain.Document) error {
	args := m.Called(ctx, doc)
	return args.Error(0)
}

func (m *MockDocumentRepository) FindByID(ctx context.Context, id string) (*domain.Document, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Document), args.Error(1)
}

// MockEmbeddingGenerator is a mock for the EmbeddingGenerator interface.
type MockEmbeddingGenerator struct {
	mock.Mock
}

func (m *MockEmbeddingGenerator) Generate(ctx context.Context, content string) ([]float32, error) {
	args := m.Called(ctx, content)
	return args.Get(0).([]float32), args.Error(1)
}

func TestIndexDocumentHandler_Handle(t *testing.T) {
	ctx := context.Background()
	repo := new(MockDocumentRepository)
	embedder := new(MockEmbeddingGenerator)
	handler := app.NewIndexDocumentHandler(repo, embedder)

	cmd := app.IndexDocumentCommand{
		ID:      "test-id",
		Content: "test content",
	}

	embedding := []float32{1.0, 2.0, 3.0}
	doc := domain.NewDocument(cmd.ID, cmd.Content)
	doc.SetEmbedding(embedding)

	embedder.On("Generate", ctx, cmd.Content).Return(embedding, nil)
	repo.On("Save", ctx, doc).Return(nil)

	err := handler.Handle(ctx, cmd)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	embedder.AssertExpectations(t)
}
