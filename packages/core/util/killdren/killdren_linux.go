//nolint:wrapcheck
package killdren

import (
	"os/exec"
	"syscall"
)

func Set(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGTERM,
	}
}

func Stop(cmd *exec.Cmd) error {
	return Signal(cmd, syscall.SIGTERM)
}

func Kill(cmd *exec.Cmd) error {
	return Signal(cmd, syscall.SIGKILL)
}

func Signal(cmd *exec.Cmd, sig syscall.Signal) error {
	return syscall.Kill(-cmd.Process.Pid, sig)
}
