//go:build windows

package hideconsole

import (
	"fmt"
	"syscall"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	getConsoleWindow = kernel32.MustFindProc("GetConsoleWindow")
	user32           = syscall.MustLoadDLL("user32.dll")
	showWindowAsync  = user32.MustFindProc("ShowWindowAsync")
)

func HideConsole() error {
	hwnd, _, _ := getConsoleWindow.Call()
	if hwnd == 0 {
		return nil
	}
	_, _, _ = showWindowAsync.Call(hwnd, syscall.SW_HIDE)
	result, _, err := showWindowAsync.Call(hwnd, syscall.SW_HIDE)
	if result != 0 {
		return nil
	}
	return fmt.Errorf("failed to hide console. Last error is: %w", err)
}
