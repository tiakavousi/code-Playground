// core Docker runner functionality
package container

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/tiakavousi/codeplayground/pkg/executor"
)

// Security and resource constants
const (
	maxCPU    = "0.5"
	maxMemory = "100m"
	maxPids   = "20"
)

// DockerRunner implements the executor.CodeRunner interface
type DockerRunner struct {
	imageName    string
	securityOpts []string
}

// NewDockerRunner creates a new Docker-based code runner
func NewDockerRunner(imageName string) *DockerRunner {
	return &DockerRunner{
		imageName: imageName,
		securityOpts: []string{
			"--cap-drop=ALL",
			"--net=none",
			"--pids-limit=" + maxPids,
			"--ulimit", "nproc=20:20",
			"--ulimit", "nofile=64:64",
			"--ulimit", "fsize=1000000:1000000",
		},
	}
}

func (d *DockerRunner) RunInteractive(ctx context.Context, req executor.ExecRequest, input <-chan string, output chan<- string) error {
	containerName := fmt.Sprintf("code-exec-%d", time.Now().UnixNano())

	cmd := d.prepareCommand(ctx, containerName, req)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("error creating stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error creating stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error creating stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go d.handleOutput(&wg, ctx, stdout, stderr, output)
	go d.handleInput(&wg, ctx, stdin, input, output)

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		d.killContainer(containerName, output)
		wg.Wait()
		return ctx.Err()
	case err := <-done:
		wg.Wait()
		return err
	}
}

func (d *DockerRunner) prepareCommand(ctx context.Context, containerName string, req executor.ExecRequest) *exec.Cmd {
	args := d.prepareBaseArgs(containerName)

	switch strings.ToLower(req.Language) {
	case "java":
		return d.prepareJavaCommand(ctx, args, req.Code)
	case "c":
		return d.prepareCCommand(ctx, args, req.Code)
	case "c++", "cpp":
		return d.prepareCppCommand(ctx, args, req.Code)
	case "javascript", "js":
		return d.prepareJavaScriptCommand(ctx, args, req.Code)
	default:
		args = append(args, d.imageName, strings.ToLower(req.Language), "-c", req.Code)
		return exec.CommandContext(ctx, "docker", args...)
	}
}

// Helper methods moved to container package
func (d *DockerRunner) prepareBaseArgs(containerName string) []string {
	args := []string{
		"run",
		"--rm",
		"--name", containerName,
		"-i",
		"--cpus=" + maxCPU,
		"-m", maxMemory,
	}

	args = append(args, d.securityOpts...)

	return args
}

// Add other prepare*Command methods here...

func (d *DockerRunner) handleOutput(wg *sync.WaitGroup, ctx context.Context, stdout, stderr io.ReadCloser, output chan<- string) {
	defer wg.Done()
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	for scanner.Scan() {
		select {
		case output <- scanner.Text():
		case <-ctx.Done():
			return
		}
	}
}

func (d *DockerRunner) handleInput(wg *sync.WaitGroup, ctx context.Context, stdin io.WriteCloser, input <-chan string, output chan<- string) {
	defer wg.Done()
	defer stdin.Close()

	for {
		select {
		case inputLine, ok := <-input:
			if !ok {
				return
			}
			if _, err := fmt.Fprintln(stdin, inputLine); err != nil {
				output <- "Error writing to stdin: " + err.Error()
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func (d *DockerRunner) killContainer(containerName string, output chan<- string) {
	killCmd := exec.Command("docker", "kill", containerName)
	if err := killCmd.Run(); err != nil {
		output <- fmt.Sprintf("Failed to kill container: %v", err)
	} else {
		output <- "Container killed successfully"
	}
}
