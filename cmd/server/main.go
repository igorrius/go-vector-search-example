package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/igorrius/go-vector-search/internal/app"
	"github.com/igorrius/go-vector-search/internal/infra/ai"
	"github.com/igorrius/go-vector-search/internal/infra/persistence"
)

type config struct {
	httpPort        int
	typesenseHost   string
	typesensePort   int
	typesenseAPIKey string
	googleAIApiKey  string
}

func loadConfig() config {
	httpPort, _ := strconv.Atoi(getEnv("HTTP_PORT", "8080"))
	typesensePort, _ := strconv.Atoi(getEnv("TYPESENSE_PORT", "8080"))

	return config{
		httpPort:        httpPort,
		typesenseHost:   getEnv("TYPESENSE_HOST", "localhost"),
		typesensePort:   typesensePort,
		typesenseAPIKey: getEnv("TYPESENSE_API_KEY", ""),
		googleAIApiKey:  getEnv("GOOGLE_API_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	cfg := loadConfig()
	ctx := context.Background()

	// Initialize infrastructure components
	typesenseRepo, err := persistence.NewTypesenseRepository(persistence.TypesenseConfig{
		Host:   cfg.typesenseHost,
		Port:   cfg.typesensePort,
		APIKey: cfg.typesenseAPIKey,
	})
	if err != nil {
		log.Fatalf("Failed to create Typesense repository: %v", err)
	}

	embeddingGenerator, err := ai.NewGoogleEmbeddingGenerator(ctx, cfg.googleAIApiKey)
	if err != nil {
		log.Fatalf("Failed to create Google embedding generator: %v", err)
	}

	summarizer, err := ai.NewGoogleSummarizer(ctx, cfg.googleAIApiKey)
	if err != nil {
		log.Fatalf("Failed to create Google summarizer: %v", err)
	}

	// Initialize application handlers
	indexDocumentHandler := app.NewIndexDocumentHandler(typesenseRepo, embeddingGenerator)
	searchDocumentsHandler := app.NewSearchDocumentsHandler(embeddingGenerator, typesenseRepo, summarizer)
	httpHandlers := app.NewHTTPHandlers(indexDocumentHandler, searchDocumentsHandler)

	// Set up HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/documents", httpHandlers.IndexDocumentHandler).Methods("POST")
	router.HandleFunc("/api/v1/search", httpHandlers.SearchDocumentsHandler).Methods("GET")
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")

	// Start server
	addr := fmt.Sprintf(":%d", cfg.httpPort)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
