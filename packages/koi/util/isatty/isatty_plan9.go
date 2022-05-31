package isatty

import "syscall"

func Isatty(fd uintptr) bool {
	path, err := syscall.Fd2path(int(fd))
	if err != nil {
		return false
	}
	return path == "/dev/cons" || path == "/mnt/term/dev/cons"
}
