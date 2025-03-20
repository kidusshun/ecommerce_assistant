package embedding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func GetEmbedding(prompt string) (*EmbeddingResponse, error) {
	// Define the API endpoint
	url := "http://localhost:11434/api/embeddings"

	// Create the request payload
	requestPayload := EmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: prompt,
	}

	// Marshal the payload into JSON
	payload, err := json.Marshal(requestPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Send the HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Read response body for additional context
		return nil, fmt.Errorf("received non-OK response: %s, body: %s", resp.Status, string(body))
	}

	// Parse the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var embeddingResponse EmbeddingResponse
	if err := json.Unmarshal(body, &embeddingResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &embeddingResponse, nil
}
