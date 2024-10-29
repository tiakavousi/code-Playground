package executor

import (
    "context"
    "fmt"
    "strings"
    "testing"
    "time"
)

// MockRunner implements CodeRunner interface for testing
type MockRunner struct {
    responses map[string]mockResponse
}

type mockResponse struct {
    output string
    err    error
}

func NewMockRunner() *MockRunner {
    return &MockRunner{
        responses: map[string]mockResponse{
            "python3": {
                output: "Hello, World!",
                err:    nil,
            },
            "javascript": {
                output: "Hello, JavaScript World!",
                err:    nil,
            },
            "java": {
                output: "Hello, Java World!",
                err:    nil,
            },
            "c": {
                output: "Hello, C World!",
                err:    nil,
            },
            "c++": {
                output: "Hello, C++ World!",
                err:    nil,
            },
            "cpp": {
                output: "Hello, C++ World!",
                err:    nil,
            },
            "invalid": {
                output: "",
                err:    fmt.Errorf("unsupported language: invalid"),
            },
            "python-infinite": {
                output: "",
                err:    fmt.Errorf("execution timed out"),
            },
        },
    }
}

func (m *MockRunner) RunInteractive(ctx context.Context, req ExecRequest, input <-chan string, output chan<- string) error {
    // Special case for infinite loop
    if req.Language == "python" && strings.Contains(req.Code, "while True: pass") {
        return m.responses["python-infinite"].err
    }

    // Get the response based on language
    response, exists := m.responses[req.Language]
    if !exists {
        return fmt.Errorf("unexpected language: %s", req.Language)
    }

    if response.err != nil {
        return response.err
    }

    output <- response.output
    return nil
}

func TestExecuteCode(t *testing.T) {
    tests := []struct {
        name     string
        req      ExecRequest
        expected string
        wantErr  bool
    }{
        {
            name: "Python Hello World",
            req: ExecRequest{
                Language: "python3",
                Code:     "print('Hello, World!')",
            },
            expected: "Hello, World!",
            wantErr:  false,
        },
        {
            name: "JavaScript Hello World",
            req: ExecRequest{
                Language: "javascript",
                Code:     "console.log('Hello, JavaScript World!');",
            },
            expected: "Hello, JavaScript World!",
            wantErr:  false,
        },
        {
            name: "Invalid Language",
            req: ExecRequest{
                Language: "invalid",
                Code:     "print('This should not run')",
            },
            expected: "",
            wantErr:  true,
        },
        {
            name: "Python Infinite Loop",
            req: ExecRequest{
                Language: "python",
                Code:     "while True: pass",
            },
            expected: "",
            wantErr:  true,
        },
        {
            name: "Java Hello World",
            req: ExecRequest{
                Language: "java",
                Code: `
public class Main {
    public static void main(String args[]) {
        System.out.println("Hello, Java World!");
    }
}`,
            },
            expected: "Hello, Java World!",
            wantErr:  false,
        },
        {
            name: "C Hello World",
            req: ExecRequest{
                Language: "c",
                Code: `
#include <stdio.h>

int main() {
    printf("Hello, C World!\n");
    return 0;
}`,
            },
            expected: "Hello, C World!",
            wantErr:  false,
        },
        {
            name: "C++ Hello World",
            req: ExecRequest{
                Language: "c++",
                Code: `
#include <iostream>

int main() {
    std::cout << "Hello, C++ World!" << std::endl;
    return 0;
}`,
            },
            expected: "Hello, C++ World!",
            wantErr:  false,
        },
    }

    // Create mock runner and service
    mockRunner := NewMockRunner()
    service := NewService(mockRunner)

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create channels for communication
            input := make(chan string)
            output := make(chan string)
            
            // Create context with timeout
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()

            // Create error channel
            errCh := make(chan error)
            
            // Run ExecuteInteractive in goroutine
            go func() {
                errCh <- service.ExecuteInteractive(ctx, tt.req, input, output)
                close(input)
            }()

            // Collect output
            var got string
            select {
            case line := <-output:
                got = line
            case err := <-errCh:
                if (err != nil) != tt.wantErr {
                    t.Errorf("ExecuteInteractive() error = %v, wantErr %v", err, tt.wantErr)
                }
                return
            case <-ctx.Done():
                t.Error("ExecuteInteractive() timed out")
                return
            }

            if !tt.wantErr && !strings.Contains(got, strings.TrimSpace(tt.expected)) {
                t.Errorf("ExecuteInteractive() got = %v, want %v", got, tt.expected)
            }
        })
    }
}

// TestExecuteInteractive tests the interactive execution functionality
func TestExecuteInteractive(t *testing.T) {
    mockRunner := NewMockRunner()
    service := NewService(mockRunner)

    tests := []struct {
        name     string
        req      ExecRequest
        input    []string
        expected []string
        wantErr  bool
    }{
        {
            name: "Simple Python",
            req: ExecRequest{
                Language: "python3",
                Code:     "print('Hello, World!')",
            },
            input:    []string{},
            expected: []string{"Hello, World!"},
            wantErr:  false,
        },
        {
            name: "Invalid Language",
            req: ExecRequest{
                Language: "invalid",
                Code:     "test",
            },
            input:    []string{},
            expected: []string{},
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
            defer cancel()

            input := make(chan string, len(tt.input))
            output := make(chan string, len(tt.expected))

            // Send test input
            go func() {
                for _, in := range tt.input {
                    input <- in
                }
                close(input)
            }()

            // Run interactive execution
            err := service.ExecuteInteractive(ctx, tt.req, input, output)
            if (err != nil) != tt.wantErr {
                t.Errorf("ExecuteInteractive() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            // Check output
            if !tt.wantErr {
                for _, expected := range tt.expected {
                    select {
                    case got := <-output:
                        if !strings.Contains(got, strings.TrimSpace(expected)) {
                            t.Errorf("ExecuteInteractive() output = %v, want %v", got, expected)
                        }
                    case <-time.After(time.Second):
                        t.Error("ExecuteInteractive() timed out waiting for output")
                    }
                }
            }
        })
    }
}

// TestValidateRequest tests the request validation function
func TestValidateRequest(t *testing.T) {
    tests := []struct {
        name    string
        req     ExecRequest
        wantErr bool
    }{
        {
            name: "Valid Request",
            req: ExecRequest{
                Language: "python",
                Code:     "print('test')",
            },
            wantErr: false,
        },
        {
            name: "Empty Language",
            req: ExecRequest{
                Language: "",
                Code:     "print('test')",
            },
            wantErr: true,
        },
        {
            name: "Empty Code",
            req: ExecRequest{
                Language: "python",
                Code:     "",
            },
            wantErr: true,
        },
        {
            name: "Whitespace Only Code",
            req: ExecRequest{
                Language: "python",
                Code:     "   ",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateRequest(tt.req)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}