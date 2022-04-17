package env

import (
	"os"
	"path/filepath"
)

func DirName() (string, error) {
	path, err := os.Executable()
	if err != nil {
		return "", err
	}
	path, err = filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	return path, nil
}
