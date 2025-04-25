package shared

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
)

type MsgType string

const (
	MSGSuccess MsgType = "success"
	MSGError   MsgType = "error"

	MSGJobsCreate MsgType = "jobs.create"
	MSGJobsAll    MsgType = "jobs.all"

	MSGTasksNext     MsgType = "tasks.next"
	MSGTasksCreate   MsgType = "tasks.create"
	MSGTasksAll      MsgType = "tasks.all"
	MSGTasksComplete MsgType = "tasks.complete"

	MSGWorkerRegister MsgType = "worker.register"
	MSGWorkerCreate   MsgType = "worker.create"
	MSGWorkerAll      MsgType = "worker.all"
)

type Message struct {
	Type MsgType
	Data interface{}
}

type RawMessage struct {
	Type MsgType
	Data json.RawMessage
}

func (self *Message) Send(conn *net.Conn) error {
	c := *conn
	body, err := json.Marshal(*self)
	if err != nil {
		return err
	}
	bodyLength := uint32(len(body))

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, bodyLength)

	_, err = c.Write(header)
	if err != nil {
		log.Printf("ERROR: Failed to send header: %v", err)
		return err
	}

	_, err = c.Write(body)
	if err != nil {
		log.Printf("ERROR: Failed to send body: %v", err)
		return err
	}
	return nil
}

func GetBodySize(conn *net.Conn, header []byte) (uint32, error) {
	c := *conn
	_, err := io.ReadFull(c, header)
	if err != nil {
		if err == io.EOF {
			log.Println("Connection closed by the server.")
			return 0, err
		}
		log.Printf("ERROR: Could not read header: %s\n", err)
		return 0, err
	}

	bodyLength := binary.BigEndian.Uint32(header)
	return bodyLength, nil
}

func ReadBody(conn *net.Conn, bodySize uint32) (*Message, error) {
	body := make([]byte, int(bodySize))
	_, err := io.ReadFull(*conn, body)
	if err != nil {
		log.Printf("ERROR: Could not read body: %s\n", err)
		return nil, err
	}

	msg := Message{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Printf("ERROR: Could not unmarshall json message: %s\n", err)
		return nil, err
	}

	return &msg, nil
}

func UnmarshallBody[T any](body any) (*T, error) {
	bodyData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	var generic T
	err = json.Unmarshal(bodyData, &generic)
	if err != nil {
		return nil, err
	}

	return &generic, nil
}

func RecvMessage[T any](conn *net.Conn) (*Message, error) {
	header := make([]byte, 4)
	bodySize, err := GetBodySize(conn, header)
	if err != nil {
		return nil, err
	}
	body, err := ReadBody(conn, bodySize)
	if err != nil {
		return nil, err
	}

	generic, err := UnmarshallBody[T](body.Data)
	if err != nil {
		return nil, err
	}

	return &Message{Type: body.Type, Data: *generic}, nil
}

type (
	MessageHandler = func(db *sql.DB, message *Message, state *ConnState) (any, error)
	MessageRouter  struct {
		Routes map[string]MessageHandler
	}
)

func NewMessageRouter() *MessageRouter {
	m := map[string]MessageHandler{}
	return &MessageRouter{Routes: m}
}

func (self *MessageRouter) IncludeRouter(router *MessageRouter) {
	for k, v := range router.Routes {
		if _, exists := self.Routes[k]; exists {
			log.Printf("Failed to include route: %v, route already exists", k)
			continue
		}
		self.Routes[k] = v
	}
}

func (self *MessageRouter) Register(path string, handler MessageHandler) {
	if _, exists := self.Routes[path]; exists {
		log.Printf("Failed to register route: %v, route already exists", path)
		return
	}
	self.Routes[path] = handler
}

func (self *MessageRouter) Handle(
	db *sql.DB,
	message *Message,
	state *ConnState,
) (any, error) {
	handler, exists := self.Routes[string(message.Type)]
	if !exists {
		return nil, fmt.Errorf("Failed to add route: %v, route already exists", message.Type)
	}

	return handler(db, message, state)
}
