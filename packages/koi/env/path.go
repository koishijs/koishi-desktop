package env

import (
	"koi/util"
	l "koi/util/logger"
	"os"
	"path/filepath"
)

var (
	DirName = dirName()
)

func dirName() string {
	path, err := os.Executable()
	if err != nil {
		l.Fatal("Cannot get executable.")
	}
	path, err = util.Resolve("", filepath.Dir(path))
	if err != nil {
		l.Fatal("Cannot get executable dir.")
	}
	return path
}
