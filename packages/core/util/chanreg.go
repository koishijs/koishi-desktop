package util

type ChannelRegistry[T any] struct {
	reg  [256]chan<- T
	next uint8
}

func (chanReg *ChannelRegistry[T]) Register(c chan<- T) {
	if chanReg.reg[chanReg.next] != nil {
		close(chanReg.reg[chanReg.next])
	}
	chanReg.reg[chanReg.next] = c
	chanReg.next++
}
