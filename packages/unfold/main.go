//nolint:typecheck
package main

import (
	_ "embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.ilharper.com/koi/core/util/pathutil"
)

//go:embed portabledata.zip
var portableData []byte

func main() {
	var err error

	if len(os.Args) != 2 {
		fmt.Println("You need to pass exactly 1 argument for mode.")
		help()
		os.Exit(1)
	}

	mode := os.Args[1]
	if mode == "help" {
		help()
		return
	}
	if mode != "ensure" && mode != "reset-data" && mode != "reset-all" {
		fmt.Println("Unknown mode.")
		help()
		os.Exit(1)
	}

	pathExe, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get executable path: %v\n", err)
		os.Exit(1)
	}
	folderBinary := filepath.Dir(pathExe)
	folderData, err := pathutil.UserDataDir()
	if err != nil {
		fmt.Printf("Failed to resolve user data directory: %v\n", err)
		os.Exit(1)
	}

	pathConfigRedirect := filepath.Join(folderBinary, "koi.yml")
	pathConfig := filepath.Join(folderData, "koi.yml")
	var extractMode uint8

	if mode == "ensure" {
		_, err := os.Stat(pathConfig)

		if errors.Is(err, fs.ErrNotExist) {
			fmt.Println("User data does not exist. Trying to migrate legacy user data.")

			if migrate(folderData) {
				fmt.Println("Migration completed. Extracting only node.")
				extractMode = EXTRACT_NODE
			} else {
				fmt.Println("Legacy user data not found or migration failed. Extracting all files.")
				extractMode = EXTRACT_DATA | EXTRACT_CONFIG | EXTRACT_NODE
			}
		} else if err == nil {
			fmt.Println("User data exists. Extracting only node.")
			extractMode = EXTRACT_NODE
		} else {
			fmt.Printf("Failed to stat config %s: %v\n", pathConfig, err)
			os.Exit(1)
		}
	} else if mode == "reset-data" {
		extractMode = EXTRACT_DATA | EXTRACT_NODE
	} else {
		extractMode = EXTRACT_DATA | EXTRACT_CONFIG | EXTRACT_NODE
	}

	fmt.Printf("Extract mode: %v\n", extractMode)

	err = extract(folderData, extractMode)
	if err != nil {
		fmt.Printf("Failed to unfold: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(pathConfigRedirect, []byte("redirect: USERDATA"), 0o644)
	if err != nil {
		fmt.Printf("Failed to setup redirect: %v\n", err)
		fmt.Println("This is not a bug and unfold will continue.")
	}

	fmt.Println("Unfold complete.")
}

func help() {
	fmt.Print(`
"unfold" for Koishi Desktop

Available modes:
	ensure      - Ensure user data folder exists. Used for installation.
	reset-data  - Reset user data folder.
	reset-all   - Reset all data and config.
`)
}
