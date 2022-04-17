package env

import (
	log "github.com/sirupsen/logrus"
	"os"
	goPath "path"
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
	return goPath.Clean(goPath.Join(base, path))
}
