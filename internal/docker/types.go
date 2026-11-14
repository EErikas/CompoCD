package docker

import (
	"context"

	"github.com/docker/compose/v5/pkg/api"
	"github.com/moby/moby/client"
)

type Docker struct {
	ctx          context.Context
	client       client.APIClient
	compose      api.Compose
	ComposeFiles []string
}

type PruneReport struct {
	TotalSpaceReclaimed uint64
	ContainersPruned    int
	ImagesPruned        int
	VolumesPruned       int
	NetworksPruned      int
}
