package server

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mitchellh/mapstructure"

	"render-box/server/db/repo"
	"render-box/server/service"
	"render-box/shared"
)

type Server struct {
	Addr     string
	Listener *net.Listener
	Db       *sql.DB
}

func NewServer(port string, db *sql.DB) *Server {
	return &Server{Addr: (":" + port), Listener: nil, Db: db}
}

func (self *Server) Run() {
	l, err := net.Listen("tcp", self.Addr)
	if err != nil {
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", self.Addr, err)
	}
	self.Listener = &l
	log.Printf("TCP Server running on %s\n", self.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("ERROR: Could not accept connection: %s\n", err)
			continue
		}
		log.Printf("New connection with %s\n", conn.RemoteAddr().String())

		go handleConnection(conn, self.Db)
	}
}

func getBodySize(conn net.Conn, header []byte) (uint32, error) {
	_, err := io.ReadFull(conn, header)
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

func readBody(conn net.Conn, bodySize uint32) (*shared.Message, error) {
	body := make([]byte, int(bodySize))
	_, err := io.ReadFull(conn, body)
	if err != nil {
		log.Printf("ERROR: Could not read body: %s\n", err)
		return nil, err
	}

	var msg shared.Message
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Printf("ERROR: Could not unmarshall json message: %s\n", err)
		return nil, err
	}

	return &msg, nil
}

func handleConnection(conn net.Conn, db *sql.DB) error {
	defer conn.Close()

	var err error
	header := make([]byte, 4)
	for {
		bodySize, err := getBodySize(conn, header)
		if err != nil {
			break
		}
		body, err := readBody(conn, bodySize)
		if err != nil {
			break
		}

		var returnMsg shared.Message
		returnData, err := handleMessage(db, body)
		if err != nil {
			returnMsg = shared.Message{Type: shared.MSGError, Data: nil}
		} else {
			returnMsg = shared.Message{Type: shared.MSGSuccess, Data: returnData}
		}
		returnMsg.Send(conn)

	}

	log.Printf("Closed connection with %s\n", conn.RemoteAddr().String())
	return err
}

func handleMessage(db *sql.DB, message *shared.Message) (interface{}, error) {
	log.Printf("MESSAGE: %+v\n", message)
	switch message.Type {
	case shared.MSGTasksCreate:
		data := &repo.CreateTaskParams{}
		mapstructure.Decode(message.Data, data)
		task, err := service.CreateTask(db, data)
		if err != nil {
			return nil, err
		}
		return task, nil
	case shared.MSGTasksAll:
		tasks, err := service.GetTasks(db)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	default:
		return nil, fmt.Errorf("Invalid message type: %v", message.Type)
	}
}
