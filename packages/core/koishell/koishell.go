package koishell

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/goccy/go-json"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/koi/core/util/killdren"
)

const deltaCh uint16 = 3000

type KoiShell struct {
	i *do.Injector

	path string
	cwd  string

	// The mutex lock.
	//
	// There's no need to use [sync.RWMutex]
	// because almost all ops are write.
	mutex sync.Mutex
	wg    sync.WaitGroup
	reg   [256]*exec.Cmd
}

func BuildKoiShell(path string) func(i *do.Injector) (*KoiShell, error) {
	return func(i *do.Injector) (*KoiShell, error) {
		return &KoiShell{
			i:    i,
			path: path,
			cwd:  filepath.Dir(path),
		}, nil
	}
}

func (shell *KoiShell) getIndex(cmd *exec.Cmd) uint8 {
	shell.mutex.Lock()
	defer shell.mutex.Unlock()

	var index uint8 = 0
	for {
		if shell.reg[index] == nil {
			shell.reg[index] = cmd

			return index
		}
		index++
	}
}

func (shell *KoiShell) exec(arg any) (map[string]any, error) {
	var err error

	l := do.MustInvoke[*logger.Logger](shell.i)

	argJson, err := json.Marshal(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal arg %#+v: %w", arg, err)
	}
	argB64 := base64.StdEncoding.EncodeToString(argJson)

	cmd := &exec.Cmd{
		Path: shell.path,
		Args: []string{argB64},
		Dir:  shell.cwd,
	}
	killdren.Set(cmd)

	var outB64 string

	// Setup IO pipes
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	go func(scn *bufio.Scanner) {
		for {
			if !scn.Scan() {
				break
			}
			scnErr := scn.Err()
			if scnErr != nil {
				l.Warn(fmt.Errorf("KoiShell scanner error: %w", scnErr))
			} else {
				outB64 = scn.Text()
			}
		}
	}(bufio.NewScanner(stdoutPipe))

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	go func(scn *bufio.Scanner) {
		for {
			if !scn.Scan() {
				break
			}
			scnErr := scn.Err()
			if scnErr != nil {
				l.Warn(fmt.Errorf("KoiShell scanner error: %w", scnErr))
			} else {
				l.Error(scn.Text())
			}
		}
	}(bufio.NewScanner(stderrPipe))

	// Wait process to stop
	index := shell.getIndex(cmd)
	shell.wg.Add(1)

	l.Debugf("Starting KoiShell process.\narg: %#+v\nargJson: %s\nargB64: %s", arg, argJson, argB64)
	err = cmd.Run()

	shell.wg.Done()

	shell.mutex.Lock()
	shell.reg[index] = nil
	shell.mutex.Unlock()

	if err != nil {
		return nil, fmt.Errorf("KoiShell exited with error: %w", err)
	} else {
		l.Debugf("KoiShell successfully exited.")
	}

	// Parse output
	if outB64 == "" {
		return nil, nil
	}

	outJson, err := base64.StdEncoding.DecodeString(outB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode KoiShell output %s: %w", outB64, err)
	}

	var out map[string]any
	err = json.Unmarshal(outJson, &out)
	if err != nil {
		return nil, fmt.Errorf("failed to parse KoiShell output %s: %w", outJson, err)
	}

	return out, nil
}

func (shell *KoiShell) Shutdown() error {
	l := do.MustInvoke[*logger.Logger](shell.i)

	l.Debug("Shutting down DaemonProcess.")

	shell.mutex.Lock()

	for _, cmd := range shell.reg {
		if cmd != nil {
			err := killdren.Stop(cmd)
			if err != nil {
				l.Debugf("failed to gracefully stop KoiShell %d: %v. Trying kill", cmd.Process.Pid, err)
				_ = killdren.Kill(cmd)
			}
		}
	}

	shell.mutex.Unlock()
	shell.wg.Wait()

	// Do not short other do.Shutdownable
	return nil
}

func (shell *KoiShell) WebView(name, url string) error {
	_, err := shell.exec(map[string]string{
		"mode": "webview",
		"name": name,
		"url":  url,
	})

	return err
}
