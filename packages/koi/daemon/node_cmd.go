package daemon

import (
	log "github.com/sirupsen/logrus"
	"io"
	"koi/config"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	// Log
	lKoishi = log.WithField("package", "koishi")
)

type NodeCmd struct {
	Cmd    exec.Cmd
	stderr io.Reader
}

func CreateNodeCmd(
	path string,
	args []string,
	dir string,
	useDataHome bool,
	useDataTemp bool,
) NodeCmd {
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

	l.Debug("Now constructing NodeCmd.")
	errReader, errWriter := io.Pipe()
	cmd := exec.Cmd{
		Path:         path,
		Args:         args,
		Env:          env,
		Dir:          dir,
		Stdin:        nil,
		Stdout:       os.Stdout,
		Stderr:       errWriter,
		ExtraFiles:   nil,
		SysProcAttr:  nil,
		Process:      nil,
		ProcessState: nil,
	}

	return NodeCmd{
		Cmd:    cmd,
		stderr: errReader,
	}
}

func (c *NodeCmd) Run() error {
	l.Debug("Now run NodeCmd.")
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *NodeCmd) Start() error {
	l.Debug("Now start stderr reader.")
	go func() {
		p := make([]byte, 1024)
		for {
			n, err := c.stderr.Read(p)
			if err == io.EOF {
				break
			}
			if err != nil {
				l.Debugf("stderr.Read() err: %s", err)
			}
			lKoishi.Info(string(p[:n]))
		}
	}()

	l.Debug("Now start NodeCmd.")
	return c.Cmd.Start()
}

func (c *NodeCmd) Wait() error {
	l.Debug("Now wait NodeCmd.")
	return c.Cmd.Wait()
}
