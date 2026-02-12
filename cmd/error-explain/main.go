package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"
	"strings"

	// Glamour for markdown rendering
	"github.com/charmbracelet/glamour"

	// Import internal packages
	"github.com/ultrakapy/error-explain/internal/config"
	errorContext "github.com/ultrakapy/error-explain/internal/context" 
	"github.com/ultrakapy/error-explain/internal/provider"
	"github.com/ultrakapy/error-explain/internal/runner"
)

// Custom flag type to handle multiple occurrences (e.g., --extra-context file1 --extra-context "string2").
type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var extraContextFlags arrayFlags

func main() {
	mode := flag.String("mode", "direct", "AI Persona: direct, deep, or teacher")
	flag.Var(&extraContextFlags, "extra-context", "Add file path or raw text to the prompt (can be used multiple times)")
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
		fmt.Printf("\n--- 🤖 [AI Thinking...] ---\n")

		// Load Configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("⚠️ Config Warning: %v\nUsing defaults...\n", err)
		}

		// Build the Chain from Config
		var chain []provider.Provider
		for _, pCfg := range cfg.Providers {
			apiKey := os.Getenv(pCfg.APIKeyEnv)
			// Skip if API key is missing (optional safety check)
			if apiKey == "" {
				continue 
			}

			switch pCfg.Type {
			case "anthropic":
				chain = append(chain, &provider.AnthropicProvider{
					APIName: pCfg.Name,
					APIKey:  apiKey,
					Model:   pCfg.Model,
				})
			case "gemini":
				chain = append(chain, &provider.GeminiProvider{
					APIKey: apiKey,
					Model:  pCfg.Model,
				})
			case "openai":
				chain = append(chain, &provider.OpenAICompatibleProvider{
					APIName: pCfg.Name,
					BaseURL: pCfg.BaseURL,
					APIKey:  apiKey,
					Model:   pCfg.Model,
				})
			}
		}

		if len(chain) == 0 {
			fmt.Println("❌ Error: No valid providers found. Check your API keys.")
			os.Exit(1)
		}

		brain := &provider.MultiProvider{Chain: chain}

		// Mine Context & Execute
		sourceContext := errorContext.Mine(result.Stderr)
		extraContext := processExtraContext(extraContextFlags)

		fullPrompt := fmt.Sprintf(
			"Compiler Output:\n%s\n\nSource Code Context:\n%s%s",
			result.Stderr, 
			sourceContext,
			extraContext,
		)
		
		sysPrompt := getSystemPrompt(*mode)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		explanation, err := brain.Explain(ctx, sysPrompt, fullPrompt)
		if err != nil {
			fmt.Printf("❌ AI Failed: %v\n", err)
		} else {
			prettyPrint(explanation)
			// Debugging output:
			//fmt.Printf("DEBUGGING: %s\n", fullPrompt);
		}
	}
	os.Exit(result.ExitCode)
}

// prettyPrint renders markdown with Glamour
func prettyPrint(markdown string) {
	// Create a glamour renderer with auto-detected terminal styling
	r, err := glamour.NewTermRenderer(
		// Use "dark" or "light" style, or glamour.WithAutoStyle() for auto-detection
		glamour.WithAutoStyle(),
		// Wrap text at 100 characters for better readability
		glamour.WithWordWrap(100),
	)
	
	if err != nil {
		// Fallback to plain text if glamour fails
		fmt.Println(markdown)
		return
	}
	
	// Render the markdown
	out, err := r.Render(markdown)
	if err != nil {
		// Fallback to plain text if rendering fails
		fmt.Println(markdown)
		return
	}
	
	// Print the beautifully formatted output
	fmt.Print(out)
}

// Helper to process extra context
func processExtraContext(extras []string) string {
    if len(extras) == 0 {
        return ""
    }

    var sb strings.Builder
    sb.WriteString("\n\n--- USER PROVIDED EXTRA CONTEXT ---\n")

    for _, item := range extras {
        // 1. Check if it looks like a file and exists
        if info, err := os.Stat(item); err == nil && !info.IsDir() {
            content, err := os.ReadFile(item)
            if err == nil {
                sb.WriteString(fmt.Sprintf("\nFile: %s\n```\n%s\n```\n", item, string(content)))
                continue
            }
        }

        // 2. Otherwise, treat as raw text
        sb.WriteString(fmt.Sprintf("\nNote: %s\n", item))
    }
    return sb.String()
}

func getSystemPrompt(mode string) string {
	const contextProtocol = `
If, and only if, the provided context is insufficient to identify the root cause with high confidence:
1. Provide your best guess.
2. Explicitly list what extra information is missing (e.g., "I need to see your go.mod").
3. Advise the user to run this tool again using the '--extra-context' flag to provide that file or information.`

	switch mode {
	case "deep":
		return "You are an expert in this area. Explain the root cause of this error in technical detail. Use markdown formatting with headers, bold text, code blocks, and lists. " + contextProtocol
	case "teacher":
		return "You are a Mentor. Explain this error simply and teach the concept behind it. Use markdown formatting with headers, bold text, code blocks, and lists to make it easy to follow. " + contextProtocol
	case "direct":
		fallthrough
	default:
		return "You are a Build Tool. Fix this error in 1-2 sentences. No fluff. Use markdown code formatting for any code snippets. " + contextProtocol
	}
}
