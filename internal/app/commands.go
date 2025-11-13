package app

import (
	"context"

	"github.com/igorrius/go-vector-search/internal/domain"
)

// IndexDocumentCommand is a command to index a document.
type IndexDocumentCommand struct {
	ID      string
	Content string
}

// IndexDocumentHandler handles the IndexDocumentCommand.
type IndexDocumentHandler struct {
	repo     domain.DocumentRepository
	embedder EmbeddingGenerator
}

// NewIndexDocumentHandler creates a new IndexDocumentHandler.
func NewIndexDocumentHandler(repo domain.DocumentRepository, embedder EmbeddingGenerator) *IndexDocumentHandler {
	return &IndexDocumentHandler{
		repo:     repo,
		embedder: embedder,
	}
}

// Handle handles the IndexDocumentCommand.
func (h *IndexDocumentHandler) Handle(ctx context.Context, cmd IndexDocumentCommand) error {
	doc := domain.NewDocument(cmd.ID, cmd.Content)

	embedding, err := h.embedder.Generate(ctx, doc.Content)
	if err != nil {
		return err
	}

	doc.SetEmbedding(embedding)

	return h.repo.Save(ctx, doc)
}
