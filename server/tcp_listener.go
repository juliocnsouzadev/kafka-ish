package server

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/settings"
)

type InternalTcpListener interface {
	Close() error
	Accept(deadline time.Time) (InternalReader, *TcpAcceptError)
}

type TcpAcceptError struct {
	err     error
	label   string
	timeout bool
}

func (i *TcpAcceptError) IsTimeout() bool {
	return i.timeout
}

func (i *TcpAcceptError) Error() string {
	return fmt.Sprintf("TcpAcceptError, %s : %v", i.label, i.err)
}

type DefaultTcpListener struct {
	tcpListener *net.TCPListener
}

func (t *DefaultTcpListener) Close() error {
	return t.tcpListener.Close()
}

func (t *DefaultTcpListener) Accept(deadline time.Time) (InternalReader, *TcpAcceptError) {
	err := t.tcpListener.SetDeadline(deadline)
	if err != nil {
		acceptError := &TcpAcceptError{err: err, label: "error setting deadline"}
		return nil, acceptError
	}

	tcpConn, err := t.tcpListener.AcceptTCP()
	if err != nil {
		acceptError := &TcpAcceptError{err: err, label: "error accepting connection"}
		if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
			acceptError.timeout = true
		}
		return nil, acceptError
	}
	return &DefaultTpcReader{reader: tcpConn}, nil
}

func NewDefaultTcpListener(config settings.Settings) (*DefaultTcpListener, error) {
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
		return nil, fmt.Errorf("ResolveTCPAddr failed %w", err)
	}

	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("ListenTCP failed %w", err)
	}
	log.Printf("tcp listening on port %s\n", port)
	return tcpListener, nil
}
