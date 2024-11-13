package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"render-box/shared"
	"render-box/shared/db/repo"
)

func handleRead(conn net.Conn) {
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

		log.Printf("MESSAGE: %v\n", body)
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

	msg := shared.Message{}
	err = json.Unmarshal(body, &msg)
	if err != nil {
		log.Printf("ERROR: Could not unmarshall json message: %s\n", err)
		return nil, err
	}

	return &msg, nil
}

func sendJsonMessage(conn net.Conn, msg *shared.Message) error {
	body, err := json.Marshal(*msg)
	if err != nil {
		return err
	}
	bodyLength := uint32(len(body))

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, bodyLength)

	_, err = conn.Write(append(header, body...))
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}
	return nil
}

func main() {
	Port := "8000"
	conn, err := net.Dial("tcp", ":"+Port)
	if err != nil {
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", Port, err)
	}
	defer conn.Close()
	log.Printf("New connection with %s\n", conn.RemoteAddr().String())

	go handleRead(conn)

	task := repo.CreateJobParams{Name: "TestJob", Priority: int64(50), State: "waiting"}
	msg := shared.Message{Type: shared.MSGJobsCreate, Data: &task}
	sendJsonMessage(conn, &msg)
}
