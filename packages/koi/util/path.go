package util

import (
	"os"
	"path/filepath"
)

// Resolve the existing file or dir.
// This will filepath.Join base and path (if has base),
// filepath.EvalSymlinks and finally
// os.Stat to ensure it exists.
func Resolve(base string, path string) (string, error) {
	if base != "" {
		path = filepath.Join(base, path)
	}
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}
	_, err = os.Stat(path)
	if err != nil {
		return "", err
	}
	return path, nil
}
