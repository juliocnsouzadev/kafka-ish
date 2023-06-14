package server

import (
	"net"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/settings"
)

type InternalTcpListener interface {
	Close() error
	Accept(dealine time.Time) (InternalReader, error)
}

type DefaultTcpListener struct {
	tcpListener *net.TCPListener
}

func (t *DefaultTcpListener) Close() error {
	return t.tcpListener.Close()
}

func (t *DefaultTcpListener) Accept(deadline time.Time) (InternalReader, error) {
	err := t.tcpListener.SetDeadline(deadline)
	if err != nil {
		return nil, err
	}

	tcpConn, err := t.tcpListener.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return &DefaultTpcReader{reader: tcpConn}, nil
}

func NewDeafultTcpListener(config settings.Settings) (*DefaultTcpListener, error) {
	tcpListener, err := getTcpListener(config.TcpPort)
	if err != nil {
		return nil, err
	}

	return &DefaultTcpListener{
		tcpListener: tcpListener,
	}, nil
}

func getTcpListener(port string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+port)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", tcpAddr)
}
