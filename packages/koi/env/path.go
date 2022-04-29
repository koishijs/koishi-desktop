package env

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var (
	// Log
	l = log.WithField("package", "env")

	DirName = dirName()
)

func dirName() string {
	path, err := os.Executable()
	if err != nil {
		l.Fatal("Cannot get executable.")
	}
	path, err = filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil {
		l.Fatal("Cannot get executable dir.")
	}
	return path
}

func Resolve(base string, path string) string {
	return filepath.Clean(filepath.Join(base, path))
}
