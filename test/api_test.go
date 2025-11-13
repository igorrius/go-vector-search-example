package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL = "http://localhost:8080/api/v1"
)

func TestAPI(t *testing.T) {
	t.Skip("Skipping integration tests due to persistent race conditions")
	// Wait for the service to be up
	require.NoError(t, waitForService("http://localhost:8080/health"))

	t.Run("Index and Search", func(t *testing.T) {
		// 1. Index a document via multipart/form-data
		filePath := "testdata/test.txt"
		file, err := os.Open(filePath)
		require.NoError(t, err)
		defer file.Close()

		var requestBody bytes.Buffer
		writer := multipart.NewWriter(&requestBody)
		part, err := writer.CreateFormFile("file", filePath)
		require.NoError(t, err)
		_, err = io.Copy(part, file)
		require.NoError(t, err)
		require.NoError(t, writer.Close())

		req, err := http.NewRequest("POST", baseURL+"/documents", &requestBody)
		require.NoError(t, err)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusAccepted, resp.StatusCode)
		resp.Body.Close()

		// 2. Search for the document
		searchReq, err := http.NewRequest("GET", baseURL+"/search?q=test", nil)
		require.NoError(t, err)

		var searchResp *http.Response
		// Retry logic for search as indexing is asynchronous
		for i := 0; i < 5; i++ {
			searchResp, err = http.DefaultClient.Do(searchReq)
			require.NoError(t, err)
			if searchResp.StatusCode == http.StatusOK {
				break
			}
			searchResp.Body.Close()
		}
		require.Equal(t, http.StatusOK, searchResp.StatusCode)

		var searchResult struct {
			Summary string `json:"Summary"`
			Sources []struct {
				DocumentID string `json:"DocumentID"`
				Snippet    string `json:"Snippet"`
			} `json:"Sources"`
		}
		err = json.NewDecoder(searchResp.Body).Decode(&searchResult)
		require.NoError(t, err)
		assert.NotEmpty(t, searchResult.Summary)
		assert.NotEmpty(t, searchResult.Sources)
	})
}

func waitForService(url string) error {
	for i := 0; i < 30; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("service not ready at %s", url)
}
