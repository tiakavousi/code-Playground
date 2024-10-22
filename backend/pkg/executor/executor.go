package executor

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

// ExecRequest defines the input for code execution
type ExecRequest struct {
	Language string `json:"language" binding:"required"`
	Code     string `json:"code" binding:"required"`
}

// CodeRunner interface defines methods that must be implemented by any code execution backend
type CodeRunner interface {
	RunInteractive(ctx context.Context, req ExecRequest, input <-chan string, output chan<- string) error
}

// Service represents the code execution service
type Service struct {
	runner CodeRunner
}

// NewService creates a new executor service with the specified runner
func NewService(runner CodeRunner) *Service {
	return &Service{
		runner: runner,
	}
}

// ExecuteInteractive runs code with interactive I/O
func (s *Service) ExecuteInteractive(ctx context.Context, req ExecRequest, input <-chan string, output chan<- string) error {
	// Validate request
	if err := validateRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Create execution context
	execCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Run the code
	err := s.runner.RunInteractive(execCtx, req, input, output)
	if err != nil {
		log.Printf("Execution error for language %s: %v", req.Language, err)
		if execCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("execution timed out after 30 seconds")
		}
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}

// Execute runs code and returns the output (non-interactive)
func (s *Service) Execute(req ExecRequest) (string, error) {
	// Validate request
	if err := validateRequest(req); err != nil {
		return "", fmt.Errorf("invalid request: %w", err)
	}

	// Create execution context
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create channels for communication
	output := make(chan string, 5)
	input := make(chan string, 5)
	defer close(input)

	// Execute the code
	errCh := make(chan error, 1)
	go func() {
		start := time.Now()
		err := s.ExecuteInteractive(ctx, req, input, output)
		log.Printf("Execution took %v", time.Since(start))
		errCh <- err
	}()

	// Collect output
	var result strings.Builder
	for {
		select {
		case line := <-output:
			result.WriteString(line + "\n")
		case err := <-errCh:
			return result.String(), err
		case <-ctx.Done():
			return result.String(), fmt.Errorf("execution timed out after %v: %w", 60*time.Second, ctx.Err())
		}
	}
}

// validateRequest checks if the request is valid
func validateRequest(req ExecRequest) error {
	if strings.TrimSpace(req.Language) == "" {
		return fmt.Errorf("language cannot be empty")
	}
	if strings.TrimSpace(req.Code) == "" {
		return fmt.Errorf("code cannot be empty")
	}
	return nil
}
