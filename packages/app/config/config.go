package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/samber/do"
	"golang.org/x/text/message"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/pathutil"
	"gopkg.ilharper.com/x/killdren"
)

var defaultConfigData = koiconfig.ConfigData{
	Mode:    "ui",
	Open:    "auto",
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
		p := do.MustInvoke[*message.Printer](i)

		exePath, err := os.Executable()
		if err != nil {
			return nil, errors.New(p.Sprintf("cannot get executable: %v", err))
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
	p := do.MustInvoke[*message.Printer](i)

	if recur >= 64 {
		return errors.New(p.Sprintf("infinite redirection detected. Check your koi.config file"))
	}

	l.Debugf("Loading config: %s", path)

	absPath := path
	if !filepath.IsAbs(absPath) {
		absPath = filepath.Join(c.Computed.DirExe, absPath)
	}
	c.Computed.DirConfig = filepath.Dir(absPath)

	l.Debug(p.Sprintf("Reading config file: %s", absPath))
	l.Debug(p.Sprintf("Config dir: %s", c.Computed.DirConfig))
	file, err := os.ReadFile(absPath)
	if err != nil {
		return errors.New(p.Sprintf("failed to read %s: %w", absPath, err))
	}

	l.Debug(p.Sprintf("Detecting redirect field."))
	var redirect string
	err = redirectPath.Read(strings.NewReader(string(file)), &redirect)
	if err == nil {
		l.Debug(p.Sprintf("'redirect' field detected: %s", redirect))
		if redirect == "USERDATA" {
			r, uddErr := pathutil.UserDataDir()
			if uddErr != nil {
				return errors.New(p.Sprintf("failed to resolve user data: %w", uddErr))
			}
			redirect = filepath.Join(r, "koi.yml")

			_, rfErr := os.ReadFile(redirect)
			if rfErr != nil {
				l.Debug(p.Sprintf("Failed to read %s: %v", redirect, rfErr))
				l.Debug(p.Sprintf("Trying to unfold."))

				var command string
				if runtime.GOOS == "windows" {
					command = "unfold.exe"
				} else {
					command = "unfold"
				}
				cmdPath := filepath.Join(c.Computed.DirExe, command)

				cmd := &exec.Cmd{
					Path: cmdPath,
					Args: []string{cmdPath, "ensure"},
					Dir:  c.Computed.DirExe,
				}
				killdren.Set(cmd)

				// Setup IO pipes
				stdoutPipe, err := cmd.StdoutPipe()
				if err != nil {
					return errors.New(p.Sprintf("failed to create stdout pipe: %v", err))
				}

				go func(scn *bufio.Scanner) {
					for {
						if !scn.Scan() {
							break
						}
						scnErr := scn.Err()
						if scnErr != nil {
							l.Warn(p.Sprintf("Scanner error: %s", scnErr))
						} else {
							l.Info(scn.Text())
						}
					}
				}(bufio.NewScanner(stdoutPipe))

				runErr := cmd.Run()
				if runErr != nil {
					l.Debug(p.Sprintf("Failed to unfold: %v", runErr))
					l.Debug(p.Sprintf("Will ignore this error."))
				}
			}

			l.Debug(p.Sprintf("Redirecting to user data: %s", redirect))
		}

		return loadConfigIntl(i, c, redirect, recur+1)
	}

	l.Debug(p.Sprintf("Parsing config."))
	err = yaml.Unmarshal(file, &(c.Data))
	if err != nil {
		return errors.New(p.Sprintf("failed to parse config %s: %v", absPath, err))
	}

	err = postConfig(i, c)
	if err != nil {
		return errors.New(p.Sprintf("failed to process postconfig: %v", err))
	}

	l.Debug(p.Sprintf("Config parsed successfully."))

	return nil
}

func postConfig(i *do.Injector, c *koiconfig.Config) error {
	var err error

	p := do.MustInvoke[*message.Printer](i)

	c.Computed.DirData, err = joinAndCreate(i, c.Computed.DirConfig, "data")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data: %v", err))
	}
	c.Computed.DirHome, err = joinAndCreate(i, c.Computed.DirData, "home")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data/home: %v", err))
	}
	c.Computed.DirNode, err = joinAndCreate(i, c.Computed.DirData, "node")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data/node: %v", err))
	}
	if runtime.GOOS == "windows" {
		c.Computed.DirNodeExe = c.Computed.DirNode
	} else {
		c.Computed.DirNodeExe, err = joinAndCreate(i, c.Computed.DirNode, "bin")
		if err != nil {
			return errors.New(p.Sprintf("failed to process dir node/bin: %v", err))
		}
	}
	c.Computed.DirLock, err = joinAndCreate(i, c.Computed.DirData, "lock")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data/lock: %v", err))
	}
	c.Computed.DirTemp, err = joinAndCreate(i, c.Computed.DirData, "tmp")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data/tmp: %v", err))
	}
	c.Computed.DirInstance, err = joinAndCreate(i, c.Computed.DirData, "instances")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data/instances: %v", err))
	}
	c.Computed.DirLog, err = joinAndCreate(i, c.Computed.DirData, "logs")
	if err != nil {
		return errors.New(p.Sprintf("failed to process dir data/log: %v", err))
	}

	return nil
}

func joinAndCreate(i *do.Injector, base, path string) (string, error) {
	var err error

	l := do.MustInvoke[*logger.Logger](i)
	p := do.MustInvoke[*message.Printer](i)

	joinedPath := filepath.Join(base, path)
	err = os.MkdirAll(joinedPath, fs.ModePerm) // -rwxrwxrwx
	if err != nil {
		return "", fmt.Errorf("failed to create data folder %s: %w", path, err)
	}
	// Set perm for directory that already exists
	err = os.Chmod(joinedPath, fs.ModePerm) // -rwxrwxrwx
	if err != nil {
		l.Warn(p.Sprintf("failed to chmod data folder %s: %v", path, err))
	}

	return joinedPath, nil
}
