package util

import (
	"koi/util/logger"
	"os"
	"path/filepath"
)

var (
	DirName = getDirName()
)

func getDirName() string {
	path, err := os.Executable()
	if err != nil {
		logger.Fatal("Cannot get executable.")
	}
	path, err = Resolve("", filepath.Dir(path))
	if err != nil {
		logger.Fatal("Cannot get executable dir.")
	}
	return path
}

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
