package logger

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"github.com/samber/do"
	"gopkg.ilharper.com/koi/core/god/proto"
	"gopkg.ilharper.com/x/rpl"
)

func FilterLog(resp <-chan *proto.Response) (<-chan *rpl.Log, <-chan *proto.Response) {
	if resp == nil {
		panic("koi/core/logger/filter: response channel is nil")
	}

	log := make(chan *rpl.Log)
	data := make(chan *proto.Response)

	go func() {
		for {
			r := <-resp
			if r == nil {
				log <- nil
				data <- nil

				break
			}

			if r.Type != proto.TypeResponseLog {
				data <- r
			} else {
				l := rpl.Log{}
				err := mapstructure.Decode(r.Data, &l)
				if err != nil {
					// Normally there won't be error here.
					// If your websocket isn't that stable,
					// fill an issue and I'll remove this panic.
					// This will not introduce a new major version so please treat carefully.
					panic(fmt.Errorf("koi/core/logger/filter: failed to decode response: %w", err))
				}
				log <- &l
			}
		}
	}()

	return log, data
}

func LogChannel(i *do.Injector, logC <-chan *rpl.Log) {
	receiver := do.MustInvoke[*rpl.Receiver](i)

	if logC == nil {
		panic("koi/core/logger/logchannel: log channel is nil")
	}

	go func() {
		for {
			log := <-logC
			if log == nil {
				break
			}
			receiver.Writer() <- log
		}
	}()
}

// Wait for a nil to ensure communication ended.
func Wait(ch <-chan *proto.Response) error {
	if <-ch != nil {
		return errors.New("got a non-nil response")
	}

	return nil
}
