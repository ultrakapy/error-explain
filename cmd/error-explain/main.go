package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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

	// 1. Run the Compiler
	result, err := runner.Run(commandArgs)
	if err != nil && result.ExitCode == 0 {
		fmt.Printf("❌ System Error: %v\n", err)
		os.Exit(1)
	}

	if result.ExitCode != 0 {
		// 2. Initialize the AI Brain (The Failover Chain)
		// In a real app, you'd load these from a config file.
		brain := &provider.MultiProvider{
			Chain: []provider.Provider{
				// First Priority: Groq (Fastest)
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

		// 3. Ask for Help
		fmt.Printf("\n--- 🤖 [AI Thinking...] ---\n")
		
		sysPrompt := getSystemPrompt(*mode)
		
		// Use a background context or one with a timeout
		ctx := context.Background()
		explanation, err := brain.Explain(ctx, sysPrompt, result.Stderr)
		
		if err != nil {
			fmt.Printf("❌ AI Failed: %v\n", err)
		} else {
			fmt.Println(explanation)
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
