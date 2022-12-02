package main

import (
	"os"
	"path/filepath"
)

func legacyUserDataDir() (string, error) {
	return filepath.Join(os.Getenv("APPDATA"), "Il Harper/Koishi"), nil
}
