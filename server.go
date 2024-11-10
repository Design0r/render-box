package main

import (
	"encoding/binary"
	"io"
	"log"
	"net"
)

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

func handleConnection(c *net.Conn) {
	conn := *c
	defer conn.Close()

	header := make([]byte, 4)
	for {
		bodySize, err := getBodySize(conn, header)
		if err != nil {
			break
		}

		body := make([]byte, int(bodySize))
		_, err = io.ReadFull(conn, body)
		if err != nil {
			log.Printf("ERROR: Could not read body: %s\n", err)
			break
		}

		log.Printf("MESSAGE: %s\n", string(body))
	}

	log.Printf("Closed connection with %s\n", conn.RemoteAddr().String())
}

func main() {
	Port := "8000"
	l, err := net.Listen("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", Port, err)
	}
	log.Printf("TCP Server running on port %s\n", Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("ERROR: Could not accept connection: %s\n", err)
		}
		log.Printf("New connection with %s\n", conn.RemoteAddr().String())

		go handleConnection(&conn)
	}
}
