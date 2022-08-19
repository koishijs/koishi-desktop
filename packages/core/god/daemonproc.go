package god

import "github.com/samber/do"

type DaemonProc struct{}

func NewDaemonProc(i *do.Injector) (*DaemonProc, error) {
	return &DaemonProc{}, nil
}

// Daemon Koishi processes.
// This will always return an error.
func (daemonProc *DaemonProc) Daemon() error {
	panic(0)
}

func (daemonProc *DaemonProc) Shutdown() error {
	panic(0)
}
