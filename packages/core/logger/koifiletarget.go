package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/util/strutil"
	"gopkg.ilharper.com/x/rpl"
)

const (
	ServiceConsoleTarget = "gopkg.ilharper.com/koi/core/logger.ConsoleTarget"
)

type KoiFileTarget struct {
	c     chan *rpl.Log
	Level int8
}

func BuildNewKoiFileTarget(target *os.File) func(i *do.Injector) (*KoiFileTarget, error) {
	return func(i *do.Injector) (*KoiFileTarget, error) {
		wg := do.MustInvoke[*sync.WaitGroup](i)

		consoleTarget := &KoiFileTarget{
			c:     make(chan *rpl.Log),
			Level: rpl.LevelInfo,
		}

		adapter := newColorAdapter(target)

		wg.Add(1)
		go func(ct *KoiFileTarget) {
			defer wg.Done()

			for {
				log := <-ct.c
				if log == nil {
					break
				}

				if log.Level > ct.Level {
					continue
				}

				lines := strings.Split(log.Value, "\n")
				for _, line := range lines {
					outLine := fmt.Sprintf(
						"%s90m%04d|%s%s%s\n",
						strutil.ColorStartCtr,
						log.Ch,
						strutil.ResetCtrlStr,
						line,
						strutil.ResetCtrlStr,
					)
					outLine = adapter.adaptColor(outLine)
					_, _ = fmt.Fprint(target, outLine)
				}
			}
		}(consoleTarget)

		return consoleTarget, nil
	}
}

func (consoleTarget *KoiFileTarget) Writer() chan<- *rpl.Log {
	return consoleTarget.c
}

func (consoleTarget *KoiFileTarget) Close() {
	consoleTarget.c <- nil
}
