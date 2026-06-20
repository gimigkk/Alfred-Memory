package embed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GeminiClient struct {
	APIKey string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{APIKey: apiKey}
}

func (c *GeminiClient) GetVector(text string) ([]float32, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/gemini-embedding-2:embedContent?key=%s", c.APIKey)

	reqBody := map[string]interface{}{
		"model": "models/gemini-embedding-2",
		"content": map[string]interface{}{
			"parts": []map[string]interface{}{
				{"text": text},
			},
		},
	}

	jsonValue, _ := json.Marshal(reqBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	embeddingData, ok := result["embedding"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected gemini response format: missing 'embedding' field")
	}

	valuesData, ok := embeddingData["values"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected gemini response format: missing 'values' array")
	}

	vector := make([]float32, len(valuesData))
	for i, v := range valuesData {
		num, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("expected float number in vector at index %d", i)
		}
		vector[i] = float32(num)
	}

	return vector, nil
}
