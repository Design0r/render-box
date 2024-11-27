package server

import (
	"database/sql"
	"log"
	"net"

	"render-box/server/service"
	"render-box/shared"
)

type Server struct {
	Addr     string
	Listener *net.Listener
	Db       *sql.DB
	Router   *shared.MessageRouter
}

func NewServer(port string, db *sql.DB, router *shared.MessageRouter) *Server {
	return &Server{Addr: (":" + port), Listener: nil, Db: db, Router: router}
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

		go self.handleConnection(&conn)
	}
}

func (self *Server) handleConnection(conn *net.Conn) {
	c := *conn
	defer c.Close()

	state := &shared.ConnState{}

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

		returnData, err := self.Router.Handle(self.Db, body, state)

		var returnMsg shared.Message
		if err != nil {
			returnMsg = shared.Message{Type: shared.MSGError, Data: nil}
		} else {
			returnMsg = shared.Message{Type: shared.MSGSuccess, Data: returnData}
		}
		log.Printf("Return Data %+v\n", returnMsg.Data)

		if sendErr := returnMsg.Send(conn); sendErr != nil {
			log.Printf("ERROR: Could not send response to client: %v\n", sendErr)
			break
		}

	}

	if state.Worker != nil {
		service.UpdateWorkerState(self.Db, "offline", state.Worker.ID)
	}
	if state.Task != nil {
		service.UpdateTaskState(self.Db, "waiting", state.Task.ID)
	}

	log.Printf("Closed connection with %s\n", c.RemoteAddr().String())
}
