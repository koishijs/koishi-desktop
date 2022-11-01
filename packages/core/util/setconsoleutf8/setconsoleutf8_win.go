//go:build windows

package setconsoleutf8

import "syscall"

func SetConsoleUTF8() error {
	var err error
	var result uintptr = 0

	k32 := syscall.MustLoadDLL("kernel32.dll")
	result, _, err = k32.MustFindProc("SetConsoleCP").Call(65001)
	if result == 0 {
		return err
	}
	result, _, err = k32.MustFindProc("SetConsoleOutputCP").Call(65001)
	if result == 0 {
		return err
	}

	return nil
}
