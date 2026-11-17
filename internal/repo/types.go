package repo

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/client"
)

type authMethod string

const (
	authNone  authMethod = ""
	authSSH   authMethod = "ssh"
	authHTTPS authMethod = "https"
)

type GitConfig struct {
	Repo       string     `yaml:"repo" json:"repo"`
	Branch     string     `yaml:"branch" json:"branch"`
	RepoPath   string     `yaml:"repoPath" json:"repoPath"`
	AuthType   authMethod `yaml:"authType" json:"authType"`
	AuthSecret string     `yaml:"authSecret,omitempty" json:"authSecret,omitempty"`
	DiskPath   string
	auth       client.Option
	git        *git.Repository
}
