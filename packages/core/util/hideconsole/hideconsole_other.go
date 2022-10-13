//go:build !windows

package hideconsole

func HideConsole() error {
	return nil
}
