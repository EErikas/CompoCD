package sops

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getsops/sops/v3/decrypt"
)

func IsEncryptedComposeFile(path string) bool {
	name := filepath.Base(path)
	return strings.HasSuffix(name, ".sops.yaml") || strings.HasSuffix(name, ".sops.yml")
}

func PrepareFile(path string) (string, func() error, error) {
	if !IsEncryptedComposeFile(path) {
		return path, func() error { return nil }, nil
	}

	encrypted, err := os.ReadFile(path)
	if err != nil {
		return "", nil, fmt.Errorf("read encrypted Compose file %q: %w", path, err)
	}
	cleartext, err := decrypt.Data(encrypted, "yaml")
	if err != nil {
		return "", nil, fmt.Errorf("decrypt Compose file %q with SOPS: %w", path, err)
	}
	defer clear(cleartext)

	temp, err := os.CreateTemp("", "compocd-compose-*.yaml")
	if err != nil {
		return "", nil, fmt.Errorf("create temporary Compose file: %w", err)
	}
	tempPath := temp.Name()
	cleanup := func() error {
		err := os.Remove(tempPath)
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if _, err := temp.Write(cleartext); err != nil {
		closeErr := temp.Close()
		return "", nil, errors.Join(
			fmt.Errorf("write temporary Compose file: %w", err),
			closeErr,
			cleanup(),
		)
	}
	if err := temp.Close(); err != nil {
		return "", nil, errors.Join(
			fmt.Errorf("close temporary Compose file: %w", err),
			cleanup(),
		)
	}

	return tempPath, cleanup, nil
}
