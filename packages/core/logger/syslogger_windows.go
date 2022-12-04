package logger

import (
	"fmt"
	"sync"

	"github.com/samber/do"
	"golang.org/x/sys/windows/svc/eventlog"
	"gopkg.ilharper.com/x/rpl"
)

type SysLogger struct {
	c chan *rpl.Log
}

func BuildNewSysLogger() func(i *do.Injector) (*SysLogger, error) {
	return func(i *do.Injector) (*SysLogger, error) {
		wg := do.MustInvoke[*sync.WaitGroup](i)

		_ = eventlog.InstallAsEventCreate(
			logName,
			eventlog.Info,
		)

		e, err := eventlog.Open(logName)
		if err != nil {
			return nil, fmt.Errorf("failed to open event log: %w", err)
		}

		adapter := newColorAdapter(nil)

		sysLogger := &SysLogger{
			c: make(chan *rpl.Log),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				_ = e.Close()
			}()

			for {
				log := <-sysLogger.c
				if log == nil {
					break
				}

				entry := fmt.Sprintf("%04d|%01d|%s", log.Ch, log.Level, adapter.adaptColor(log.Value))
				_ = e.Info(1, entry)
			}
		}()

		return sysLogger, nil
	}
}

func (sysLogger *SysLogger) Writer() chan<- *rpl.Log {
	return sysLogger.c
}

func (sysLogger *SysLogger) Close() {
	sysLogger.c <- nil
}
