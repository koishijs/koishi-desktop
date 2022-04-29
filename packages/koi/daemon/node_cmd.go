package daemon

import (
	log "github.com/sirupsen/logrus"
	"io"
	"koi/config"
	"koi/util"
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
	Cmd       exec.Cmd
	errReader *io.PipeReader
	errWriter *io.PipeWriter
}

func RunNodeCmd(
	nodeExe string,
	args []string,
	dir string,
) error {
	cmd := CreateNodeCmd(nodeExe, args, dir)
	return cmd.Run()
}

func CreateNodeCmd(
	nodeExe string,
	args []string,
	dir string,
) NodeCmd {
	l.Debug("Getting env.")
	env := os.Environ()

	if config.Config.UseDataHome {
		l.Debug("Now replace HOME/USERPROFILE.")
		for {
			notFound := true
			for i, e := range env {
				if strings.HasPrefix(e, "HOME=") || strings.HasPrefix(e, "USERPROFILE=") {
					env = append(env[:i], env[i+1:]...)
					notFound = false
					break
				}
			}

			if notFound {
				break
			}
		}

		env = append(env, "HOME="+config.Config.InternalHomeDir)
		env = append(env, "USERPROFILE="+config.Config.InternalHomeDir)
		l.Debugf("HOME=%s", config.Config.InternalHomeDir)
	}

	if config.Config.UseDataTemp {
		l.Debug("Now replace TMPDIR/TEMP/TMP.")
		for {
			notFound := true
			for i, e := range env {
				if strings.HasPrefix(e, "TMPDIR=") || strings.HasPrefix(e, "TEMP=") || strings.HasPrefix(e, "TMP=") {
					env = append(env[:i], env[i+1:]...)
					notFound = false
					break
				}
			}

			if notFound {
				break
			}
		}

		env = append(env, "TMPDIR="+config.Config.InternalTempDir)
		env = append(env, "TEMP="+config.Config.InternalTempDir)
		env = append(env, "TMP="+config.Config.InternalTempDir)
		l.Debugf("TEMP=%s", config.Config.InternalTempDir)
	}

	l.Debug("Now replace PATH.")
	pathEnv := ""
	for {
		notFound := true
		for i, e := range env {
			if strings.HasPrefix(e, "PATH=") {
				pathEnv = e[5:]
				env = append(env[:i], env[i+1:]...)
				notFound = false
				break
			}
		}

		if notFound {
			break
		}
	}
	var pathSepr string
	if runtime.GOOS == "windows" {
		pathSepr = ";"
	} else {
		pathSepr = ":"
	}
	if pathEnv != "" && !config.Config.Strict {
		pathEnv = config.Config.InternalNodeExeDir + pathSepr + pathEnv
	} else {
		pathEnv = config.Config.InternalNodeExeDir
	}
	env = append(env, "PATH="+pathEnv)
	l.Debugf("PATH=%s", pathEnv)

	l.Debugf("PWD=%s", dir)

	l.Debug("Now constructing NodeCmd.")
	errReader, errWriter := io.Pipe()
	cmdPath, err := util.Resolve(config.Config.InternalNodeExeDir, nodeExe, true)
	if err != nil {
		l.Fatalf("Failed to resolve node executable: %s", nodeExe)
	}
	cmdArgs := []string{cmdPath}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Cmd{
		Path:         cmdPath,
		Args:         cmdArgs,
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
		Cmd:       cmd,
		errReader: errReader,
		errWriter: errWriter,
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
			n, err := c.errReader.Read(p)
			if err == io.EOF {
				break
			}
			if err != nil {
				l.Debugf("stderr.Read() err: %s", err)
			}
			lKoishi.Info(util.Trim(string(p[:n])))
		}
	}()

	l.Debug("Now start NodeCmd.")
	return c.Cmd.Start()
}

func (c *NodeCmd) Wait() error {
	l.Debug("Now wait NodeCmd.")

	defer func() {
		err := c.errWriter.Close()
		if err != nil {
			l.Debug("Stderr closed with err.")
			l.Debug(err)
		}
	}()

	return c.Cmd.Wait()
}
