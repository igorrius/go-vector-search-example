package app

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
)

// HTTPHandlers holds the command and query handlers.
type HTTPHandlers struct {
	indexDocumentHandler   *IndexDocumentHandler
	searchDocumentsHandler *SearchDocumentsHandler
}

// NewHTTPHandlers creates a new HTTPHandlers.
func NewHTTPHandlers(indexDocumentHandler *IndexDocumentHandler, searchDocumentsHandler *SearchDocumentsHandler) *HTTPHandlers {
	return &HTTPHandlers{
		indexDocumentHandler:   indexDocumentHandler,
		searchDocumentsHandler: searchDocumentsHandler,
	}
}

// IndexDocumentRequest is the request body for indexing a document.
type IndexDocumentRequest struct {
	ID      string `json:"id"`
	Content string `json:"content"`
}

// IndexDocumentHandler handles the POST /api/v1/documents endpoint.
func (h *HTTPHandlers) IndexDocumentHandler(w http.ResponseWriter, r *http.Request) {
	var cmd IndexDocumentCommand

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		var req IndexDocumentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		cmd.ID = req.ID
		cmd.Content = req.Content
	} else if _, _, err := r.FormFile("file"); err == nil {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Invalid file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		cmd.Content = string(content)
		cmd.ID = r.FormValue("id")
	} else {
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}

	if cmd.ID == "" {
		cmd.ID = uuid.New().String()
	}

	if err := h.indexDocumentHandler.Handle(r.Context(), cmd); err != nil {
		http.Error(w, "Failed to index document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

// SearchDocumentsHandler handles the GET /api/v1/search endpoint.
func (h *HTTPHandlers) SearchDocumentsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	searchQuery := SearchDocumentsQuery{Query: query}
	result, err := h.searchDocumentsHandler.Handle(r.Context(), searchQuery)
	if err != nil {
		http.Error(w, "Failed to search documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
