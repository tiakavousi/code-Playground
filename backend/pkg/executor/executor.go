package executor

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	var cmd *exec.Cmd
	var err error

	switch req.Language {
	case "python":
		cmd = exec.CommandContext(ctx, "python3", "-c", req.Code)
	case "javascript":
		cmd = exec.CommandContext(ctx, "node", "-e", req.Code)
	case "bash":
		cmd = exec.CommandContext(ctx, "bash", "-c", req.Code)
	case "java":
		file := "Main.java"
		err = ioutil.WriteFile(file, []byte(req.Code), 0644)
		if err != nil {
			return err
		}
		defer os.Remove(file)

		compileCmd := exec.CommandContext(ctx, "javac", file)
		if err = compileCmd.Run(); err != nil {
			return err
		}

		cmd = exec.CommandContext(ctx, "java", "Main")
	case "c":
		sourceFile := "main.c"
		err = ioutil.WriteFile(sourceFile, []byte(req.Code), 0644)
		if err != nil {
			return err
		}
		defer os.Remove(sourceFile)

		binaryFile := filepath.Join(os.TempDir(), "main_c")
		compileCmd := exec.CommandContext(ctx, "gcc", sourceFile, "-o", binaryFile)
		if compileOutput, err := compileCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("compilation error: %v\n%s", err, compileOutput)
		}

		cmd = exec.CommandContext(ctx, binaryFile)
		defer os.Remove(binaryFile)
	case "cpp":
		sourceFile := "main.cpp"
		err = ioutil.WriteFile(sourceFile, []byte(req.Code), 0644)
		if err != nil {
			return err
		}
		defer os.Remove(sourceFile)

		binaryFile := filepath.Join(os.TempDir(), "main_cpp")
		compileCmd := exec.CommandContext(ctx, "g++", sourceFile, "-o", binaryFile)
		if compileOutput, err := compileCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("compilation error: %v\n%s", err, compileOutput)
		}

		cmd = exec.CommandContext(ctx, binaryFile)
		defer os.Remove(binaryFile)
	default:
		return fmt.Errorf("unsupported language: %s", req.Language)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

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

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		if err := cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %v", err)
		}
		return ctx.Err()
	case err := <-done:
		return err
	}
}

// ExecuteCode runs the submitted code and returns the output or an error (non-interactive version)
func ExecuteCode(req ExecRequest) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
