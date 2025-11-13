package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/igorrius/go-vector-search/internal/app"
	"github.com/igorrius/go-vector-search/internal/domain"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
)

const (
	collectionName = "documents"
)

// TypesenseRepository implements the domain.DocumentRepository and app.VectorStore interfaces.
type TypesenseRepository struct {
	client *typesense.Client
}

// TypesenseConfig holds the configuration for the Typesense client.
type TypesenseConfig struct {
	Host   string
	Port   int
	APIKey string
}

// NewTypesenseRepository creates a new TypesenseRepository.
func NewTypesenseRepository(config TypesenseConfig) (*TypesenseRepository, error) {
	client := typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("http://%s:%d", config.Host, config.Port)),
		typesense.WithAPIKey(config.APIKey),
	)

	repo := &TypesenseRepository{
		client: client,
	}

	var err error
	for i := 0; i < 30; i++ {
		err = repo.ensureCollectionExists(context.Background())
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *TypesenseRepository) ensureCollectionExists(ctx context.Context) error {
	_, err := r.client.Collection(collectionName).Retrieve(ctx)
	if err == nil {
		r.client.Collection(collectionName).Delete(ctx)
	}

	// NumDim should match the dimension of your embeddings, here assumed to be 8.
	// For gemini-embedding-001 model it can be: Flexible, supports: 128 - 3072, Recommended: 768, 1536, 3072
	numDim := 8
	schema := &api.CollectionSchema{
		Name: collectionName,
		Fields: []api.Field{
			{Name: "id", Type: "string"},
			{Name: "content", Type: "string"},
			{Name: "embedding", Type: "float[]", Index: boolPtr(true), Optional: boolPtr(true), NumDim: intPtr(numDim)},
		},
	}

	_, err = r.client.Collections().Create(ctx, schema)
	if err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return err
		}
	}

	return nil
}

// Save persists a document to Typesense.
func (r *TypesenseRepository) Save(ctx context.Context, doc *domain.Document) error {
	document := map[string]interface{}{
		"id":        doc.ID,
		"content":   doc.Content,
		"embedding": doc.Embedding,
	}

	_, err := r.client.Collection(collectionName).Documents().Upsert(ctx, document)
	return err
}

// FindByID retrieves a document from Typesense by its ID.
func (r *TypesenseRepository) FindByID(ctx context.Context, id string) (*domain.Document, error) {
	doc, err := r.client.Collection(collectionName).Document(id).Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	embedding, ok := doc["embedding"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("embedding is not a []interface{}")
	}

	floatEmbedding := make([]float32, len(embedding))
	for i, v := range embedding {
		floatEmbedding[i] = float32(v.(float64))
	}

	return &domain.Document{
		ID:        doc["id"].(string),
		Content:   doc["content"].(string),
		Embedding: floatEmbedding,
	}, nil
}

// Search performs a vector similarity search in Typesense.
func (r *TypesenseRepository) Search(ctx context.Context, embedding []float32) ([]domain.Document, error) {
	vectorQuery := fmt.Sprintf("embedding:([%s], k:10)", floatsToString(embedding))
	searchRequest := &api.SearchCollectionParams{
		Q:           "*",
		QueryBy:     "content",
		VectorQuery: &vectorQuery,
	}

	res, err := r.client.Collection(collectionName).Documents().Search(ctx, searchRequest)
	if err != nil {
		return nil, err
	}

	var documents []domain.Document
	for _, hit := range *res.Hits {
		doc := *hit.Document
		embedding, ok := doc["embedding"].([]interface{})
		if !ok {
			return nil, fmt.Errorf("embedding is not a []interface{}")
		}

		floatEmbedding := make([]float32, len(embedding))
		for i, v := range embedding {
			floatEmbedding[i] = float32(v.(float64))
		}

		documents = append(documents, domain.Document{
			ID:        doc["id"].(string),
			Content:   doc["content"].(string),
			Embedding: floatEmbedding,
		})
	}

	return documents, nil
}

func floatsToString(floats []float32) string {
	var b strings.Builder
	for i, f := range floats {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%f", f))
	}
	return b.String()
}

var _ domain.DocumentRepository = (*TypesenseRepository)(nil)
var _ app.VectorStore = (*TypesenseRepository)(nil)

func boolPtr(b bool) *bool {
	return &b
}

func intPtr(i int) *int {
	return &i
}
