package provider

import (
	"context"
	"fmt"
	"log"
)

// MultiProvider wraps a list of providers and tries them in order.
type MultiProvider struct {
	Chain []Provider
}

func (mp *MultiProvider) Explain(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	var lastErr error

	for _, p := range mp.Chain {
		// Log which provider we are attempting (useful for debugging)
		// In production, you might hide this unless --verbose is on.
		log.Printf("Trying provider: %s...", p.Name())

		explanation, err := p.Explain(ctx, systemPrompt, userPrompt)
		if err == nil {
			return explanation, nil
		}

		// Log the failure but continue
		log.Printf("⚠️ %s failed: %v", p.Name(), err)
		lastErr = err
	}

	return "", fmt.Errorf("all providers failed. Last error: %w", lastErr)
}
