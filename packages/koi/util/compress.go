package util

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func UnzipFile(src string, dest string, clean bool, strip bool) error {
	l.Debugf("Unzip: %s", src)

	srcFile, err := os.Open(src)
	if err != nil {
		l.Error("Failed to open source file.")
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	return Unzip(srcFile, dest, clean, strip)
}

func Unzip(src io.Reader, dest string, clean bool, strip bool) error {
	var err error

	l.Debugf("To: %s", dest)
	if clean {
		err = os.RemoveAll(dest)
		if err != nil {
			l.Error("Failed to clean destination dir.")
			return err
		}
	}
	err = os.MkdirAll(dest, os.ModePerm)
	if err != nil {
		l.Error("Failed to create destination dir.")
		return err
	}

	gReader, err := gzip.NewReader(src)
	if err != nil {
		l.Error("Failed to parse gzip.")
		return err
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
			l.Error("Tar reading error:")
			l.Error(err)
			return err
		}
		if !validRelPath(f.Name) {
			err = errors.New("Tar contains invalid name: " + f.Name)
			l.Error(err)
			return err
		}

		if f.Typeflag != tar.TypeDir {
			rel := f.Name
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
				l.Error("Failed to create:")
				l.Error(dir)
				return err
			}

			file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, os.FileMode(f.Mode))
			if err != nil {
				l.Error("Failed to create:")
				l.Error(name)
				return err
			}
			_, err = io.Copy(file, tReader)
			if err != nil {
				l.Error("Failed to write:")
				l.Error(name)
				return err
			}
		}
	}

	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
