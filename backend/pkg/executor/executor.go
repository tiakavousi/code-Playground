package executor

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// ExecRequest defines the input for code execution
type ExecRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// ExecuteCode runs the submitted code and returns the output or an error
func ExecuteCode(req ExecRequest) (string, error) {
	var cmd *exec.Cmd
	var output []byte
	var err error

	switch req.Language {
	case "python":
		// Run Python code directly
		cmd = exec.Command("python3", "-c", req.Code)
	case "javascript":
		// Run JavaScript using Node.js
		cmd = exec.Command("node", "-e", req.Code)
	case "bash":
		// Run Bash script
		cmd = exec.Command("bash", "-c", req.Code)
	case "java":
		// Handle Java by compiling and running
		file := "Main.java"
		err = ioutil.WriteFile(file, []byte(req.Code), 0644)
		if err != nil {
			return "", err
		}
		defer os.Remove(file)

		// Compile Java code
		cmd = exec.Command("javac", file)
		if err = cmd.Run(); err != nil {
			return "", err
		}

		// Run the compiled Java program
		cmd = exec.Command("java", "Main")

	case "c":
		// Handle C by compiling and running
		sourceFile := "main.c"
		err = ioutil.WriteFile(sourceFile, []byte(req.Code), 0644)
		if err != nil {
			return "", err
		}
		defer os.Remove(sourceFile)

		// Compile C code
		binaryFile := filepath.Join(os.TempDir(), "main_c")
		cmd = exec.Command("gcc", sourceFile, "-o", binaryFile)
		compileOutput, err := cmd.CombinedOutput()
		if err != nil {
			return string(compileOutput), fmt.Errorf("compilation error: %v", err)
		}

		// Run the compiled binary
		cmd = exec.Command(binaryFile)
		defer os.Remove(binaryFile) // Clean up binary after execution

	case "cpp":
		// Handle C++ by compiling and running
		sourceFile := "main.cpp"
		err = ioutil.WriteFile(sourceFile, []byte(req.Code), 0644)
		if err != nil {
			return "", err
		}
		defer os.Remove(sourceFile)

		// Compile C++ code
		binaryFile := filepath.Join(os.TempDir(), "main_cpp")
		cmd = exec.Command("g++", sourceFile, "-o", binaryFile)
		compileOutput, err := cmd.CombinedOutput()
		if err != nil {
			return string(compileOutput), fmt.Errorf("compilation error: %v", err)
		}

		// Run the compiled binary
		cmd = exec.Command(binaryFile)
		defer os.Remove(binaryFile) // Clean up binary after execution

	default:
		return "", fmt.Errorf("unsupported language: %s", req.Language)
	}

	// Capture the output from running the code
	output, err = cmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	return string(output), nil
}
