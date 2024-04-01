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
	container struct {
		client *dockerClient.Client
		repo   *repo
		id     string
		state  string
	}
)

func initContainer() (container, error) {
	var c container
	var err error

	c.client, err = dockerClient.NewClientWithOpts(dockerClient.FromEnv, dockerClient.WithAPIVersionNegotiation())
	if err != nil {
		return c, err
	}

	// Pull the Docker image (optional, if you already have the image, you can skip this)
	out, err := c.client.ImagePull(context.Background(), POSTGRES_IMAGE, dockerImage.PullOptions{})
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
	resp, err := c.client.ContainerCreate(context.Background(), config, nil, nil, nil, "")
	if err != nil {
		return c, err
	}

	// Start the container
	if err := c.client.ContainerStart(context.Background(), resp.ID, dockerContainer.StartOptions{}); err != nil {
		return c, err
	}

	// Inspect the container to get its details
	containerJSON, err := c.client.ContainerInspect(context.Background(), resp.ID)
	if err != nil {
		return c, err
	}

	// Print container details
	c.id = containerJSON.ID
	c.state = containerJSON.State.Status

	c.repo, err = initRepo()

	return c, nil
}

func (c *container) shutdownContainer() error {
	err := c.client.ContainerStop(context.Background(), c.id, dockerContainer.StopOptions{
		Timeout: nil,
	})
	if err != nil {
		return err
	}

	err = c.repo.db.Close()
	if err != nil {
		return err
	}

	containerJSON, err := c.client.ContainerInspect(context.Background(), c.id)
	if err != nil {
		return err
	}

	c.state = containerJSON.State.Status
	return nil
}
