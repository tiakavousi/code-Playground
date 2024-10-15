package executor

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func ExecuteInDocker(req ExecRequest) (string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return "", err
	}

	// Define the container configuration
	containerConfig := &container.Config{
		Image: "python:3.9-alpine", // Change based on language
		Cmd:   []string{"python3", "-c", req.Code},
		Tty:   false,
	}

	// Create the container
	resp, err := cli.ContainerCreate(ctx, containerConfig, nil, nil, nil, "")
	if err != nil {
		return "", err
	}

	// Start the container
	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}

	// Get the logs
	result, err := cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{})
	if err != nil {
		return "", err
	}
	defer result.Close()

	// Remove the container
	defer cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{})

	// Read and return the container's output
	output, err := io.ReadAll(result)
	if err != nil {
		return "", err
	}

	return string(output), nil
}
