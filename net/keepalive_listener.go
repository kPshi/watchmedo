package net

import (
	"net"
	"time"
)

type TcpKeepAliveListener struct {
	*net.TCPListener
}

func Listen(network string, address string) (net.Listener, error) {
	listener, err := net.Listen(network, address)
	if err != nil { return listener, err}

	if tcpListener, ok := listener.(*net.TCPListener); ok {
		return TcpKeepAliveListener{tcpListener}, nil
	}

	return listener, nil
}

func (ln TcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

