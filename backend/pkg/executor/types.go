package executor

import "errors"

var (
	// ErrInvalidLanguage is returned when the requested language is not supported
	ErrInvalidLanguage = errors.New("invalid or unsupported language")

	// ErrEmptyCode is returned when the submitted code is empty
	ErrEmptyCode = errors.New("code cannot be empty")

	// ErrExecutionTimeout is returned when code execution takes too long
	ErrExecutionTimeout = errors.New("execution timed out")
)

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
	Error    string `json:"error,omitempty"`
}
