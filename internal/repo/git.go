package repo

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/go-git/go-git/v6/plumbing"
	"github.com/go-git/go-git/v6/plumbing/client"
	httptransport "github.com/go-git/go-git/v6/plumbing/transport/http"
	sshtransport "github.com/go-git/go-git/v6/plumbing/transport/ssh"
)

func (cfg *GitConfig) configureAuth() error {
	cfg.auth = nil
	switch cfg.AuthType {
	case authNone:
		return nil
	case authHTTPS:
		if cfg.AuthSecret == "" {
			return fmt.Errorf("HTTPS authentication requires an access token")
		}
		cfg.auth = client.WithHTTPAuth(&httptransport.BasicAuth{
			Username: "oauth2",
			Password: cfg.AuthSecret,
		})
		return nil
	case authSSH:
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

func (cfg *GitConfig) clientOptions() []client.Option {
	if cfg.auth == nil {
		return nil
	}
	return []client.Option{cfg.auth}
}

func (cfg *GitConfig) Init() error {
	if cfg.Repo == "" {
		return fmt.Errorf("git repo is required")
	}
	if cfg.Branch == "" {
		return fmt.Errorf("git branch is required")
	}
	if err := cfg.configureAuth(); err != nil {
		return fmt.Errorf("auth error: %w", err)
	}

	repoExists, err := cfg.checkRepo()
	if err != nil {
		return err
	}
	if repoExists {
		log.Print("Repo already exists, skipping cloning... ")
		return nil
	}

	repository, err := git.PlainClone(cfg.DiskPath, &git.CloneOptions{
		URL:           cfg.Repo,
		ReferenceName: plumbing.NewBranchReferenceName(cfg.Branch),
		SingleBranch:  true,
		Depth:         1,
		NoCheckout:    true,
		ClientOptions: cfg.clientOptions(),
	})
	if err != nil {
		return fmt.Errorf("clone repository: %w", err)
	}
	cfg.git = repository
	return nil
}

func (cfg *GitConfig) checkRepo() (bool, error) {
	repository, err := git.PlainOpen(cfg.DiskPath)
	if err == nil {
		cfg.git = repository
		return true, nil
	}
	if errors.Is(err, git.ErrRepositoryNotExists) {
		return false, nil
	}
	return false, fmt.Errorf("open repository: %w", err)
}

func (cfg *GitConfig) Fetch() error {
	if cfg.git == nil {
		return fmt.Errorf("repository is not initialized")
	}
	branch := plumbing.NewBranchReferenceName(cfg.Branch)
	remoteBranch := plumbing.NewRemoteReferenceName("origin", cfg.Branch)
	refSpec := config.RefSpec("+" + branch.String() + ":" + remoteBranch.String())
	if err := cfg.git.Fetch(&git.FetchOptions{
		RemoteName:    "origin",
		RefSpecs:      []config.RefSpec{refSpec},
		Force:         true,
		ClientOptions: cfg.clientOptions(),
	}); err != nil && !errors.Is(err, git.NoErrAlreadyUpToDate) {
		return fmt.Errorf("fetch repository: %w", err)
	}
	remoteRef, err := cfg.git.Reference(remoteBranch, false)
	if err != nil {
		return fmt.Errorf("find fetched branch %q: %w", cfg.Branch, err)
	}

	worktree, err := cfg.git.Worktree()
	if err != nil {
		return fmt.Errorf("open worktree: %w", err)
	}
	sparseDirs := sparseCheckoutDirectories(cfg.RepoPath)
	if err := worktree.Checkout(&git.CheckoutOptions{
		Hash:                      remoteRef.Hash(),
		Force:                     true,
		SparseCheckoutDirectories: sparseDirs,
	}); err != nil {
		return fmt.Errorf("reset worktree to fetched branch %q: %w", cfg.Branch, err)
	}
	return nil
}

func sparseCheckoutDirectories(repoPath string) []string {
	repoPath = path.Clean(repoPath)
	if repoPath == "." {
		return nil
	}
	return []string{repoPath}
}
