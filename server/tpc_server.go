package server

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/producer"
	"github.com/juliocnsouzadev/kafka-ish/settings"
)

type TcpServer struct {
	listener        *net.TCPListener
	commands        chan Command
	commandHandlers map[CommandType]CommandHandler
	prod            producer.Producer
}

func NewTCPServer(config settings.Settings) (*TcpServer, error) {
	listerner, err := getTcpListener(config.TcpPort)
	if err != nil {
		return nil, err
	}

	return &TcpServer{
		listener: listerner,
		commands: make(chan Command),
		prod:     producer.NewProducer(config),
	}, nil
}

func getTcpListener(port string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "localhost:"+port)
	if err != nil {
		return nil, err
	}

	return net.ListenTCP("tcp", tcpAddr)
}

func (t *TcpServer) CommandHandlers() map[CommandType]CommandHandler {
	if t.commandHandlers == nil {
		t.commandHandlers = map[CommandType]CommandHandler{
			CommandPublish: t.publish,
			CommandConsume: t.consume,
			CommandClose:   t.close,
		}
	}
	return t.commandHandlers
}

func (t *TcpServer) publish(command Command) error {
	return nil
}

func (t *TcpServer) consume(command Command) error {
	return nil
}

func (t *TcpServer) close(command Command) error {
	return nil
}

func (t *TcpServer) Cancel() {
	close(t.commands)
	t.prod.Cancel()
	t.listener.Close()
}

func (t *TcpServer) Start() {

	for {
		conn, err := t.connect()
		if err != nil {
			log.Printf("unable to accept tcp connection: %s\n", err)
			continue
		}

		if err := conn.SetKeepAlive(true); err != nil {
			log.Printf("unable to set keep alive: %s\n", err)
		}

		go t.readCommands(conn)

		go t.handleCommands()
	}
}

func (t *TcpServer) handleCommands() {
	for command := range t.commands {
		handler := t.CommandHandlers()[command.Type]
		if handler == nil {
			log.Printf("unknown command: %s\n", command.Type)
			continue
		}

		if err := handler(command); err != nil {
			log.Printf("error on command %s: %s\n", command.Type, err)
		}
	}
}

func (t *TcpServer) readCommands(conn *net.TCPConn) {

	defer conn.Close()
	defer log.Printf("connection closed: %s\n", conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	for {
		line, _, err := reader.ReadLine()

		command := Command{}
		if err == io.EOF {
			command.Type = CommandClose
			t.commands <- command
			return
		}

		if err != nil {
			if closedConnection(err) {
				return
			}

			log.Printf("unable to read connection: %s\n", err)
			continue
		}

		if err = json.Unmarshal(line, &command); err != nil {
			log.Printf("error on json: %s\n", err)
			continue
		}
	}
}

func (t *TcpServer) connect() (*net.TCPConn, error) {
	t.listener.SetDeadline(time.Now().Add(200 * time.Millisecond))
	conn, err := t.listener.AcceptTCP()
	return conn, err
}

func closedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
