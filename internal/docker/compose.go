package docker

import (
	"context"
	"log"

	composetypes "github.com/compose-spec/compose-go/v2/types"
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
		ctx:     context.Background(),
		client:  dockerCLI.Client(),
		compose: service,
	}

}

func (d *Docker) loadProject(name string, files []string) (*composetypes.Project, error) {

	return d.compose.LoadProject(d.ctx, api.ProjectLoadOptions{
		ConfigPaths: files,
		ProjectName: name,
	})
}

func (d *Docker) StartService(name string, files []string) {
	// Load the Compose project from a compose file
	project, loadErr := d.loadProject(name, files)
	if loadErr != nil {
		log.Fatalf("Failed to load project: %v", loadErr)
	}

	upOptions := api.UpOptions{
		Create: api.CreateOptions{
			RemoveOrphans: true,
			Inherit:       true,
		},
		Start: api.StartOptions{},
	}

	// Start the services defined in the Compose file
	if err := d.compose.Up(d.ctx, project, upOptions); err != nil {
		log.Fatalf("Failed to start services: %v", err)
	}

	log.Printf("Successfully started project: %s", project.Name)
}

func (d *Docker) RemoveService(name string) {

	downOptions := api.DownOptions{
		RemoveOrphans: true,
		Volumes:       true,
		Images:        api.ImagePruneAll,
	}

	// Start the services defined in the Compose file
	if err := d.compose.Down(d.ctx, name, downOptions); err != nil {
		log.Fatalf("Failed to start services: %v", err)
	}

	log.Printf("Successfully removed project: %s", name)
}
