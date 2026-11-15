package repo

import (
	"errors"
	"fmt"
	"os"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	httptransport "github.com/go-git/go-git/v6/plumbing/transport/http"
	sshtransport "github.com/go-git/go-git/v6/plumbing/transport/ssh"
	"go.yaml.in/yaml/v4"
)

func SetupGit(configFile string) (GitConfig, error) {
	var cfg GitConfig
	data, readErr := os.ReadFile(configFile)
	if readErr != nil {
		return cfg, readErr
	}
	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (cfg *GitConfig) configureAuth() error {
	cfg.auth = nil
	switch cfg.AuthType {
	case AuthNone:
		return nil
	case AuthHTTPS:
		if cfg.AuthSecret == "" {
			return fmt.Errorf("HTTPS authentication requires an access token")
		}
		cfg.auth = client.WithHTTPAuth(&httptransport.BasicAuth{
			Username: "oauth2",
			Password: cfg.AuthSecret,
		})
		return nil
	case AuthSSH:
		if cfg.AuthSecret == "" {
			return fmt.Errorf("SSH authentication requires a private key file path")
		}
		auth, err := sshtransport.NewPublicKeysFromFile("git", cfg.AuthSecret, "")
		if err != nil {
			return fmt.Errorf("load SSH private key: %w", err)
		}
		cfg.auth = client.WithSSHAuth(auth)
		return nil
	default:
		return fmt.Errorf("unsupported authentication type %q", cfg.AuthType)
	}
}

func (cfg *GitConfig) Init(projectPath string) error {
	if err := cfg.configureAuth(); err != nil {
		return err
	}

	repository, err := git.PlainClone(projectPath, &git.CloneOptions{
		URL:           cfg.Repo,
		ReferenceName: plumbing.NewBranchReferenceName(cfg.Branch),
		SingleBranch:  true,
		Depth:         1,
		NoCheckout:    true,
		ClientOptions: []client.Option{cfg.auth},
	})
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}
	cfg.git = repository
	return nil
}

func (cfg *GitConfig) Fetch() error {
	if cfg.git == nil {
		return fmt.Errorf("repository is not initialized")
	}
	if err := cfg.git.Fetch(&git.FetchOptions{
		RemoteName:    "origin",
		ClientOptions: []client.Option{cfg.auth},
	}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("fetch repository: %w", err)
	}
	worktree, err := cfg.git.Worktree()
	if err != nil {
		return err
	}
	checkoutOptions := &git.CheckoutOptions{
		Branch:                    plumbing.NewBranchReferenceName(cfg.Branch),
		SparseCheckoutDirectories: []string{cfg.Path},
	}
	return worktree.Checkout(checkoutOptions)
}
