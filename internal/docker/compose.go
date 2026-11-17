package docker

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path/filepath"

	"github.com/EErikas/CompoCD/internal/sops"
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

func (d *Docker) loadProject(name string, files []string, workingDir string) (*composetypes.Project, error) {

	return d.compose.LoadProject(d.ctx, api.ProjectLoadOptions{
		ConfigPaths: files,
		ProjectName: name,
		WorkingDir:  workingDir,
	})
}

func (d *Docker) StartService(name string, files []string) {
	if err := d.startService(name, files); err != nil {
		log.Fatalf("Failed to start project %s: %v", name, err)
	}
}

func (d *Docker) startService(name string, files []string) (resultErr error) {
	if len(files) == 0 {
		return fmt.Errorf("no Compose files provided")
	}

	composeFiles := make([]string, 0, len(files))
	cleanups := make([]func() error, 0, len(files))
	defer cleanupSOPSFiles(&cleanups, &resultErr)
	for _, file := range files {
		prepared, cleanup, err := sops.PrepareFile(file)
		if err != nil {
			return err
		}
		composeFiles = append(composeFiles, prepared)
		cleanups = append(cleanups, cleanup)
	}

	workingDir, err := filepath.Abs(filepath.Dir(files[0]))
	if err != nil {
		return fmt.Errorf("resolve Compose project directory: %w", err)
	}

	// Load the Compose project from a compose file
	project, loadErr := d.loadProject(name, composeFiles, workingDir)
	if loadErr != nil {
		return fmt.Errorf("load project: %w", loadErr)
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
		return fmt.Errorf("start services: %w", err)
	}

	log.Printf("Successfully started project: %s", project.Name)
	return nil
}

func cleanupSOPSFiles(cleanups *[]func() error, resultErr *error) {
	var cleanupErr error
	for _, cleanup := range *cleanups {
		if err := cleanup(); err != nil {
			cleanupErr = errors.Join(cleanupErr, fmt.Errorf("remove temporary Compose file: %w", err))
		}
	}
	*resultErr = errors.Join(*resultErr, cleanupErr)
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
