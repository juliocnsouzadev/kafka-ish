package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/juliocnsouzadev/kafka-ish/model"
	"github.com/juliocnsouzadev/kafka-ish/producer"
)

type JsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PublishPayload struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

type HttpServer struct {
	server          *http.Server
	commandHandlers map[CommandType]CommandHandler
	prod            producer.Producer
}

func NewHttpServer(producer producer.Producer, server *http.Server) (*HttpServer, error) {
	return &HttpServer{
		server: server,
		prod:   producer,
	}, nil
}

func (h *HttpServer) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https:*", "http:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))

	mux.Post("/command", h.handle)

	return mux
}

func (h *HttpServer) handle(writer http.ResponseWriter, request *http.Request) {
	command, err := h.readCommand(writer, request)
	if err != nil {
		h.errorJson(writer, err, http.StatusBadRequest)
		log.Println(err)
		return
	}

	handler := h.CommandHandlers()[command.Type]
	if handler == nil {
		h.errorJson(writer, fmt.Errorf("command %s not found", command.Type), http.StatusNotFound)
		log.Println(fmt.Errorf("command %s not found", command.Type))
		return
	}

	err = handler(*command)
	if err != nil {
		h.errorJson(writer, err, http.StatusBadRequest)
		log.Println(fmt.Errorf("error publishing command: %w", err))
		return
	}
}

func (h *HttpServer) readCommand(writer http.ResponseWriter, request *http.Request) (*Command, error) {
	var requestPayload PublishPayload
	err := h.readJson(writer, request, &requestPayload)
	if err != nil {
		return nil, fmt.Errorf("error reading json: %w", err)
	}
	command := Command{
		Type: CommandType(requestPayload.Type),
		Body: requestPayload.Data,
	}
	return &command, nil
}

func (h *HttpServer) readJson(writer http.ResponseWriter, request *http.Request, data interface{}) error {
	maxBytes := 1048576

	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	decode := json.NewDecoder(request.Body)
	err := decode.Decode(data)
	if err != nil {
		return err
	}

	err = decode.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("Body should have one a single JSON")
	}
	return nil
}

func (h *HttpServer) writeJson(writer http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			writer.Header()[key] = value
		}
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_, err = writer.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (h *HttpServer) errorJson(writer http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := JsonResponse{
		Error:   true,
		Message: err.Error(),
	}
	return h.writeJson(writer, statusCode, payload)
}

func (t *HttpServer) CommandHandlers() map[CommandType]CommandHandler {
	if t.commandHandlers == nil {
		t.commandHandlers = map[CommandType]CommandHandler{
			CommandPublish: t.publish,
			CommandConsume: t.consume,
			CommandClose:   t.close,
		}
	}
	return t.commandHandlers
}

func (t *HttpServer) publish(command Command) error {
	message := model.Message{}

	if err := json.Unmarshal([]byte(command.Body), &message); err != nil {
		return err
	}

	return t.prod.Publish(message)
}

func (t *HttpServer) consume(command Command) error {
	panic("not implemented")
}

func (t *HttpServer) close(command Command) error {
	return nil
}

func (t *HttpServer) Cancel() {
	t.prod.Cancel()
	t.server.Close()
}

func (t *HttpServer) Start() {
	http.Handle("/", t.routes())
	t.server.ListenAndServe()
}
