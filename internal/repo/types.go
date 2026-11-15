package repo

import (
	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing/client"
)

type AuthMethod string

const (
	AuthNone  AuthMethod = ""
	AuthSSH   AuthMethod = "ssh"
	AuthHTTPS AuthMethod = "https"
)

type GitConfig struct {
	Repo       string     `yaml:"repo" json:"repo"`
	URL        string     `yaml:"url" json:"url"`
	Branch     string     `yaml:"branch" json:"branch"`
	Path       string     `yaml:"path" json:"path"`
	AuthType   AuthMethod `yaml:"authType" json:"authType"`
	AuthSecret string     `yaml:"authSecret,omitempty" json:"authSecret,omitempty"`
	auth       client.Option
	git        *git.Repository
}
