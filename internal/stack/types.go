package stack

import "github.com/EErikas/CompoCD/internal/docker"

const LockFileName = "compocd.lock"

type Project struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
}

type NodeConfig struct {
	Project Project `yaml:"project"`
}

type Node struct {
	Config   NodeConfig
	LockFile string
	docker   docker.Docker
	stacks   map[string][]string
}

type LockConfig struct {
	Config NodeConfig `yaml:"config"`
	Stacks []string   `yaml:"stacks"`
}
