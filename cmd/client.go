package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"render-box/shared"
)

func handleRead(conn net.Conn) {
	buffer := make([]byte, 512)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("ERROR: Could not read from connection: %s\n", err)
			break
		}

		log.Printf("MESSAGE: %s\n", string(buffer[:n]))
	}
}

func sendMessage(conn net.Conn, message string) error {
	// Step 1: Convert the message to bytes
	body := []byte(message)
	bodyLength := uint32(len(body))

	// Step 2: Create a 4-byte header with the body length
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, bodyLength)

	// Step 3: Write the header followed by the body
	_, err := conn.Write(append(header, body...))
	if err != nil {
		return fmt.Errorf("could not send message: %w", err)
	}
	return nil
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

func handleWrite(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		sendMessage(conn, message)
	}
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
	// handleWrite(conn)

	task := shared.CreateTask{Name: "Peter", Age: 32}
	msg := shared.Message{Type: shared.TaskCreate, Data: &task}
	sendJsonMessage(conn, &msg)
}
