package provider

import "context"

// Provider defines the behavior for any AI backend.
type Provider interface {
	// Name returns the identifier (e.g., "groq", "gemini") for logging.
	Name() string
	
	// Explain sends the prompt and returns the text response.
	Explain(ctx context.Context, systemPrompt, userPrompt string) (string, error)
}
