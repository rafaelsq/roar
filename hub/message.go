package hub

type MessageType uint8

const (
	MessageTypeNotification MessageType = iota
)

type Message struct {
	Type    MessageType
	Payload interface{}
}
