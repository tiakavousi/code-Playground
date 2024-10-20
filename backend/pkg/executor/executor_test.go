package executor

import (
	"strings"
	"testing"
)

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
			expected: "Hello, World!\n",
			wantErr:  false,
		},
		{
			name: "JavaScript Hello World",
			req: ExecRequest{
				Language: "javascript",
				Code:     "console.log('Hello, JavaScript World!');",
			},
			expected: "Hello, JavaScript World!\n",
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
}	
`,
			},
			expected: "Hello, Java World!\n",
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
			expected: "Hello, C World!\n",
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
			expected: "Hello, C++ World!\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExecuteCode(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.expected) {
				t.Errorf("ExecuteCode() got = %v, want %v", got, tt.expected)
			}
		})
	}
}
