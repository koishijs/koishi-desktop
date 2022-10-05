//go:build windows

package pathutil

import (
	"os"
	"path/filepath"
)

func UserDataDir() (string, error) {
	return filepath.Join(os.Getenv("APPDATA"), "koishijs/koishi-desktop"), nil
}
