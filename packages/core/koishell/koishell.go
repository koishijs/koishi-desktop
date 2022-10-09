package koishell

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os/exec"
	"path/filepath"

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
		Args: []string{shell.path, argB64},
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
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("KoiShell exited with error: %w", err)
	}

	// Parse output
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
