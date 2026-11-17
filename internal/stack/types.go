package stack

import (
	"github.com/EErikas/CompoCD/internal/docker"
	"github.com/EErikas/CompoCD/internal/repo"
)

const LockFileName = "compocd.lock"

type Project struct {
	Name       string `yaml:"name"`
	Dir        string `yaml:"dir"`
	SopsAgeKey string `yaml:"sops_age_key,omitempty"`
}

type NodeConfig struct {
	Project   Project        `yaml:"project"`
	GitConfig repo.GitConfig `yaml:"git"`
}

type Node struct {
	Config   NodeConfig
	LockFile string
	docker   docker.Docker
	stacks   map[string][]string
}

type LockConfig struct {
	Stacks []string `yaml:"stacks"`
}
