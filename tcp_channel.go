package tcpchannel

import "net"

type Receiver[T any] struct {
	Chan chan T
	listenAddress string
	listener net.Listener
}

func NewReceiver[T any](listenAddress string) (*Receiver[T], error) {
	receiver := &Receiver[T]{
		Chan: make(chan T),
		listenAddress: listenAddress,
	}

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}

	receiver.listener = ln


}

type Sender struct {
}