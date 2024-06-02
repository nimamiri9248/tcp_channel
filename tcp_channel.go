package tcpchannel

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

type TCPChannel[T any] struct {
	listenAddress string
	remoteAddress string
	sendChan      chan T
	recvChan      chan T
	outboundConn  net.Conn
	listener      net.Listener
}

func NewTCPChannel[T any](listenAddress, remoteAddress string) (*TCPChannel[T], error) {
	tcpChannel := &TCPChannel[T]{
		listenAddress: listenAddress,
		remoteAddress: remoteAddress,
		sendChan:      make(chan T, 10),
		recvChan:      make(chan T, 10),
	}

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		return nil, err
	}
	tcpChannel.listener = ln

	go tcpChannel.sendLoop()
	go tcpChannel.acceptConnections()
	go tcpChannel.connectToRemote()

	return tcpChannel, nil
}

func (tc *TCPChannel[T]) sendLoop() {
	for {
		msg := <-tc.sendChan
		log.Println("Sending message over the wire:", msg)
		if err := gob.NewEncoder(tc.outboundConn).Encode(&msg); err != nil {
			log.Println("Encoding error:", err)
		}
	}
}

func (tc *TCPChannel[T]) acceptConnections() {
	defer func() {
		tc.listener.Close()
	}()

	for {
		conn, err := tc.listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			return
		}

		log.Printf("Sender connected: %s", conn.RemoteAddr())

		go tc.handleConnection(conn)
	}
}

func (tc *TCPChannel[T]) handleConnection(conn net.Conn) {
	for {
		var msg T
		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Println("Decoding error:", err)
			continue
		}
		tc.recvChan <- msg
	}
}

func (tc *TCPChannel[T]) connectToRemote() {
	for {
		conn, err := net.Dial("tcp", tc.remoteAddress)
		if err != nil {
			log.Printf("Dial error: %s. Retrying in 3 seconds...", err)
			time.Sleep(time.Second * 3)
			continue
		}

		tc.outboundConn = conn
		return
	}
}
