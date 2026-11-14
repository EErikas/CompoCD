package stack

import "github.com/EErikas/CompoCD/internal/docker"

type Project struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
}

type NodeConfig struct {
	Project Project `yaml:"project"`
}

type Node struct {
	Config NodeConfig
	docker docker.Docker
	stacks map[string][]string
}
