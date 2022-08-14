package logger

import (
	"fmt"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/util/strutil"
	"gopkg.ilharper.com/x/rpl"
	"os"
	"strings"
)

type ConsoleTarget struct {
	c     chan *rpl.Log
	Level int8
}

func NewConsoleTarget(i *do.Injector) (*ConsoleTarget, error) {
	targetStream := os.Stderr

	consoleTarget := &ConsoleTarget{
		c:     make(chan *rpl.Log),
		Level: rpl.LevelInfo,
	}

	adapter := newColorAdapter(targetStream)

	go func(ct *ConsoleTarget) {
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
				_, _ = fmt.Fprint(targetStream, outLine)
			}
		}
	}(consoleTarget)

	return consoleTarget, nil
}

func (consoleTarget *ConsoleTarget) Writer() chan<- *rpl.Log {
	return consoleTarget.c
}

func (consoleTarget *ConsoleTarget) Close() {
	consoleTarget.c <- nil
}
