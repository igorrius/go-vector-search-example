package domain

import "context"

// Document is the aggregate root for our domain.
type Document struct {
	ID        string
	Content   string
	Embedding []float32
}

// NewDocument creates a new Document.
func NewDocument(id, content string) *Document {
	return &Document{
		ID:      id,
		Content: content,
	}
}

// SetEmbedding sets the vector embedding for the document.
func (d *Document) SetEmbedding(embedding []float32) {
	d.Embedding = embedding
}

// DocumentRepository defines the contract for storing and retrieving Document aggregates.
type DocumentRepository interface {
	Save(ctx context.Context, doc *Document) error
	FindByID(ctx context.Context, id string) (*Document, error)
}
