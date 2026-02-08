package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	// Import internal packages
	errorContext "github.com/ultrakapy/error-explain/internal/context" 
	"github.com/ultrakapy/error-explain/internal/provider"
	"github.com/ultrakapy/error-explain/internal/runner"
)

func main() {
	mode := flag.String("mode", "direct", "AI Persona: direct, deep, or teacher")
	flag.Parse()
	commandArgs := flag.Args()

	if len(commandArgs) == 0 {
		fmt.Println("Usage: error-explain -- <command>")
		os.Exit(1)
	}

	// 1. Run the Compiler (The "Sidecar")
	result, err := runner.Run(commandArgs)
	if err != nil && result.ExitCode == 0 {
		fmt.Printf("❌ System Error: %v\n", err)
		os.Exit(1)
	}

	// 2. If Compiler Failed, Trigger the Voice
	if result.ExitCode != 0 {
		fmt.Printf("\n--- 🤖 [AI Thinking...] ---\n")

		// A. Mine the Code Context (Sprint 3)
		// We pass the raw stderr to the miner.
		sourceContext := errorContext.Mine(result.Stderr)

		// B. Combine Raw Error + Source Code
		// The prompt now has "Eyes"
		fullPrompt := fmt.Sprintf("Compiler Output:\n%s\n\n%s", result.Stderr, sourceContext)
		
		// C. Initialize the Brain (Sprint 2)
		brain := &provider.MultiProvider{
			Chain: []provider.Provider{
				// First Priority: Groq (Free Tier, Fastest)
				&provider.OpenAICompatibleProvider{
					APIName: "Groq",
					BaseURL: "https://api.groq.com/openai/v1",
					APIKey:  os.Getenv("GROQ_API_KEY"),
					Model:   "llama-3.3-70b-versatile",
				},
				// Fallback: Gemini (Free Tier)
				&provider.GeminiProvider{
					APIKey: os.Getenv("GEMINI_API_KEY"),
					//Model:  "gemini-2.0-flash-lite",
					Model: "gemini-3-flash-preview",
				},
			},
		}

		// D. Ask the AI
		sysPrompt := getSystemPrompt(*mode)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		explanation, err := brain.Explain(ctx, sysPrompt, fullPrompt)
		if err != nil {
			fmt.Printf("❌ AI Failed: %v\n", err)
		} else {
			fmt.Println(explanation)
			// For debugging:
			//fmt.Printf("FULL PROMPT:\n%s\n", fullPrompt)
		}
	}

	os.Exit(result.ExitCode)
}

func getSystemPrompt(mode string) string {
	switch mode {
	case "deep":
		return "You are a C++ Expert. Explain the root cause of this error in technical detail."
	case "teacher":
		return "You are a Mentor. Explain this error simply and teach the concept behind it."
	default: // direct
		return "You are a Build Tool. Fix this error in 1-2 sentences. No fluff."
	}
}
