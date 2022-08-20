package compress

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ExtractTgzFile(src string, dest string, strip bool) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", src, err)
	}
	defer func() {
		_ = srcFile.Close()
	}()

	return ExtractTgz(srcFile, dest, strip)
}

func ExtractTgz(src io.Reader, dest string, strip bool) error {
	var err error

	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}

	gReader, err := gzip.NewReader(src)
	if err != nil {
		return fmt.Errorf("failed to parse gzip: %w", err)
	}
	defer func() {
		_ = gReader.Close()
	}()

	tReader := tar.NewReader(gReader)

	for {
		f, err := tReader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("tar reading error: %w", err)
		}

		if !validRelPath(f.Name) {
			return fmt.Errorf("zipslip file detected: %s", f.Name)
		}

		if f.Typeflag != tar.TypeDir {
			rel := f.Name
			if rel == "pax_global_header" {
				continue
			}
			if strip {
				var i int
				le := len(f.Name)
				for i = 0; i < le; i++ {
					if f.Name[i] == '/' {
						break
					}
				}
				if i < le-1 {
					rel = f.Name[i+1:]
				}
			}

			name := filepath.Join(dest, rel)
			dir := filepath.Dir(name)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", dir, err)
			}

			file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(f.Mode))
			if err != nil {
				return fmt.Errorf("failed to create %s: %w", name, err)
			}
			_, err = io.Copy(file, tReader)
			if err != nil {
				return fmt.Errorf("failed to write %s: %w", name, err)
			}
			_ = file.Close()
		}
	}

	return nil
}
