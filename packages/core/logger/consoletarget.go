package logger

import (
	"fmt"
	"gopkg.ilharper.com/x/rpl"
	"strings"
)

type ConsoleTarget struct {
	c     chan rpl.Log
	Level int8
}

func NewConsoleTarget() *ConsoleTarget {
	consoleTarget := &ConsoleTarget{
		Level: rpl.LevelInfo,
	}

	go func(ct *ConsoleTarget) {
		for {
			log := <-ct.c
			if log.Level > ct.Level {
				continue
			}

			lines := strings.Split(log.Value, "\n")
			for _, line := range lines {
				fmt.Printf("%04d|%s\n", log.Ch, line)
			}
		}
	}(consoleTarget)

	return consoleTarget
}

func (consoleTarget *ConsoleTarget) Writer() chan<- rpl.Log {
	return consoleTarget.c
}
