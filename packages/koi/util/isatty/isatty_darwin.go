package isatty

import "golang.org/x/sys/unix"

func Isatty(fd uintptr) bool {
	_, err := unix.IoctlGetTermios(int(fd), unix.TIOCGETA)
	return err == nil
}
