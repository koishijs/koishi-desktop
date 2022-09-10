package tray

import (
	"sync"

	"github.com/samber/do"
)

const serviceTrayChannelRegistry = "gopkg.ilharper.com/koi/app/ui/tray.ChannelRegistry"

type ChannelRegistry struct {
	lock  sync.Mutex
	reg   [256]chan struct{}
	index uint8
}

func NewChannelRegistry(i *do.Injector) (*ChannelRegistry, error) {
	return &ChannelRegistry{}, nil
}

func (cr *ChannelRegistry) Insert(ch chan struct{}) {
	cr.lock.Lock()
	defer cr.lock.Unlock()

	if cr.reg[cr.index] != nil {
		close(cr.reg[cr.index])
	}
	cr.reg[cr.index] = ch
	cr.index++
}
