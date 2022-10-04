//nolint:wrapcheck
package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func extract(dest string) error {
	var err error

	reader, err := zip.NewReader(bytes.NewReader(portableData), int64(len(portableData)))
	if err != nil {
		return fmt.Errorf("failed to open portable data: %w", err)
	}

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}

	for _, f := range reader.File {
		err = extractIntl(dest, f)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", f.Name, err)
		}
	}

	return nil
}

func extractIntl(dest string, f *zip.File) error {
	var err error

	path, err := sanitizeArchivePath(dest, f.Name)
	if err != nil {
		return err
	}

	if f.FileInfo().IsDir() {
		err = os.MkdirAll(path, f.Mode())
		if err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(path), f.Mode())
		if err != nil {
			return err
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer func() {
			_ = file.Close()
		}()

		reader, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			_ = reader.Close()
		}()

		_, err = io.Copy(file, reader) //nolint:gosec
		if err != nil {
			return err
		}
	}

	return nil
}

// https://snyk.io/research/zip-slip-vulnerability
// https://github.com/securego/gosec/issues/324#issuecomment-935927967
func sanitizeArchivePath(dest string, filename string) (string, error) {
	v := filepath.Join(dest, filename)
	if strings.HasPrefix(v, filepath.Clean(dest)) {
		return v, nil
	}

	return "", fmt.Errorf("%s: %s", "content filepath is tainted", filename)
}
