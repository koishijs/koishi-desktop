package config

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"gopkg.ilharper.com/koi/core/logger"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Data     ConfigData
	Computed ConfigComputed
}

//goland:noinspection GoNameStartsWithPackageName
type ConfigData struct {
	Mode    string `yaml:"mode"`
	Open    bool   `yaml:"open"`
	Isolate string `yaml:"isolate"`
	Start   []string
	Env     []string `yaml:"env"`
}

var defaultConfigData = ConfigData{
	Mode:    "cli",
	Open:    true,
	Isolate: "normal",
	Start:   nil,
	Env:     nil,
}

//goland:noinspection GoNameStartsWithPackageName
type ConfigComputed struct {
	DirExe      string
	DirConfig   string
	DirData     string
	DirHome     string
	DirNode     string
	DirNodeExe  string
	DirLock     string
	DirTemp     string
	DirInstance string
}

var (
	redirectPath = (func() *yaml.Path {
		r, _ := yaml.PathString("$.redirect")
		return r
	})()
)

func LoadConfig(l *logger.Logger, path string) (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("cannot get executable: %w", err)
	}
	dirExe := filepath.Dir(exePath)

	config := &Config{
		Data: defaultConfigData,
		Computed: ConfigComputed{
			DirExe: dirExe,
		},
	}
	return config, loadConfigIntl(config, l, path, 1)
}

func loadConfigIntl(c *Config, l *logger.Logger, path string, recur uint8) (err error) {
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
		return loadConfigIntl(c, l, filepath.Join(c.Computed.DirConfig, redirect), recur+1)
	}

	l.Debug("Parsing config.")
	err = yaml.Unmarshal(file, &(c.Data))
	if err != nil {
		return fmt.Errorf("failed to parse config %s: %w", absPath, err)
	}

	l.Debug("Config parsed successfully.")
	l.Debug("Now processing postConfig.")
	err = postConfig(c)
	if err != nil {
		return fmt.Errorf("failed to process postconfig: %w", err)
	}

	return
}

func postConfig(c *Config) (err error) {
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

	return
}

func joinAndCreate(base, path string) (joinedPath string, err error) {
	joinedPath = filepath.Join(base, path)
	err = os.MkdirAll(joinedPath, fs.ModePerm) // -rwxrwxrwx
	if err != nil {
		return "", err
	}
	// Set perm for directory that already exists
	err = os.Chmod(joinedPath, fs.ModePerm) // -rwxrwxrwx
	if err != nil {
		return "", err
	}
	return joinedPath, nil
}
