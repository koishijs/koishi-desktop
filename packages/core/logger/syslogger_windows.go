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
		var err error

		wg := do.MustInvoke[*sync.WaitGroup](i)

		err = eventlog.InstallAsEventCreate(
			logName,
			eventlog.Info,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to execute EventCreate: %w", err)
		}

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

			var eid uint32 = 1

			for {
				log := <-sysLogger.c
				if log == nil {
					break
				}

				entry := fmt.Sprintf("%04d|%01d|%s", log.Ch, log.Level, adapter.adaptColor(log.Value))
				_ = e.Info(eid, entry)

				eid++
				if eid == 1000 {
					eid = 1
				}
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
