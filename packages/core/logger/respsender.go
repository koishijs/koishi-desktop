package logger

import (
	"sync"

	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/x/rpl"
)

type ResponseSender struct {
	c chan *rpl.Log
}

func NewResponseSender(i *do.Injector) (*ResponseSender, error) {
	wg := do.MustInvoke[*sync.WaitGroup](i)

	r := &ResponseSender{
		c: make(chan *rpl.Log),
	}
	// Actually chan<- *proto.Response
	// But do don't support implicit conversion between channels
	ch := do.MustInvoke[chan *proto.Response](i)

	wg.Add(1)
	go func(r1 *ResponseSender, ch1 chan<- *proto.Response) {
		defer wg.Done()

		for {
			log := <-r1.c
			if log == nil {
				break
			}
			ch1 <- proto.NewLog(log)
		}
	}(r, ch)

	return r, nil
}

func (responseSender *ResponseSender) Writer() chan<- *rpl.Log {
	return responseSender.c
}

func (responseSender *ResponseSender) Close() {
	responseSender.c <- nil
}
