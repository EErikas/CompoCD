package stack

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"

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
	if cfg.Project.SopsAgeKey != "" {
		if err := os.Setenv("SOPS_AGE_KEY", cfg.Project.SopsAgeKey); err != nil {
			return fmt.Errorf("set SOPS_AGE_KEY: %w", err)
		}
	}
	s.Config = cfg
	s.LockFile = filepath.Join(s.Config.Project.Dir, LockFileName)
	s.stacks = make(map[string][]string)
	s.docker = docker.SetupDocker()
	s.Config.GitConfig.DiskPath = s.Config.Project.Dir
	return nil
}

func (s *Node) UpdateRepo() error {
	var err error
	err = s.Config.GitConfig.Init()
	err = s.Config.GitConfig.Fetch()
	return err
}

func (s *Node) StartStack() {
	if err := s.getComposeFiles(); err != nil {
		log.Fatalf("Failed to find Compose files: %v", err)
	}
	lockConfig, exists, err := s.GetLockFile()
	if err != nil {
		log.Fatalf("Failed due to error %e", err)
	}
	incoming := slices.Collect(maps.Keys(s.stacks))
	if exists {
		fmt.Println("Reading existing lockfile")
		for _, i := range lockConfig.Stacks {
			if !slices.Contains(incoming, i) {
				fmt.Printf("Removing %s", i)
				s.docker.RemoveService(i)
			}
		}
	}
	s.writeLockFile()

	for stack, files := range s.stacks {
		s.docker.StartService(stack, files)
	}

	report, err := s.docker.Prune()
	if err != nil {
		log.Fatalf("Failed to prune Docker resources: %v", err)
	}
	report.Print()
}

func (s *Node) getComposeFiles() error {
	s.stacks = make(map[string][]string)
	composeFileDir := filepath.Join(
		s.Config.Project.Dir,
		s.Config.GitConfig.RepoPath,
	)
	return filepath.WalkDir(composeFileDir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && (d.Name() == "compose.yaml" || d.Name() == "compose.sops.yaml") {
				key := filepath.Base(filepath.Dir(path))
				s.stacks[key] = append(s.stacks[key], path)
			}
			return nil
		})
}

func (s *Node) writeLockFile() (LockConfig, error) {
	lockConfig := LockConfig{
		Stacks: slices.Collect(maps.Keys(s.stacks)),
	}
	data, marshallErr := yaml.Marshal(&lockConfig)
	if marshallErr == nil {
		writeErr := os.WriteFile(s.LockFile, data, 0644)
		return lockConfig, writeErr
	}
	return lockConfig, marshallErr
}

func (s *Node) GetLockFile() (LockConfig, bool, error) {
	var lockConfig LockConfig

	data, readErr := os.ReadFile(s.LockFile)
	if readErr != nil && errors.Is(readErr, os.ErrNotExist) {
		lockConfig, err := s.writeLockFile()
		return lockConfig, false, err
	} else if readErr == nil {
		parseErr := yaml.Unmarshal(data, &lockConfig)
		return lockConfig, true, parseErr
	} else {
		return lockConfig, false, readErr
	}

}
