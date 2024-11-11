package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"log"
	"net"

	"render-box/shared"
)

type Server struct {
	Addr     string
	Listener *net.Listener
}

func NewServer(port string) *Server {
	return &Server{Addr: (":" + port), Listener: nil}
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

		go handleConnection(&conn)
	}
}

func getBodySize(conn *net.Conn, header []byte) (uint32, error) {
	_, err := io.ReadFull(*conn, header)
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

func readBody(conn *net.Conn, bodySize uint32) (*shared.Message, error) {
	body := make([]byte, int(bodySize))
	_, err := io.ReadFull(*conn, body)
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

func handleConnection(c *net.Conn) error {
	conn := *c
	defer conn.Close()

	var err error
	header := make([]byte, 4)
	for {
		bodySize, err := getBodySize(&conn, header)
		if err != nil {
			break
		}
		body, err := readBody(&conn, bodySize)
		if err != nil {
			break
		}

		handleMessage(body)
	}

	log.Printf("Closed connection with %s\n", conn.RemoteAddr().String())
	return err
}

func handleMessage(message *shared.Message) {
	log.Printf("MESSAGE: %+v\n", message)
}

func main() {
	Port := "8000"
	server := NewServer(Port)
	server.Run()
}
