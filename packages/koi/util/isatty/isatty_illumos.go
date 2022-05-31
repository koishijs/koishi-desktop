package isatty

import "golang.org/x/sys/unix"

func Isatty(fd uintptr) bool {
	// https://src.illumos.org/source/xref/illumos-gate/usr/src/lib/libc/port/gen/isatty.c
	_, err := unix.IoctlGetTermio(int(fd), unix.TCGETA)
	return err == nil
}
