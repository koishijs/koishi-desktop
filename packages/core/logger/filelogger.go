package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/koiconfig"
	"gopkg.ilharper.com/x/rpl"
)

// FileLogger defaults.
const (
	// logFormat known as "yyyy-MM-dd".
	// See [package time].
	//
	// [package time]: https://pkg.go.dev/time#pkg-constants
	logFormat = "2006-01-02"

	maxAge = 30 * 24 * time.Hour
)

type FileLogger struct {
	c chan *rpl.Log
}

func BuildNewFileLogger() func(i *do.Injector) (*FileLogger, error) {
	return func(i *do.Injector) (*FileLogger, error) {
		var err error

		cfg := do.MustInvoke[*koiconfig.Config](i)
		wg := do.MustInvoke[*sync.WaitGroup](i)

		now := time.Now()
		date := now.Format(logFormat)

		go func() {
			var rErr error
			files, rErr := os.ReadDir(cfg.Computed.DirLog)
			if rErr != nil {
				return
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				name := file.Name()
				if len(name) < 14 {
					continue
				}

				fileDate, rErr := time.Parse(logFormat, name[:10])
				if rErr != nil {
					continue
				}

				if fileDate.Add(maxAge).Before(now) {
					_ = os.Remove(filepath.Join(cfg.Computed.DirLog, name))
				}
			}
		}()

		path := filepath.Join(cfg.Computed.DirLog, fmt.Sprintf("%s.log", date))
		f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}

		adapter := newColorAdapter(nil)

		fileLogger := &FileLogger{
			c: make(chan *rpl.Log),
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				_ = f.Close()
			}()

			for {
				log := <-fileLogger.c
				if log == nil {
					break
				}

				entry := fmt.Sprintf("%04d|%01d|%s", log.Ch, log.Level, adapter.adaptColor(log.Value))
				_, _ = f.Write([]byte(entry))
				_, _ = f.Write([]byte{'\n'})
			}
		}()

		return fileLogger, nil
	}
}

func (fileLogger *FileLogger) Writer() chan<- *rpl.Log {
	return fileLogger.c
}

func (fileLogger *FileLogger) Close() {
	fileLogger.c <- nil
}
