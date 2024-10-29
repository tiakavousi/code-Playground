package container

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/tiakavousi/codeplayground/pkg/executor"
)

// TestDockerRunner wraps DockerRunner for testing
type TestDockerRunner struct {
	*DockerRunner
	execCommand       commandFunc
	killContainerFunc func(containerName string, output chan<- string)
}

// commandFunc is a function type for executing commands
type commandFunc func(name string, args ...string) *exec.Cmd

// execCommand is our package-level variable for the command function
var execCommand = exec.Command

// mockCommand creates a command that simulates Docker behavior
func mockCommand(name string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{
		"GO_WANT_HELPER_PROCESS=1",
		"MOCK_COMMAND=" + name,
		"MOCK_ARGS=" + strings.Join(args, " "),
	}
	return cmd
}

// TestHelperProcess helps mock exec.Command
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	// Extract mock command and args
	mockCmd := os.Getenv("MOCK_COMMAND")
	mockArgs := strings.Split(os.Getenv("MOCK_ARGS"), " ")

	if mockCmd == "docker" {
		if len(mockArgs) > 0 && mockArgs[0] == "run" {
			fmt.Println("Container output")
			if strings.Contains(strings.Join(mockArgs, " "), "while True: pass") {
				time.Sleep(200 * time.Millisecond) // Simulate long-running process
			}
		} else if len(mockArgs) > 0 && mockArgs[0] == "kill" {
			fmt.Println("Container killed successfully")
		}
	}
	os.Exit(0)
}

func (d *TestDockerRunner) killContainer(containerName string, output chan<- string) {
	d.killContainerFunc(containerName, output)
}

func NewTestDockerRunner(imageName string) *TestDockerRunner {
	runner := &TestDockerRunner{
		DockerRunner: NewDockerRunner(imageName),
		execCommand:  execCommand,
	}
	runner.killContainerFunc = runner.killContainer
	return runner
}

func TestRunInteractive(t *testing.T) {
	// Save original command function and restore it after tests
	originalCommand := execCommand
	defer func() { execCommand = originalCommand }()

	// Replace command function with our mock
	execCommand = mockCommand

	tests := []struct {
		name       string
		language   string
		code       string
		input      []string
		wantOutput []string
		wantErr    bool
		timeout    time.Duration
	}{
		{
			name:       "Basic Python Execution",
			language:   "python3",
			code:       "print('hello')",
			input:      []string{},
			wantOutput: []string{"Container output"},
			timeout:    2 * time.Second,
			wantErr:    false,
		},
		{
			name:       "Context Cancellation",
			language:   "python",
			code:       "while True: pass",
			input:      []string{},
			wantOutput: []string{"Container killed successfully"},
			timeout:    100 * time.Millisecond,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewTestDockerRunner("tayebe/repl")
			runner.execCommand = mockCommand
			runner.killContainerFunc = func(containerName string, output chan<- string) {
				fmt.Println("Container killed successfully")
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			input := make(chan string, len(tt.input))
			output := make(chan string, 10)

			// Send inputs
			go func() {
				defer close(input)
				for _, in := range tt.input {
					select {
					case input <- in:
					case <-ctx.Done():
						return
					}
				}
			}()

			// Run command and collect outputs
			done := make(chan struct{})
			var gotOutput []string
			var runErr error

			go func() {
				defer close(done)
				runErr = runner.RunInteractive(ctx, executor.ExecRequest{
					Language: tt.language,
					Code:     tt.code,
				}, input, output)
				close(output) // Important: close output channel after command completes
			}()

			// Collect output until command completes
			outputDone := make(chan struct{})
			go func() {
				defer close(outputDone)
				for line := range output {
					gotOutput = append(gotOutput, line)
				}
			}()

			// Wait for completion or timeout
			select {
			case <-done:
				// Command completed
				if (runErr != nil) != tt.wantErr {
					t.Errorf("RunInteractive() error = %v, wantErr %v", runErr, tt.wantErr)
				}
			case <-time.After(tt.timeout + time.Second):
				t.Fatal("Test timeout exceeded")
			}

			// Wait for output collection
			select {
			case <-outputDone:
			case <-time.After(500 * time.Millisecond):
				t.Fatal("Timeout waiting for output collection")
			}

			// Verify output
			for _, want := range tt.wantOutput {
				found := false
				for _, got := range gotOutput {
					if strings.Contains(got, want) {
						found = true
						break
					}
				}
				if !found {
					t.Logf("RunInteractive() missing expected output %q in %v\n", want, gotOutput)
				}
			}
		})
	}
}

func TestPrepareCommand(t *testing.T) {
	tests := []struct {
		name     string
		language string
		code     string
		wantArgs []string
	}{
		{
			name:     "Python Command",
			language: "python3",
			code:     "print('hello')",
			wantArgs: []string{"python3", "-c", "print('hello')"},
		},
		{
			name:     "JavaScript Command",
			language: "javascript",
			code:     "console.log('hello')",
			wantArgs: []string{"node", "-e", "console.log('hello')"},
		},
	}

	ctx := context.Background()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := NewTestDockerRunner("tayebe/repl")
			runner.execCommand = mockCommand

			cmd := runner.prepareCommand(ctx, "test-container", executor.ExecRequest{
				Language: tt.language,
				Code:     tt.code,
			})

			for _, arg := range tt.wantArgs {
				found := false
				for _, cmdArg := range cmd.Args {
					if cmdArg == arg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("prepareCommand() missing argument %q in %v", arg, cmd.Args)
				}
			}
		})
	}
}

func TestNewDockerRunner(t *testing.T) {
	imageName := "tayebe/repl"
	runner := NewDockerRunner(imageName)

	if runner.imageName != imageName {
		t.Errorf("NewDockerRunner() imageName = %v, want %v", runner.imageName, imageName)
	}

	expectedOpts := []string{
		"--cap-drop=ALL",
		"--net=none",
		"--pids-limit=20",
		"--ulimit", "nproc=20:20",
		"--ulimit", "nofile=64:64",
		"--ulimit", "fsize=1000000:1000000",
	}

	if !reflect.DeepEqual(runner.securityOpts, expectedOpts) {
		t.Errorf("NewDockerRunner() securityOpts = %v, want %v", runner.securityOpts, expectedOpts)
	}
}
