package shared

import (
	"log"
	"net"
	"render-box/shared/db/repo"
)

type ConnState struct {
	Worker *repo.Worker
	Task   *repo.Task
}

type TCPListener struct {
	Port string
	Conn *net.Conn
}

func NewTcpListener(port string) *TCPListener {
	return &TCPListener{Port: port, Conn: nil}
}

func (self *TCPListener) Run() (*net.Conn, error) {
	conn, err := net.Dial("tcp", ":"+self.Port)
	if err != nil {
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", self.Port, err)
		return nil, err
	}
	self.Conn = &conn
	log.Printf("New connection with %s\n", conn.RemoteAddr().String())

	return &conn, nil
}
