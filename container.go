package main

import (
	"context"
	"io"
	"os"

	dockerContainer "github.com/docker/docker/api/types/container"
	dockerImage "github.com/docker/docker/api/types/image"
	dockerClient "github.com/docker/docker/client"
)

const (
	POSTGRES_IMAGE = "docker.io/library/alpine"
)

type (
	Container struct {
		Client *dockerClient.Client
		ID     string
		State  string
	}
)

func startContainer() (Container, error) {
	var c Container
	var err error

	c.Client, err = dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		return c, err
	}

	// Pull the Docker image (optional, if you already have the image, you can skip this)
	out, err := c.Client.ImagePull(context.Background(), POSTGRES_IMAGE, dockerImage.PullOptions{})
	if err != nil {
		return c, err
	}
	defer out.Close()
	io.Copy(os.Stdout, out)

	// Define container configuration
	config := &dockerContainer.Config{
		Image: "alpine",
		Cmd:   []string{"echo", "Hello, World!"},
	}
	resp, err := c.Client.ContainerCreate(context.Background(), config, nil, nil, nil, "")
	if err != nil {
		return c, err
	}

	// Start the container
	if err := c.Client.ContainerStart(context.Background(), resp.ID, dockerContainer.StartOptions{}); err != nil {
		return c, err
	}

	// Inspect the container to get its details
	containerJSON, err := c.Client.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		return c, err
	}

	// Print container details
	c.ID = containerJSON.ID
	c.State = containerJSON.State.Status

	return c, nil
}

func (c *Container) shutdownContainer() error {
	err := c.Client.ContainerStop(context.Background(), c.ID, dockerContainer.StopOptions{
		Timeout: nil,
	})
	if err != nil {
		return err
	}

	containerJSON, err := c.Client.ContainerInspect(context.Background(), c.ID)
	if err != nil {
		return err
	}
	c.State = containerJSON.State.Status
	return nil
}
