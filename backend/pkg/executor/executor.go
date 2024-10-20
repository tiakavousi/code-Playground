package executor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// ExecRequest defines the input for code execution
type ExecRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// ExecuteInteractiveCode runs the submitted code and supports interactive I/O
func ExecuteInteractiveCode(ctx context.Context, req ExecRequest, input <-chan string, output chan<- string) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Generate a unique container name
	containerName := fmt.Sprintf("code-exec-%d", time.Now().UnixNano())
	dockerName := "phantasm/busybox"

	// Prepare Docker command based on language
	var dockerCmd *exec.Cmd
	switch strings.ToLower(req.Language) {
	case "java":
		dockerCmd = prepareJavaCommand(timeoutCtx, containerName, dockerName, req.Code)
	case "c":
		dockerCmd = prepareCCommand(timeoutCtx, containerName, dockerName, req.Code)
	case "c++", "cpp":
		dockerCmd = prepareCppCommand(timeoutCtx, containerName, dockerName, req.Code)
	case "javascript", "js":
		dockerCmd = prepareJavaScriptCommand(timeoutCtx, containerName, dockerName, req.Code)
	default:
		// For interpreted languages, use the existing approach
		dockerCmd = exec.CommandContext(timeoutCtx, "docker", "run", "--rm",
			"--name", containerName,
			"-i", "--cpus=0.5",
			"-m", "100m",
			dockerName,
			strings.ToLower(req.Language), "-c", req.Code)
	}

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

	// Start the Docker command
	if err := dockerCmd.Start(); err != nil {
		output <- "Error starting Docker command: " + err.Error()
		return err
	}

	// Use a WaitGroup to manage goroutines
	var wg sync.WaitGroup
	wg.Add(2)

	// Goroutine for handling stdout and stderr
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
		for scanner.Scan() {
			select {
			case output <- scanner.Text():
			case <-timeoutCtx.Done():
				return
			}
		}
	}()

	// Goroutine for handling stdin
	go func() {
		defer wg.Done()
		for {
			select {
			case inputLine, ok := <-input:
				if !ok {
					return
				}
				_, err := fmt.Fprintln(stdin, inputLine)
				if err != nil {
					output <- "Error writing to stdin: " + err.Error()
					return
				}
			case <-timeoutCtx.Done():
				return
			}
		}
	}()

	// Wait for the command to finish or the context to be done
	done := make(chan error, 1)
	go func() {
		done <- dockerCmd.Wait()
	}()

	// Wait for the command to finish or the context to be done
	select {
	case <-timeoutCtx.Done():
		// Force kill the container
		killCmd := exec.Command("docker", "kill", containerName)
		if err := killCmd.Run(); err != nil {
			output <- fmt.Sprintf("Failed to kill container: %v", err)
		} else {
			output <- "Container killed successfully"
		}

		// Wait for goroutines to finish
		wg.Wait()

		return timeoutCtx.Err()
	case err := <-done:
		// Wait for goroutines to finish
		wg.Wait()
		return err
	}
}

func prepareJavaCommand(ctx context.Context, containerName, dockerName, code string) *exec.Cmd {
	return exec.CommandContext(ctx, "docker", "run", "--rm",
		"--name", containerName,
		"-i", "--cpus=0.5", "-m", "100m",
		"-v", "/tmp:/tmp",
		dockerName,
		"bash", "-c", fmt.Sprintf(`
			echo '%s' > /tmp/Main.java &&
			javac /tmp/Main.java &&
			java -cp /tmp Main
		`, code))
}

func prepareCCommand(ctx context.Context, containerName, dockerName, code string) *exec.Cmd {
	return exec.CommandContext(ctx, "docker", "run", "--rm",
		"--name", containerName,
		"-i", "--cpus=0.5", "-m", "100m",
		"-v", "/tmp:/tmp",
		dockerName,
		"bash", "-c", fmt.Sprintf(`
			echo '%s' > /tmp/main.c &&
			gcc /tmp/main.c -o /tmp/main &&
			/tmp/main
		`, code))
}

func prepareCppCommand(ctx context.Context, containerName, dockerName, code string) *exec.Cmd {
	return exec.CommandContext(ctx, "docker", "run", "--rm",
		"--name", containerName,
		"-i", "--cpus=0.5", "-m", "100m",
		"-v", "/tmp:/tmp",
		dockerName,
		"bash", "-c", fmt.Sprintf(`
			echo '%s' > /tmp/main.cpp &&
			g++ /tmp/main.cpp -o /tmp/main &&
			/tmp/main
		`, code))
}

func prepareJavaScriptCommand(ctx context.Context, containerName, dockerName, code string) *exec.Cmd {
	return exec.CommandContext(ctx, "docker", "run", "--rm",
		"--name", containerName,
		"-i", "--cpus=0.5", "-m", "100m",
		dockerName,
		"node", "-e", code)
}

// ExecuteCode runs the submitted code and returns the output or an error (non-interactive version)
func ExecuteCode(req ExecRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	output := make(chan string, 1)
	input := make(chan string)
	defer close(input)

	errCh := make(chan error, 1)
	go func() {
		start := time.Now()
		err := ExecuteInteractiveCode(ctx, req, input, output)
		fmt.Printf("ExecuteInteractiveCode took %v\n", time.Since(start))
		errCh <- err
	}()

	var result strings.Builder
	for {
		select {
		case line := <-output:
			result.WriteString(line + "\n")
		case err := <-errCh:
			return result.String(), err
		case <-ctx.Done():
			return result.String(), fmt.Errorf("execution timed out after %v: %w", 10*time.Second, ctx.Err())
		}
	}
}
