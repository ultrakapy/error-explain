package runner

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Result holds the captured data after execution
type Result struct {
	ExitCode int
	Stderr   string // This is what we will feed to the AI
}

// Run executes the command while streaming output to the user AND capturing it.
func Run(command []string) (Result, error) {
	if len(command) == 0 {
		return Result{}, fmt.Errorf("no command provided")
	}

	// 1. Prepare the subprocess
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout // Direct passthrough for standard output

	// 2. The "Transparent Pipe" Architecture
	// We create a buffer to hold the error text for the AI.
	var stderrCapture bytes.Buffer
	
	// MultiWriter sends data to TWO places at once:
	// A. os.Stderr (The User's Screen) - so they see errors instantly.
	// B. stderrCapture (The AI's Ear) - so we can explain it later.
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderrCapture)

	// 3. Execute
	err := cmd.Run()

	// 4. Extract Exit Code
	exitCode := 0
	if err != nil {
		// Try to get the actual exit code from the error
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// If the command failed to start (e.g. "g++ not found"), return error
			return Result{ExitCode: 1, Stderr: err.Error()}, err
		}
	}

	return Result{
		ExitCode: exitCode,
		Stderr:   stderrCapture.String(),
	}, nil
}
