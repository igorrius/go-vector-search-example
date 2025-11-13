package app

import (
	"context"
)

// SearchDocumentsQuery represents a query to search for documents.
type SearchDocumentsQuery struct {
	Query string
}

// SearchResult represents the result of a document search.
type SearchResult struct {
	Summary string
	Sources []Source
}

// Source represents a source document for a search result.
type Source struct {
	DocumentID string
	Snippet    string
}

// SearchDocumentsHandler handles the SearchDocumentsQuery.
type SearchDocumentsHandler struct {
	embedder   EmbeddingGenerator
	store      VectorStore
	summarizer Summarizer
}

// NewSearchDocumentsHandler creates a new SearchDocumentsHandler.
func NewSearchDocumentsHandler(embedder EmbeddingGenerator, store VectorStore, summarizer Summarizer) *SearchDocumentsHandler {
	return &SearchDocumentsHandler{
		embedder:   embedder,
		store:      store,
		summarizer: summarizer,
	}
}

// Handle handles the SearchDocumentsQuery.
func (h *SearchDocumentsHandler) Handle(ctx context.Context, query SearchDocumentsQuery) (*SearchResult, error) {
	embedding, err := h.embedder.Generate(ctx, query.Query)
	if err != nil {
		return nil, err
	}

	docs, err := h.store.Search(ctx, embedding)
	if err != nil {
		return nil, err
	}

	var content []string
	for _, doc := range docs {
		content = append(content, doc.Content)
	}

	summary, err := h.summarizer.Summarize(ctx, content)
	if err != nil {
		return nil, err
	}

	var sources []Source
	for _, doc := range docs {
		sources = append(sources, Source{
			DocumentID: doc.ID,
			Snippet:    doc.Content, // Using full content as snippet for now
		})
	}

	return &SearchResult{
		Summary: summary,
		Sources: sources,
	}, nil
}
