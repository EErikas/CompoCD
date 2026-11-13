package stack

import "github.com/EErikas/CompoCD/internal/compose"

type Project struct {
	Name string `yaml:"name"`
	Dir  string `yaml:"dir"`
}

type NodeConfig struct {
	Project Project `yaml:"project"`
}

type Node struct {
	Config  NodeConfig
	compose compose.Compose
	stacks  map[string][]string
}
