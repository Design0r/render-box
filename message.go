package main

type MsgType string

const (
	Task MsgType = "Task"
	Job  MsgType = "Job"
)

type Message struct {
	Message MsgType     `json:"message"`
	Data    interface{} `json:"data"`
}
