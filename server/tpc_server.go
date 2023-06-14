package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/producer"
)

type TcpServer struct {
	listener        InternalTcpListener
	commands        chan Command
	commandHandlers map[CommandType]CommandHandler
	prod            producer.Producer
	done            chan bool
}

func NewTCPServer(producer producer.Producer, listener InternalTcpListener, done chan bool) (*TcpServer, error) {
	return &TcpServer{
		listener: listener,
		commands: make(chan Command),
		prod:     producer,
		done:     done,
	}, nil
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
	message := model.Message{}

	if err := json.Unmarshal([]byte(command.Body), &message); err != nil {
		return err
	}

	return t.prod.Publish(message)
}

func (t *TcpServer) consume(command Command) error {
	panic("not implemented")
}

func (t *TcpServer) close(command Command) error {
	t.done <- true
	return nil
}

func (t *TcpServer) Cancel() {
	close(t.commands)
	t.prod.Cancel()
	t.listener.Close()
}

func (t *TcpServer) Start() {

	defer t.Cancel()

	for {
		fmt.Println("waiting for connection...")

		deadline := time.Now().Add(200 * time.Millisecond)
		conn, err := t.listener.Accept(deadline)
		if err != nil {
			log.Printf("unable to accept tcp connection: %s\n", err)
			continue
		}

		if err := conn.SetKeepAlive(true); err != nil {
			log.Printf("unable to set keep alive: %s\n", err)
		}

		go t.readCommands(conn)

		go t.handleCommands()

		select {
		case <-t.done:
			fmt.Println("done")
			return
		}
	}
}

func (t *TcpServer) readCommands(conn InternalReader) {

	defer conn.Close()
	defer log.Printf("connection closed: %v\n", conn)

	reader := bufio.NewReader(conn.Reader())
	for {
		line, _, err := reader.ReadLine()

		fmt.Printf("line: %s\n", line)

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
		t.commands <- command
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

		if command.Type == CommandClose {
			t.done <- true
			return
		}
	}
}

func closedConnection(err error) bool {
	return strings.Contains(err.Error(), "use of closed network connection")
}
