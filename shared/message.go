package shared

import (
	"encoding/binary"
	"encoding/json"
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

	MSGTasksCreate MsgType = "tasks.create"
	MSGTasksAll    MsgType = "tasks.all"
)

type Message struct {
	Type MsgType
	Data interface{}
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
