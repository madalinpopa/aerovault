package dockerbackup

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// The Inspector defines methods to inspect the state and details of a container using its container ID.
type Inspector interface {
	ContainerInspect(ctx context.Context, containerID string) (types.ContainerJSON, error)
}

// The Creator defines an interface for creating Docker containers with the specified configurations and context.
type Creator interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
}

// Starter defines methods to start a container with given options in a specific context.
type Starter interface {
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
}

// APIClient defines an interface for container operations including inspect, create, and start.
type APIClient interface {
	Inspector
	Creator
	Starter
}
