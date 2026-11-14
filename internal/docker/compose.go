package docker

import (
	"context"
	"log"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

func SetupDocker() Docker {
	dockerCLI, err := command.NewDockerCli()
	if err != nil {
		log.Fatalf("Failed to create docker CLI: %v", err)
	}
	err = dockerCLI.Initialize(&flags.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to initialize docker CLI: %v", err)
	}
	// Create a new Compose service instance
	service, err := compose.NewComposeService(dockerCLI)
	if err != nil {
		log.Fatalf("Failed to create compose service: %v", err)
	}

	return Docker{
		ctx:          context.Background(),
		client:       dockerCLI.Client(),
		compose:      service,
		ComposeFiles: []string{},
	}

}

func (c *Docker) StartServices(name string) {
	ctx := context.Background()

	// Load the Compose project from a compose file
	project, err := c.compose.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: c.ComposeFiles,
		ProjectName: name,
	})
	if err != nil {
		log.Fatalf("Failed to load project: %v", err)
	}

	// Start the services defined in the Compose file
	err = c.compose.Up(ctx, project, api.UpOptions{
		Create: api.CreateOptions{
			RemoveOrphans: true,
			Inherit:       true,
		},
		Start: api.StartOptions{},
	})
	if err != nil {
		log.Fatalf("Failed to start services: %v", err)
	}

	log.Printf("Successfully started project: %s", project.Name)
}

func (c *Docker) StartService(name string, files []string) {
	// Load the Compose project from a compose file
	project, err := c.compose.LoadProject(c.ctx, api.ProjectLoadOptions{
		ConfigPaths: files,
		ProjectName: name,
	})
	if err != nil {
		log.Fatalf("Failed to load project: %v", err)
	}

	// Start the services defined in the Compose file
	err = c.compose.Up(c.ctx, project, api.UpOptions{
		Create: api.CreateOptions{
			RemoveOrphans: true,
			Inherit:       true,
		},
		Start: api.StartOptions{},
	})
	if err != nil {
		log.Fatalf("Failed to start services: %v", err)
	}

	log.Printf("Successfully started project: %s", project.Name)
}
