package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GeminiProvider handles Google's Generative Language API.
type GeminiProvider struct {
	APIKey string
	Model  string // e.g. "gemini-1.5-flash"
}

func (p *GeminiProvider) Name() string {
	return "Gemini"
}

// Minimal request structure for Gemini REST API
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

func (p *GeminiProvider) Explain(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	// Combine system prompt and user prompt for Gemini's simple API
	fullPrompt := fmt.Sprintf("%s\n\nUser Error:\n%s", systemPrompt, userPrompt)

	reqBody := geminiRequest{
		Contents: []geminiContent{{
			Parts: []geminiPart{{Text: fullPrompt}},
		}},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.Model, p.APIKey)

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error: %d", resp.StatusCode)
	}

	var result geminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini")
	}

	return result.Candidates[0].Content.Parts[0].Text, nil
}
