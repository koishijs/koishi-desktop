//nolint:wrapcheck
package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func migrate(folderData string) bool {
	var err error

	ldd, err := legacyUserDataDir()
	if err != nil {
		fmt.Printf("Failed to resolve legacy user data directory: %v\n", err)

		return false
	}

	if ldd == "" {
		fmt.Println("Legacy user data does not exist. Nothing to migrate.")

		return false
	}

	_, err = os.Stat(ldd)
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Println("Legacy user data does not exist. Nothing to migrate.")

		return false
	} else if err != nil {
		fmt.Printf("Failed to check legacy user data directory: %v\n", err)

		return false
	}

	err = copyDir(ldd, folderData)
	if err != nil {
		fmt.Printf("Failed to migrate legacy user data: %v\n", err)

		return false
	}

	return true
}

func copyDir(src, dst string) error {
	var err error

	err = os.MkdirAll(dst, 0o755)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		sourcePath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dst, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		switch fileInfo.Mode() & os.ModeType { //nolint:exhaustive
		case os.ModeDir:
			err = copyDir(sourcePath, destPath)
			if err != nil {
				return err
			}

			_ = os.Chmod(destPath, fileInfo.Mode())
		case os.ModeSymlink:
			err = copySymLink(sourcePath, destPath)
			if err != nil {
				return err
			}
		default:
			err = copyFile(sourcePath, destPath, fileInfo.Mode())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string, perm fs.FileMode) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	return err
}

func copySymLink(src, dst string) error {
	link, err := os.Readlink(src)
	if err != nil {
		return err
	}

	return os.Symlink(link, dst)
}
