package daemon

import (
	"bufio"
	"koi/config"
	"koi/util"
	envUtil "koi/util/env"
	l "koi/util/logger"
	"koi/util/strutil"
	"koi/util/supcolor"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type KoiCmdOut struct {
	IsErr bool
	Text  string
}

type KoiCmd struct {
	Cmd *exec.Cmd

	// The output of KoiCmd process.
	//
	// If you set Out to nil, you need to process received KoiCmdOut
	// on a new goroutine. Otherwise, Stdout/Stderr will block.
	Out *chan KoiCmdOut
}

func ResolveYarn() (string, error) {
	yarnPath, err := util.Resolve(config.Config.InternalNodeExeDir, "yarn.cjs")
	if err != nil {
		l.Error("Cannot resolve yarn.")
		return "", err
	}
	return yarnPath, nil
}

func RunYarnCmd(
	args []string,
	dir string,
) error {
	cmd, err := CreateYarnCmd(args, dir)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func CreateYarnCmd(
	args []string,
	dir string,
) (*KoiCmd, error) {
	yarnPath, err := ResolveYarn()
	if err != nil {
		return nil, err
	}
	args = append([]string{yarnPath}, args...)
	return CreateNodeCmd(args, dir, true)
}

func RunNodeCmd(
	args []string,
	dir string,
) error {
	cmd, err := CreateNodeCmd(args, dir, true)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func CreateNodeCmd(
	args []string,
	dir string,
	handleStdout bool,
) (*KoiCmd, error) {
	return CreateKoiCmd(
		config.Config.InternalNodeExeDir,
		"node",
		args,
		dir,
		handleStdout,
	)
}

func CreateKoiCmd(
	path string,
	exe string,
	args []string,
	dir string,
	handleStdout bool,
) (*KoiCmd, error) {
	l.Debug("Getting env.")
	env := os.Environ()

	if config.Config.UseDataHome {
		l.Debug("Now replace HOME/USERPROFILE.")
		envUtil.UseEnv(&env, "HOME", config.Config.InternalHomeDir)
		envUtil.UseEnv(&env, "USERPROFILE", config.Config.InternalHomeDir)
		l.Debugf("HOME=%s", config.Config.InternalHomeDir)

		if runtime.GOOS == "windows" {
			l.Debug("Now replace APPDATA.")
			localPath := filepath.Join(config.Config.InternalHomeDir, "AppData", "Local")
			envUtil.UseEnv(&env, "LOCALAPPDATA", localPath)
			l.Debugf("LOCALAPPDATA=%s", localPath)
		}
	}

	if config.Config.UseDataTemp {
		l.Debug("Now replace TMPDIR/TEMP/TMP.")
		envUtil.UseEnv(&env, "TMPDIR", config.Config.InternalTempDir)
		envUtil.UseEnv(&env, "TEMP", config.Config.InternalTempDir)
		envUtil.UseEnv(&env, "TMP", config.Config.InternalTempDir)
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
		pathEnv = path + pathSepr + pathEnv
	} else {
		pathEnv = path
	}
	env = append(env, "PATH="+pathEnv)
	l.Debugf("PATH=%s", pathEnv)

	koiEnv := "KOI=" + config.Version
	env = append(env, koiEnv)
	l.Debug(koiEnv)

	supcolor.UseColorEnv(&env, supcolor.Stderr)
	config.UseConfigEnv(&env)

	l.Debugf("PWD=%s", dir)

	l.Debug("Now constructing KoiCmd.")
	cmdPath := filepath.Join(path, exe)
	cmdArgs := []string{cmdPath}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Cmd{
		Path: cmdPath,
		Args: cmdArgs,
		Env:  env,
		Dir:  dir,
	}
	koiCmd := KoiCmd{Cmd: &cmd}

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
				if koiCmd.Out != nil {
					*koiCmd.Out <- KoiCmdOut{
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
			if koiCmd.Out != nil {
				*koiCmd.Out <- KoiCmdOut{
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

	return &koiCmd, nil
}

func (c *KoiCmd) Run() error {
	l.Debug("Now run KoiCmd.")
	// Can use c.Cmd.Run() instead,
	// but remain KoiCmd method call for future refactoring.
	if err := c.Start(); err != nil {
		return err
	}
	return c.Wait()
}

func (c *KoiCmd) Start() error {
	l.Debug("Now start KoiCmd.")
	return c.Cmd.Start()
}

func (c *KoiCmd) Wait() error {
	l.Debug("Now wait KoiCmd.")
	return c.Cmd.Wait()
}
