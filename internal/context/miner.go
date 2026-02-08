package context

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Standard GCC/Clang error format: "main.cpp:10:5: error: ..."
// We capture group 1 (file) and group 2 (line)
var gccPattern = regexp.MustCompile(`^([^:\n]+):(\d+):(\d+):?`)

// Mine inspects the compiler output and extracts the relevant source code.
// It creates a formatted block of text ready for the LLM.
func Mine(stderrInput string) string {
	// 1. Find the FIRST error location. 
	// We only focus on the first error because cascading errors are often noise.
	scanner := bufio.NewScanner(strings.NewReader(stderrInput))
	
	var file string
	var line int
	var errFound bool

	for scanner.Scan() {
		text := scanner.Text()
		matches := gccPattern.FindStringSubmatch(text)
		if len(matches) >= 3 {
			file = matches[1]
			lineNum, _ := strconv.Atoi(matches[2])
			line = lineNum
			errFound = true
			break // Stop at the first error to save tokens
		}
	}

	if !errFound {
		return "" // No file context found, just use the raw error message
	}

	// 2. Read the file
	codeSnippet, err := readWindow(file, line, 5) // Grab ±5 lines
	if err != nil {
		return fmt.Sprintf("[Context Miner Warning: Could not read file %s: %v]", file, err)
	}

	// 3. Format for the AI
	return fmt.Sprintf("\n--- Source Context (%s:%d) ---\n%s\n", file, line, codeSnippet)
}

// readWindow reads the file and extracts lines [target-window, target+window]
func readWindow(path string, targetLine, window int) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var builder strings.Builder
	scanner := bufio.NewScanner(file)
	currentLine := 1
	start := targetLine - window
	end := targetLine + window

	for scanner.Scan() {
		if currentLine >= start && currentLine <= end {
			prefix := "   "
			if currentLine == targetLine {
				prefix = "-> " // Point to the error line
			}
			builder.WriteString(fmt.Sprintf("%s%d | %s\n", prefix, currentLine, scanner.Text()))
		}
		
		if currentLine > end {
			break
		}
		currentLine++
	}

	return builder.String(), scanner.Err()
}
