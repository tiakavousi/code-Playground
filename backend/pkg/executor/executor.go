package executor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// ExecRequest defines the input for code execution
type ExecRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// ExecuteInteractiveCode runs the submitted code and supports interactive I/O
func ExecuteInteractiveCode(ctx context.Context, req ExecRequest, input <-chan string, output chan<- string) error {
	var err error

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	// Generate a unique container name
	containerName := fmt.Sprintf("code-exec-%d", time.Now().UnixNano())

	// Construct the Docker command with resource limits and your Docker image
	dockerCmd := exec.CommandContext(timeoutCtx, "docker", "run", "--rm",
		"--name", containerName, // Add a name to the container
		"-i", "--cpus=0.5", // Limit CPU usage to 0.5 cores
		"-m", "100m", // Limit memory usage to 100MB
		"phantasm/busybox", // Use the correct image name
		strings.ToLower(req.Language), "-c", req.Code)

	stdin, err := dockerCmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	stdout, err := dockerCmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %w", err)
	}

	stderr, err := dockerCmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %w", err)
	}

	done := make(chan error, 1)

	// Start the Docker command in a separate goroutine
	go func() {
		if err := dockerCmd.Start(); err != nil {
			output <- "Error starting Docker command: " + err.Error()
			return
		}

		// Wait for the command to finish
		done <- dockerCmd.Wait()
	}()

	go func() {
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			output <- scanner.Text()
		}
	}()

	go func() {
		for inputLine := range input {
			fmt.Fprintln(stdin, inputLine)
		}
		stdin.Close()
	}()

	// Wait for the command to finish or the context to be done
	select {
	case <-timeoutCtx.Done():
		fmt.Println("Timeout reached. Killing container...")

		// Kill the Docker container using its name
		killCmd := exec.Command("docker", "kill", containerName)
		if err := killCmd.Run(); err != nil {
			fmt.Printf("Failed to kill container: %v\n", err)
		} else {
			fmt.Println("Container killed successfully")
		}

		return timeoutCtx.Err()
	case err := <-done:
		return err
	}
}

// ExecuteCode runs the submitted code and returns the output or an error (non-interactive version)
func ExecuteCode(req ExecRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	output := make(chan string)
	input := make(chan string)
	defer close(input)

	errCh := make(chan error, 1)
	go func() {
		errCh <- ExecuteInteractiveCode(ctx, req, input, output)
	}()

	var result strings.Builder
	for {
		select {
		case line := <-output:
			result.WriteString(line + "\n")
		case err := <-errCh:
			return result.String(), err
		case <-ctx.Done():
			return result.String(), ctx.Err()
		}
	}
}
