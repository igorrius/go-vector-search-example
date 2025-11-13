package test

import (
	"context"
	"testing"

	"github.com/igorrius/go-vector-search/internal/domain"
	"github.com/igorrius/go-vector-search/internal/infra/persistence"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypesenseRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	config := persistence.TypesenseConfig{
		Host:   "localhost",
		Port:   8108,
		APIKey: "xyz",
	}

	repo, err := persistence.NewTypesenseRepository(config)
	require.NoError(t, err)

	doc := &domain.Document{
		ID:        "test-id",
		Content:   "this is a test document",
		Embedding: []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8},
	}

	// Save the document
	err = repo.Save(context.Background(), doc)
	require.NoError(t, err)

	// Find the document by ID
	foundDoc, err := repo.FindByID(context.Background(), "test-id")
	require.NoError(t, err)
	assert.Equal(t, doc.ID, foundDoc.ID)
	assert.Equal(t, doc.Content, foundDoc.Content)

	// Search for the document
	results, err := repo.Search(context.Background(), []float32{0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8})
	require.NoError(t, err)
	require.NotEmpty(t, results)
	for _, result := range results {
		assert.Equal(t, "test-id", result.ID)
		assert.Equal(t, "this is a test document", result.Content)
		assert.InDeltaSlice(t, doc.Embedding, result.Embedding, 0.001)
	}
}
