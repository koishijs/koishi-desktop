package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func UnzipFile(src string, dest string, clean bool) error {
	l.Debugf("Unzip: %s", src)

	srcFile, err := os.Open(src)
	if err != nil {
		l.Error("Failed to open source file.")
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	return Unzip(srcFile, dest, clean)
}

func Unzip(src io.Reader, dest string, clean bool) error {
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
		header, err := tReader.Next()
		if err == io.EOF {
			break
		}

		if header.Typeflag != tar.TypeDir {
			name := filepath.Join(dest, header.Name)
			dir := filepath.Dir(name)
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				l.Error("Failed to create:")
				l.Error(dir)
				return err
			}

			file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
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
