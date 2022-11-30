package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/mitchellh/go-homedir"
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

		home, err := homedir.Dir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		logPath := filepath.Join(home, "Library/Logs/Koishi.log")

		f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		adapter := newColorAdapter(nil)

		sysLogger := &SysLogger{
			c: make(chan *rpl.Log),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				_ = f.Close()
			}()

			for {
				log := <-sysLogger.c
				if log == nil {
					break
				}

				entry := fmt.Sprintf("%04d|%01d|%s", log.Ch, log.Level, adapter.adaptColor(log.Value))
				_, _ = f.Write([]byte(entry))
				_, _ = f.Write([]byte{'\n'})
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
