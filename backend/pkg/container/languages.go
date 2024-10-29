// language-specific command preparation
package container

import (
	"context"
	"fmt"
	"os/exec"
)

func (d *DockerRunner) prepareJavaCommand(ctx context.Context, baseArgs []string, code string) *exec.Cmd {
	// Add the Docker image name
	args := append(baseArgs, d.imageName)

	// Add the command to run
	command := fmt.Sprintf(`
		cd /sandbox/tmp &&
		echo '%s' > Main.java &&
		javac Main.java &&
		java Main
	`, code)

	args = append(args, "bash", "-c", command)
	return exec.CommandContext(ctx, "docker", args...)
}

func (d *DockerRunner) prepareCCommand(ctx context.Context, baseArgs []string, code string) *exec.Cmd {
	// Add the Docker image name
	args := append(baseArgs, d.imageName)

	// Add the command to run
	command := fmt.Sprintf(`
		cd /sandbox/tmp &&
		echo '%s' > main.c &&
		gcc main.c -o main &&
		./main
	`, code)

	args = append(args, "bash", "-c", command)
	return exec.CommandContext(ctx, "docker", args...)
}

func (d *DockerRunner) prepareCppCommand(ctx context.Context, baseArgs []string, code string) *exec.Cmd {
	// Add the Docker image name
	args := append(baseArgs, d.imageName)

	// Add the command to run
	command := fmt.Sprintf(`
		cd /sandbox/tmp &&
		echo '%s' > main.cpp &&
		g++ main.cpp -o main &&
		./main
	`, code)

	args = append(args, "bash", "-c", command)
	return exec.CommandContext(ctx, "docker", args...)
}

func (d *DockerRunner) prepareJavaScriptCommand(ctx context.Context, baseArgs []string, code string) *exec.Cmd {
	// Add the Docker image name
	args := append(baseArgs, d.imageName)

	// For JavaScript, we use node directly instead of writing to a file
	args = append(args, "node", "-e", code)

	return exec.CommandContext(ctx, "docker", args...)
}
