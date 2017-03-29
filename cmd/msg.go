package cmd

type Type uint8

const (
	Default Type = iota
	Error
)

type Msg struct {
	Type    Type
	Text    string
	Command string
}
