package proc

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/logger"
	"gopkg.ilharper.com/x/rpl"
)

type KoiProc struct {
	i  *do.Injector
	ch uint16

	cmd *exec.Cmd

	logTargets []rpl.Target

	HookOutput func(output string)
}

func NewKoiProc(
	i *do.Injector,
	ch uint16,
	path string,
	command string,
	args []string,
	cwd string,
) *KoiProc {
	cmdPath := filepath.Join(path, command)
	cmdArgs := append([]string{cmdPath}, args...)
	env := environ(i, path)

	return &KoiProc{
		i:  i,
		ch: ch,

		cmd: &exec.Cmd{
			Path: cmdPath,
			Args: cmdArgs,
			Env:  *env,
			Dir:  cwd,
		},
	}
}

func (koiProc *KoiProc) Register(target rpl.Target) {
	koiProc.logTargets = append(koiProc.logTargets, target)
}

func (koiProc *KoiProc) Close() {
	panic(errors.New("no need to call Close(). Channel will close automatically after subprocess dead"))
}

func (koiProc *KoiProc) Run() error {
	var err error

	l := do.MustInvoke[*logger.Logger](koiProc.i)

	// Make output channel
	out := make(chan *string)
	defer close(out)

	// Setup log targets
	go func() {
		for {
			str := <-out
			if str == nil {
				break
			}

			if koiProc.HookOutput != nil {
				koiProc.HookOutput(*str)
			}

			log := &rpl.Log{
				Ch:    koiProc.ch,
				Level: rpl.LevelInfo,
				Value: *str,
			}

			for _, target := range koiProc.logTargets {
				go func(t rpl.Target, l *rpl.Log) {
					t.Writer() <- l
				}(target, log)
			}
		}
	}()

	// Setup IO pipes
	stdoutPipe, err := koiProc.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := koiProc.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}
	scanners := []*bufio.Scanner{
		bufio.NewScanner(stdoutPipe),
		bufio.NewScanner(stderrPipe),
	}
	for _, scanner := range scanners {
		go func(scn *bufio.Scanner) {
			for {
				if !scn.Scan() {
					break
				}
				scnErr := scn.Err()
				if scnErr != nil {
					l.Warn(fmt.Errorf("koiProc scanner err: %w", scnErr))
				} else {
					s := scn.Text()
					out <- &s
				}
			}
		}(scanner)
	}

	// Run process
	err = koiProc.cmd.Run()
	if err != nil {
		// Here err is likely to be an ExitError,
		// Which is normal (killed by god daemon).
		// No need to wrap this error.
		return err //nolint:wrapcheck
	}

	return nil
}

// Stop sends [syscall.SIGTERM] to process.
//
// This just sends the signal and do not wait for anything.
func (koiProc *KoiProc) Stop() error {
	if koiProc.cmd.Process == nil {
		return nil
	}

	err := koiProc.cmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to send SIGTERM to process: %w", err)
	}

	return nil
}

// Kill sends [syscall.SIGKILL] to process.
//
// This just sends the signal and do not wait for anything.
// If possible, use [KoiProc.Stop].
func (koiProc *KoiProc) Kill() error {
	if koiProc.cmd.Process == nil {
		return nil
	}

	err := koiProc.cmd.Process.Kill()
	if err != nil {
		return fmt.Errorf("failed to send SIGKILL to process: %w", err)
	}

	return nil
}
