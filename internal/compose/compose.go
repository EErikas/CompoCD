package compose

import (
	"context"
	"log"

	"github.com/docker/cli/cli/command"
	"github.com/docker/cli/cli/flags"
	"github.com/docker/compose/v5/pkg/api"
	"github.com/docker/compose/v5/pkg/compose"
)

type Compose struct {
	ctx          context.Context
	service      api.Compose
	ComposeFiles []string
}

func SetupCompose() Compose {
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

	return Compose{
		ctx:          context.Background(),
		service:      service,
		ComposeFiles: []string{},
	}

}

func (c *Compose) StartServices(name string) {
	ctx := context.Background()

	// Load the Compose project from a compose file
	project, err := c.service.LoadProject(ctx, api.ProjectLoadOptions{
		ConfigPaths: c.ComposeFiles,
		ProjectName: name,
	})
	if err != nil {
		log.Fatalf("Failed to load project: %v", err)
	}

	// Start the services defined in the Compose file
	err = c.service.Up(ctx, project, api.UpOptions{
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

func (c *Compose) StartService(name string, files []string) {
	// Load the Compose project from a compose file
	project, err := c.service.LoadProject(c.ctx, api.ProjectLoadOptions{
		ConfigPaths: files,
		ProjectName: name,
	})
	if err != nil {
		log.Fatalf("Failed to load project: %v", err)
	}

	// Start the services defined in the Compose file
	err = c.service.Up(c.ctx, project, api.UpOptions{
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
