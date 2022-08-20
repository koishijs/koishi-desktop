package compress

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ExtractZipFile(src string, dest string) error {
	var err error

	reader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", src, err)
	}
	defer func() {
		_ = reader.Close()
	}()

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}

	for _, f := range reader.File {
		err = extractZipFileIntl(dest, f)
		if err != nil {
			return fmt.Errorf("failed to write %s: %w", f.Name, err)
		}
	}

	return nil
}

func extractZipFileIntl(dest string, f *zip.File) error {
	var err error

	if !validRelPath(f.Name) {
		return fmt.Errorf("zipslip file detected: %s", f.Name)
	}

	path := filepath.Join(dest, f.Name)

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

		_, err = io.Copy(file, reader)
		if err != nil {
			return err
		}
	}

	return nil
}
