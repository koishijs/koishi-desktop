package env

import (
	log "github.com/sirupsen/logrus"
	"koi/util"
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
	path, err = util.Resolve("", filepath.Dir(path), true)
	if err != nil {
		l.Fatal("Cannot get executable dir.")
	}
	return path
}
