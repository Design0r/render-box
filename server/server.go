package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/mitchellh/mapstructure"

	"render-box/server/service"
	"render-box/shared"
	"render-box/shared/db/repo"
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

		go handleConnection(&conn, self.Db)
	}
}

func handleConnection(conn *net.Conn, db *sql.DB) error {
	c := *conn
	defer c.Close()

	var err error
	header := make([]byte, 4)
	for {
		bodySize, err := shared.GetBodySize(conn, header)
		if err != nil {
			break
		}
		body, err := shared.ReadBody(conn, bodySize)
		if err != nil {
			break
		}

		returnData, err := handleMessage(db, body)
		var returnMsg shared.Message
		if err != nil {
			returnMsg = shared.Message{Type: shared.MSGError, Data: nil}
		} else {
			returnMsg = shared.Message{Type: shared.MSGSuccess, Data: returnData}
		}
		fmt.Printf("%+v\n", returnMsg)

		if sendErr := returnMsg.Send(conn); sendErr != nil {
			log.Printf("ERROR: Could not send response to client: %v\n", sendErr)
			break
		}

	}

	log.Printf("Closed connection with %s\n", c.RemoteAddr().String())
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
	case shared.MSGJobsCreate:
		data := &repo.CreateJobParams{}
		mapstructure.Decode(message.Data, data)
		task, err := service.CreateJob(db, data)
		if err != nil {
			return nil, err
		}
		return task, nil
	case shared.MSGJobsAll:
		tasks, err := service.GetJobs(db)
		if err != nil {
			return nil, err
		}
		return tasks, nil

	default:
		return nil, fmt.Errorf("Invalid message type: %v", message.Type)
	}
}
