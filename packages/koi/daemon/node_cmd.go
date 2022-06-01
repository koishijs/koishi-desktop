package daemon

import (
	"bufio"
	"koi/config"
	"koi/util"
	l "koi/util/logger"
	"koi/util/strutil"
	"koi/util/supcolor"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type NodeCmdOut struct {
	IsErr bool
	Text  string
}

type NodeCmd struct {
	Cmd *exec.Cmd

	// The output of NodeCmd process.
	//
	// If you set Out to nil, you need to process received NodeCmdOut
	// on a new goroutine. Otherwise, Stdout/Stderr will block.
	Out *chan NodeCmdOut
}

func RunNode(
	entry string,
	args []string,
	dir string,
) error {
	args = append([]string{entry}, args...)
	return RunNodeCmd("node", args, dir)
}

func ResolveYarn() (string, error) {
	yarnPath, err := util.Resolve(config.Config.InternalNodeExeDir, "yarn.cjs")
	if err != nil {
		l.Error("Cannot resolve yarn.")
		return "", err
	}
	return yarnPath, nil
}

func RunYarn(
	args []string,
	dir string,
) error {
	yarnPath, err := ResolveYarn()
	if err != nil {
		return err
	}
	return RunNode(yarnPath, args, dir)
}

func RunNodeCmd(
	nodeExe string,
	args []string,
	dir string,
) error {
	cmd, err := CreateNodeCmd(nodeExe, args, dir, true)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func CreateNodeCmd(
	nodeExe string,
	args []string,
	dir string,
	handleStdout bool,
) (*NodeCmd, error) {
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

		if runtime.GOOS == "windows" {
			l.Debug("Now replace APPDATA.")
			for {
				notFound := true
				for i, e := range env {
					if strings.HasPrefix(e, "APPDATA=") {
						env = append(env[:i], env[i+1:]...)
						notFound = false
						break
					}
				}

				if notFound {
					break
				}
			}

			roamingPath := filepath.Join(config.Config.InternalHomeDir, "AppData", "Roaming")
			env = append(env, "APPDATA="+roamingPath)
			l.Debugf("APPDATA=%s", roamingPath)

			l.Debug("Now replace LOCALAPPDATA.")
			for {
				notFound := true
				for i, e := range env {
					if strings.HasPrefix(e, "LOCALAPPDATA=") {
						env = append(env[:i], env[i+1:]...)
						notFound = false
						break
					}
				}

				if notFound {
					break
				}
			}

			localPath := filepath.Join(config.Config.InternalHomeDir, "AppData", "Local")
			env = append(env, "LOCALAPPDATA="+localPath)
			l.Debugf("LOCALAPPDATA=%s", localPath)
		}
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

	koiEnv := "KOI=" + config.Version
	env = append(env, koiEnv)
	l.Debug(koiEnv)

	env = supcolor.UseEnvironColor(env, supcolor.Stderr)

	l.Debugf("PWD=%s", dir)

	l.Debug("Now constructing NodeCmd.")
	cmdPath := filepath.Join(config.Config.InternalNodeExeDir, nodeExe)
	cmdArgs := []string{cmdPath}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Cmd{
		Path: cmdPath,
		Args: cmdArgs,
		Env:  env,
		Dir:  dir,
	}
	nodeCmd := NodeCmd{Cmd: &cmd}

	l.Debug("Now constructing io.")
	if handleStdout {
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			l.Error("Err constructing cmd.StdoutPipe():")
			l.Error(err)
			return nil, err
		}
		stdoutScanner := bufio.NewScanner(stdoutPipe)
		go func() {
			for stdoutScanner.Scan() {
				s := stdoutScanner.Text() + strutil.ResetCtrlStr
				_, _ = os.Stderr.WriteString(s + "\n")
				if nodeCmd.Out != nil {
					*nodeCmd.Out <- NodeCmdOut{
						IsErr: false,
						Text:  s,
					}
				}
			}
			if err := stdoutScanner.Err(); err != nil {
				l.Error("Err reading stdout:")
				l.Error(err)
			}
		}()
	} else {
		cmd.Stdout = os.Stdout
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		l.Error("Err constructing cmd.StderrPipe():")
		l.Error(err)
		return nil, err
	}
	stderrScanner := bufio.NewScanner(stderrPipe)
	go func() {
		for stderrScanner.Scan() {
			s := stderrScanner.Text() + strutil.ResetCtrlStr
			_, _ = os.Stderr.WriteString(s + "\n")
			if nodeCmd.Out != nil {
				*nodeCmd.Out <- NodeCmdOut{
					IsErr: true,
					Text:  s,
				}
			}
		}
		if err := stderrScanner.Err(); err != nil {
			l.Error("Err reading stdout:")
			l.Error(err)
		}
	}()

	return &nodeCmd, nil
}

func (c *NodeCmd) Run() error {
	l.Debug("Now run NodeCmd.")
	// Can use c.Cmd.Run() instead,
	// but remain NodeCmd method call for future refactoring.
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *NodeCmd) Start() error {
	l.Debug("Now start NodeCmd.")
	return c.Cmd.Start()
}

func (c *NodeCmd) Wait() error {
	l.Debug("Now wait NodeCmd.")
	return c.Cmd.Wait()
}
