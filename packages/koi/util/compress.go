package util

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func Unzip(src string, dest string, clean bool) error {
	rSrc, err := Resolve("", src, true)
	if err != nil {
		l.Errorf("Failed to resolve %s", src)
		return err
	}
	l.Debugf("Unzip: %s", rSrc)

	rDest, err := Resolve("", dest, false)
	if err != nil {
		l.Errorf("Failed to resolve %s", dest)
		return err
	}
	l.Debugf("To: %s", rDest)
	if clean {
		err = os.RemoveAll(rDest)
		if err != nil {
			l.Error("Failed to clean destination dir.")
			return err
		}
	}
	err = os.MkdirAll(rDest, os.ModePerm)
	if err != nil {
		l.Error("Failed to create destination dir.")
		return err
	}

	srcFile, err := os.Open(rSrc)
	if err != nil {
		l.Error("Failed to open source file.")
		return err
	}
	defer func() {
		_ = srcFile.Close()
	}()

	gReader, err := gzip.NewReader(srcFile)
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
			name := filepath.Join(rDest, header.Name)
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
