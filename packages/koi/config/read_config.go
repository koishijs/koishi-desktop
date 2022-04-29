package config

import (
	"errors"
	"github.com/goccy/go-yaml"
	"koi/env"
	"koi/util"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	redirectPath = createRedirectPath()
)

func ReadConfig(path string) (*KoiConfig, error) {
	config, err := readConfigIntl(path, 1)
	if err != nil {
		l.Error(err.Error())
		return nil, err
	}
	return config, nil
}

func createRedirectPath() *yaml.Path {
	r, err := yaml.PathString("$.redirect")
	if err != nil {
		l.Fatal("Err create redirect yaml path.")
	}
	return r
}

func readConfigIntl(path string, recur int) (*KoiConfig, error) {
	if recur >= 64 {
		return nil, errors.New("infinite redirection detected. Check your koi.config file")
	}

	l.Debugf("Loading config: %s", path)

	absPath := path
	if !filepath.IsAbs(absPath) {
		newAbsPath, err := util.Resolve(env.DirName, absPath, true)
		if err != nil {
			l.Error("Failed to resolve config path:")
			l.Fatal(absPath)
		}
		absPath = newAbsPath
	}
	configDir := filepath.Dir(absPath)

	l.Debugf("Reading config file: %s", absPath)
	l.Debugf("Config dir: %s", configDir)
	file, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	l.Debug("Detecting redirect field.")
	var redirect string
	err = redirectPath.Read(strings.NewReader(string(file)), &redirect)
	if err == nil {
		l.Debugf("'redirect' field detected: %s", redirect)
		return readConfigIntl(filepath.Join(configDir, redirect), recur+1)
	}

	l.Debug("Parsing config.")
	config := new(KoiConfig)
	*config = defaultConfig
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	l.Debug("Config parsed successfully.")

	l.Debug("Now processing postConfig.")
	config.InternalConfigDir = configDir
	postConfig(config)

	return config, nil
}

func postConfig(c *KoiConfig) {
	dir, err := util.Resolve(c.InternalConfigDir, "data", true)
	if err != nil {
		l.Error("Failed to resolve data dir. Config dir:")
		l.Fatal(c.InternalConfigDir)
	}
	c.InternalDataDir = dir

	if c.UseDataHome {
		dir, err = util.Resolve(c.InternalDataDir, "home", true)
		if err != nil {
			l.Error("Failed to resolve home dir. Data dir:")
			l.Fatal(c.InternalDataDir)
		}
		c.InternalHomeDir = dir
	}

	dir, err = util.Resolve(c.InternalDataDir, "node", true)
	if err != nil {
		l.Error("Failed to resolve node dir. Data dir:")
		l.Fatal(c.InternalDataDir)
	}
	c.InternalNodeDir = dir
	if runtime.GOOS == "windows" {
		c.InternalNodeExeDir = c.InternalNodeDir
	} else {
		dir, err = util.Resolve(c.InternalNodeDir, "bin", true)
		if err != nil {
			l.Error("Failed to resolve node binary dir. Node dir:")
			l.Fatal(c.InternalNodeDir)
		}
		c.InternalNodeExeDir = dir
	}

	if c.UseDataTemp {
		dir, err = util.Resolve(c.InternalDataDir, "tmp", true)
		if err != nil {
			l.Error("Failed to resolve temp dir. Data dir:")
			l.Fatal(c.InternalDataDir)
		}
		c.InternalTempDir = dir
	}

	dir, err = util.Resolve(c.InternalDataDir, "instances", true)
	if err != nil {
		l.Error("Failed to resolve instance dir. Data dir:")
		l.Fatal(c.InternalDataDir)
	}
	c.InternalInstanceDir = dir
}
