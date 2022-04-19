package daemon

import (
	"koi/config"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func CreateNodeCmd(
	path string,
	args []string,
	dir string,
	useDataHome bool,
	useDataTemp bool,
) exec.Cmd {
	l.Debug("Getting env.")
	env := os.Environ()

	if useDataHome {
		l.Debug("Now replace HOME.")
		for {
			flag := true
			for i, e := range env {
				if strings.HasPrefix(e, "HOME=") {
					env = append(env[:i], env[i+1:]...)
					flag = false
					break
				}
			}

			if flag {
				break
			}
		}

		env = append(env, "HOME="+config.Config.InternalHomeDir)
	}

	if useDataTemp {
		l.Debug("Now replace TEMP/TMP.")
		for {
			flag := true
			for i, e := range env {
				if strings.HasPrefix(e, "TEMP=") || strings.HasPrefix(e, "TMP=") {
					env = append(env[:i], env[i+1:]...)
					flag = false
					break
				}
			}

			if flag {
				break
			}
		}

		env = append(env, "TEMP="+config.Config.InternalTempDir)
		env = append(env, "TMP="+config.Config.InternalTempDir)
	}

	l.Debug("Now replace PATH.")
	pathEnv := ""
	for _, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			pathEnv = e
			break
		}
	}
	if pathEnv == "" {
		pathEnv = "PATH="
	}
	pathEnv = pathEnv[5:]
	var pathSepr string
	if runtime.GOOS == "windows" {
		pathSepr = ";"
	} else {
		pathSepr = ":"
	}
	if pathEnv != "" {
		pathEnv += pathSepr
	}
	pathEnv += config.Config.InternalNodeDir

	return exec.Cmd{
		Path:         path,
		Args:         args,
		Env:          env,
		Dir:          dir,
		Stdin:        nil,
		Stdout:       nil,
		Stderr:       nil,
		ExtraFiles:   nil,
		SysProcAttr:  nil,
		Process:      nil,
		ProcessState: nil,
	}
}
