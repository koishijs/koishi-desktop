package util

import (
	"os"
	"path/filepath"
)

func Resolve(base string, path string, ensureExists bool) (string, error) {
	if base != "" {
		path = filepath.Join(base, path)
	}
	if ensureExists {
		ePath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return "", err
		}
		_, err = os.Stat(ePath)
		if err != nil {
			return "", err
		}
		path = ePath
	}
	return path, nil
}
