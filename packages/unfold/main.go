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

	switch mode {
	case "ensure":
		_, err := os.Stat(pathConfig)

		switch {
		case errors.Is(err, fs.ErrNotExist):
			fmt.Println("User data does not exist. Extracting all files.")
			extractMode = ExtractData | ExtractConfig
		case err == nil:
			fmt.Println("User data exists. Setting up redirect.")
			extractMode = 0
		default:
			fmt.Printf("Failed to stat config %s: %v\n", pathConfig, err)
			os.Exit(1)
		}
	case "reset-data":
		extractMode = ExtractData
	default:
		extractMode = ExtractData | ExtractConfig
	}

	fmt.Printf("Extract mode: %v\n", extractMode)

	err = extract(folderData, extractMode)
	if err != nil {
		fmt.Printf("Failed to unfold: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(pathConfigRedirect, []byte("redirect: USERDATA"), 0o644) //nolint:gosec
	if err != nil {
		fmt.Printf("Failed to setup redirect: %v\n", err)
		fmt.Println("This is not a bug and unfold will continue.")
	}

	fmt.Println("Unfold complete.")
}

func help() {
	fmt.Print(`
"unfold" for Cordis Desktop

Available modes:
	ensure      - Ensure user data folder exists. Used for installation.
	reset-data  - Reset user data folder.
	reset-all   - Reset all data and config.
`)
}
