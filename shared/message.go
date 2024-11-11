package shared

type MsgType string

const (
	TaskCreate MsgType = "task.create"
)

type Message struct {
	Type MsgType
	Data interface{}
}
