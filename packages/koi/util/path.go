package util

import (
	"os"
	"path/filepath"
)

func Resolve(base string, path string, ensureExists bool) (string, error) {
	if base != "" {
		path = filepath.Join(base, path)
	}
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	if ensureExists {
		_, err = os.Stat(path)
		if err != nil {
			return "", err
		}
	}
	return path, nil
}
