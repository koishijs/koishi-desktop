package killdren

import (
	"os"
	"os/exec"
	"syscall"
)

func Set(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func Stop(cmd *exec.Cmd) error {
	return Signal(cmd, syscall.SIGTERM)
}

func Kill(cmd *exec.Cmd) error {
	return Signal(cmd, syscall.SIGKILL)
}

func Signal(cmd *exec.Cmd, sig os.Signal) error {
	syscall.Kill(-cmd.Process.Pid, sig)
}
