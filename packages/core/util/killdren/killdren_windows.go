package killdren

import (
	"os"
	"os/exec"
	"syscall"
)

func Set(cmd *exec.Cmd) {
}

func Stop(cmd *exec.Cmd) error {
	return cmd.Process.Signal(syscall.SIGTERM)
}

func Kill(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}

func Signal(cmd *exec.Cmd, sig os.Signal) error {
	return cmd.Process.Signal(sig)
}
