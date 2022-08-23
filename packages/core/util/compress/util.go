package compress

import (
	"fmt"
	"path/filepath"
	"strings"
)

// https://snyk.io/research/zip-slip-vulnerability
// https://github.com/securego/gosec/issues/324#issuecomment-935927967
func sanitizeArchivePath(dest string, filename string) (string, error) {
	v := filepath.Join(dest, filename)
	if strings.HasPrefix(v, filepath.Clean(dest)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", filename)
}
