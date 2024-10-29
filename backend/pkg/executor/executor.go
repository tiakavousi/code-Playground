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
func (s *Service) ExecuteInteractive(
	ctx context.Context, req ExecRequest, 
	input <-chan string, output chan<- string) error {
	// Validate request
	if err := validateRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Create execution context
	execCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Run the code
	err := s.runner.RunInteractive(execCtx, req, input, output)
	if err != nil {
		log.Printf("Execution error for language %s: %v", req.Language, err)
		if execCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("execution timed out after 10 seconds")
		}
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
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
