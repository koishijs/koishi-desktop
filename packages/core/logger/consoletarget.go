package logger

import (
	"fmt"
	"github.com/samber/do"
	"gopkg.ilharper.com/x/rpl"
	"strings"
	"sync"
)

type ConsoleTarget struct {
	c     chan *rpl.Log
	Level int8
}

func NewConsoleTarget(i *do.Injector) (*ConsoleTarget, error) {
	wg := do.MustInvoke[*sync.WaitGroup](i)

	consoleTarget := &ConsoleTarget{
		c:     make(chan *rpl.Log),
		Level: rpl.LevelInfo,
	}

	wg.Add(1)
	go func(ct *ConsoleTarget) {
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
				fmt.Printf("%04d|%s\n", log.Ch, line)
			}
		}
	}(consoleTarget)

	return consoleTarget, nil
}

func (consoleTarget *ConsoleTarget) Writer() chan<- *rpl.Log {
	return consoleTarget.c
}

func (consoleTarget *ConsoleTarget) Close() {
	close(consoleTarget.c)
}
