package logger

import (
	"fmt"
	"log/syslog"
	"sync"

	"github.com/samber/do"
	"gopkg.ilharper.com/x/rpl"
)

type SysLogger struct {
	c chan *rpl.Log
}

func BuildNewSysLogger() func(i *do.Injector) (*SysLogger, error) {
	return func(i *do.Injector) (*SysLogger, error) {
		var err error

		wg := do.MustInvoke[*sync.WaitGroup](i)

		l, err := syslog.New(syslog.LOG_INFO|syslog.LOG_USER, logName)
		if err != nil {
			return nil, fmt.Errorf("failed to create syslog logger: %w", err)
		}

		adapter := newColorAdapter(nil)

		sysLogger := &SysLogger{
			c: make(chan *rpl.Log),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				_ = l.Close()
			}()

			for {
				log := <-sysLogger.c
				if log == nil {
					break
				}

				entry := fmt.Sprintf("%04d|%01d|%s", log.Ch, log.Level, adapter.adaptColor(log.Value))
				_ = l.Notice(entry)
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
