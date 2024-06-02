package task

import (
	"context"
	"io"
	"log"
	"math"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Docker struct {
	Client *client.Client
	Config Config
}

type DockerResult struct {
	ContainerID string
	Error       error
	Action      string
	Result      string
}

func (d *Docker) Run() DockerResult {
	ctx := context.Background()

	reader, err := d.Client.ImagePull(ctx, d.Config.Image, image.PullOptions{})
	if err != nil {
		log.Printf("Error pulling image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}
	defer func(reader io.ReadCloser) {
		err := reader.Close()
		if err != nil {
			log.Printf("Error closing ImagePull reader: %v\n", err)
		}
	}(reader)

	_, err = io.Copy(os.Stdout, reader)
	if err != nil {
		log.Printf("Error copying from reader from ImagePull: %v\n", err)
		return DockerResult{Error: err}
	}

	rp := container.RestartPolicy{
		Name:              d.Config.RestartPolicy,
		MaximumRetryCount: d.Config.RestartMaximumRetryCount,
	}

	res := container.Resources{
		Memory:   d.Config.Memory,
		NanoCPUs: int64(d.Config.CPU * math.Pow(10, 9)),
	}

	cc := container.Config{
		Image:        d.Config.Image,
		Tty:          false,
		Env:          d.Config.Env,
		ExposedPorts: d.Config.ExposedPorts,
	}

	hc := container.HostConfig{
		RestartPolicy:   rp,
		Resources:       res,
		PublishAllPorts: true,
	}

	resp, err := d.Client.ContainerCreate(ctx, &cc, &hc, nil, nil, d.Config.Name)
	if err != nil {
		log.Printf("Error creating container using image %s: %v\n", d.Config.Image, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		log.Printf("Error starting container %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	out, err := d.Client.ContainerLogs(ctx, resp.ID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		log.Printf("Error getting container logs %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	if err != nil {
		log.Printf("Error calling StdCopy for container logs %s: %v\n", resp.ID, err)
		return DockerResult{Error: err}
	}

	log.Printf("Warnings encountered when starting container %s: %v", resp.ID, resp.Warnings)

	return DockerResult{
		ContainerID: resp.ID,
		Action:      "start",
		Result:      "success",
	}
}

func (d *Docker) Stop(containerID string) DockerResult {
	log.Printf("Trying to stop container %s", containerID)

	ctx := context.Background()

	err := d.Client.ContainerStop(ctx, containerID, container.StopOptions{})
	if err != nil {
		log.Printf("Error stopping container %s: %v\n", containerID, err)
		return DockerResult{Error: err}
	}

	err = d.Client.ContainerRemove(ctx, containerID, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   false,
	})
	if err != nil {
		log.Printf("Error removing container %s: %v\n", containerID, err)
		return DockerResult{Error: err}
	}

	return DockerResult{
		ContainerID: containerID,
		Action:      "stop",
		Result:      "success",
	}
}
