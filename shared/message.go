package shared

import (
	"encoding/binary"
	"encoding/json"
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

func (self *Message) Send(conn net.Conn) error {
	body, err := json.Marshal(*self)
	if err != nil {
		return err
	}
	bodyLength := uint32(len(body))

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, bodyLength)

	_, err = conn.Write(header)
	if err != nil {
		log.Printf("ERROR: Failed to send header: %v", err)
		return err
	}

	_, err = conn.Write(body)
	if err != nil {
		log.Printf("ERROR: Failed to send body: %v", err)
		return err
	}
	return nil
}
