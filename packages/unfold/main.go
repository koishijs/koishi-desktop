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

	if mode == "ensure" {
		_, err := os.Stat(pathConfig)
		if err == nil {
			fmt.Println("User data exists. Skip unfolding.")

			err = os.WriteFile(pathConfigRedirect, []byte("redirect: USERDATA"), 0o644)
			if err != nil {
				fmt.Printf("Failed to setup redirect: %v\n", err)
				os.Exit(1)
			}

			return
		}
		if !errors.Is(err, fs.ErrNotExist) {
			fmt.Printf("Failed to stat config %s: %v\n", pathConfig, err)
			os.Exit(1)
		}
	}

	if mode == "reset-data" {
		err = extract(folderData, false)
	} else {
		err = extract(folderData, true)
	}
	if err != nil {
		fmt.Printf("Failed to unfold: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(pathConfigRedirect, []byte("redirect: USERDATA"), 0o644)
	if err != nil {
		fmt.Printf("Failed to setup redirect: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Unfold complete.")
}

func setupRedirect(pathConfigRedirect string) {
	err := os.WriteFile(pathConfigRedirect, []byte(""), 0o644)
	if err != nil {
		fmt.Printf("Failed to setup redirect: %v", err)
		os.Exit(1)
	}
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
