package hub

type MessageType uint8

const (
	MessageTypeNotification MessageType = iota
	MessageTypeError
	MessageTypeSuccess
)

type Message struct {
	Type    MessageType
	Payload interface{}
}
