package config

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
)

var defaultConfigData = koiconfig.ConfigData{
	Mode:    "cli",
	Open:    true,
	Isolate: "normal",
	Start:   nil,
	Env:     nil,
}

var redirectPath = (func() *yaml.Path {
	r, _ := yaml.PathString("$.redirect")

	return r
})()

func BuildLoadConfig(path string) func(i *do.Injector) (*koiconfig.Config, error) {
	return func(i *do.Injector) (*koiconfig.Config, error) {
		exePath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("cannot get executable: %w", err)
		}
		exePath = filepath.Clean(exePath)
		dirExe := filepath.Dir(exePath)

		config := &koiconfig.Config{
			Data: defaultConfigData,
			Computed: koiconfig.ConfigComputed{
				Exe:    exePath,
				DirExe: dirExe,
			},
		}

		return config, loadConfigIntl(i, config, path, 1)
	}
}

func loadConfigIntl(i *do.Injector, c *koiconfig.Config, path string, recur uint8) error {
	var err error

	l := do.MustInvoke[*logger.Logger](i)

	if recur >= 64 {
		return fmt.Errorf("infinite redirection detected. Check your koi.config file")
	}

	l.Debugf("Loading config: %s", path)

	absPath := path
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(c.Computed.DirExe, absPath)
	}
	c.Computed.DirConfig = filepath.Dir(absPath)

	l.Debugf("Reading config file: %s", absPath)
	l.Debugf("Config dir: %s", c.Computed.DirConfig)
	file, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", absPath, err)
	}

	l.Debug("Detecting redirect field.")
	var redirect string
	err = redirectPath.Read(strings.NewReader(string(file)), &redirect)
	if err == nil {
		l.Debugf("'redirect' field detected: %s", redirect)

		return loadConfigIntl(i, c, filepath.Join(c.Computed.DirConfig, redirect), recur+1)
	}

	l.Debug("Parsing config.")
	err = yaml.Unmarshal(file, &(c.Data))
	if err != nil {
		return fmt.Errorf("failed to parse config %s: %w", absPath, err)
	}

	err = postConfig(c)
	if err != nil {
		return fmt.Errorf("failed to process postconfig: %w", err)
	}

	l.Debug("Config parsed successfully.")

	return nil
}

func postConfig(c *koiconfig.Config) error {
	var err error

	c.Computed.DirData, err = joinAndCreate(c.Computed.DirConfig, "data")
	if err != nil {
		return fmt.Errorf("failed to process dir data: %w", err)
	}
	c.Computed.DirHome, err = joinAndCreate(c.Computed.DirData, "home")
	if err != nil {
		return fmt.Errorf("failed to process dir data/home: %w", err)
	}
	c.Computed.DirNode, err = joinAndCreate(c.Computed.DirData, "node")
	if err != nil {
		return fmt.Errorf("failed to process dir data/node: %w", err)
	}
	if runtime.GOOS == "windows" {
		c.Computed.DirNodeExe = c.Computed.DirNode
	} else {
		c.Computed.DirNodeExe, err = joinAndCreate(c.Computed.DirNode, "bin")
		if err != nil {
			return fmt.Errorf("failed to process dir node/bin: %w", err)
		}
	}
	c.Computed.DirLock, err = joinAndCreate(c.Computed.DirData, "lock")
	if err != nil {
		return fmt.Errorf("failed to process dir data/lock: %w", err)
	}
	c.Computed.DirTemp, err = joinAndCreate(c.Computed.DirData, "tmp")
	if err != nil {
		return fmt.Errorf("failed to process dir data/tmp: %w", err)
	}
	c.Computed.DirInstance, err = joinAndCreate(c.Computed.DirData, "instances")
	if err != nil {
		return fmt.Errorf("failed to process dir data/instances: %w", err)
	}

	return nil
}

func joinAndCreate(base, path string) (string, error) {
	var err error

	joinedPath := filepath.Join(base, path)
	err = os.MkdirAll(joinedPath, fs.ModePerm) // -rwxrwxrwx
	if err != nil {
		return "", fmt.Errorf("failed to create data folder %s: %w", path, err)
	}
	// Set perm for directory that already exists
	err = os.Chmod(joinedPath, fs.ModePerm) // -rwxrwxrwx
	if err != nil {
		return "", fmt.Errorf("failed to chmod data folder %s: %w", path, err)
	}

	return joinedPath, nil
}
