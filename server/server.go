package server

const (
	CommandPublish CommandType = "publish"
	CommandConsume CommandType = "consume"
	CommandClose   CommandType = "close"
)

type CommandHandler func(command Command) error

type CommandType string

type Command struct {
	Type CommandType `json:"type"`
	Body string      `json:"body"`
}

type Server interface {
	Start()
	Cancel()
	CommandHandlers() map[CommandType]CommandHandler
}
