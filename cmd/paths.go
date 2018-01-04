package cmd

import (
	"os"
	"path/filepath"
)

func findRootFolder(cwd string) (string, error) {
	_, err := os.Stat(filepath.Join(cwd, "config.toml"))
	if err != nil {
		if os.IsNotExist(err) {
			parent := filepath.Dir(cwd)
			if parent == cwd || parent == "" {
				return "", nil
			}
			return findRootFolder(parent)
		}
		return "", err
	}
	return cwd, nil
}
