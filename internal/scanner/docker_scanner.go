package scanner

import (
	"context"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mhshahzad/portman/internal/ports"
)

// DockerScanner implements PortScanner using the Docker SDK.
type DockerScanner struct {
	client *client.Client
}

func (s *DockerScanner) Name() string { return "docker" }

// NewDockerScanner initializes a new Docker client.
func NewDockerScanner() (*DockerScanner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerScanner{client: cli}, nil
}

func (s *DockerScanner) Scan() ([]ports.Port, error) {
	ctx := context.Background()
	containers, err := s.client.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var result []ports.Port
	for _, c := range containers {
		serviceName := s.resolveServiceName(c)

		for _, p := range c.Ports {
			if p.PublicPort == 0 {
				continue
			}

			result = append(result, ports.Port{
				Number:      int(p.PublicPort),
				Protocol:    p.Type,
				Service:     serviceName,
				ContainerID: c.ID[:12],
				Source:      "docker",
			})
		}
	}

	return result, nil
}

func (s *DockerScanner) resolveServiceName(c container.Summary) string {
	// 1. Docker Compose Label
	if service, ok := c.Labels["com.docker.compose.service"]; ok {
		return service
	}

	// 2. Container Name (strip leading slash)
	if len(c.Names) > 0 {
		return strings.TrimPrefix(c.Names[0], "/")
	}

	// 3. Image Name
	return c.Image

	// ID fallback is handled by return or last resort elsewhere
}
