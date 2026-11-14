package stack

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/EErikas/CompoCD/internal/docker"
	"go.yaml.in/yaml/v3"
)

func (s *Node) ParseStackConfig(path string) error {
	var cfg NodeConfig
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to find file")
		return err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal("Failed to parse config")
		return err
	}
	s.Config = cfg
	s.stacks = make(map[string][]string)
	s.docker = docker.SetupDocker()
	s.getComposeFiles()

	return nil
}

func (s *Node) StartStack() {
	// s.compose.StartServices(s.Config.Project.Name)
	for stack, files := range s.stacks {
		s.docker.StartService(
			fmt.Sprintf("%s-%s", s.Config.Project.Name, stack),
			files,
		)
	}
	report, err := s.docker.Prune()
	if err != nil {
		log.Fatalf("Failed to prune Docker resources: %v", err)
	}
	report.Print()
}

func (s *Node) getComposeFiles() error {
	return filepath.WalkDir(s.Config.Project.Dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && d.Name() == "compose.yaml" {
				key := filepath.Base(filepath.Dir(path))
				s.stacks[key] = append(s.stacks[key], path)
			}
			return nil
		})
}
