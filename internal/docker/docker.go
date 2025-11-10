package docker

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

const (
	labelExclude    = "docker-status.exclude"
	labelCustomName = "docker-status.name"
)

type ContainerInfo struct {
	ID        string
	Name      string
	Status    string
	State     string
	Image     string
	CreatedAt time.Time
}

func New() (*Docker, error) {
	cli, err := client.NewClientWithOpts(
		client.WithHost("unix:///var/run/docker.sock"),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("docker client init err: %s", err.Error())
	}

	_, err = cli.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("docker client init err: %s", err.Error())
	}

	return &Docker{
		cli: cli,
	}, nil
}

type Docker struct {
	cli *client.Client
}

func (d *Docker) Close() error {
	return d.cli.Close()
}

func (d *Docker) GetContainers(ctx context.Context) ([]ContainerInfo, error) {

	containers, err := d.cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var result []ContainerInfo

	for _, c := range containers {
		if exclude, ok := c.Labels[labelExclude]; ok {
			if excludeBool, err := strconv.ParseBool(exclude); err == nil && excludeBool {
				continue
			}
		}

		name := c.Names[0]
		if name[0] == '/' {
			name = name[1:]
		}
		if customName, ok := c.Labels[labelCustomName]; ok && customName != "" {
			name = customName
		}

		imageName := c.Image
		if idx := strings.LastIndex(imageName, "/"); idx != -1 {
			imageName = imageName[idx+1:]
		}

		containerInfo := ContainerInfo{
			ID:        c.ID[:12],
			Name:      name,
			Status:    strings.ToUpper(c.State),
			State:     strings.ToLower(c.State),
			Image:     imageName,
			CreatedAt: time.Unix(c.Created, 0),
		}

		result = append(result, containerInfo)
	}

	return result, nil
}
