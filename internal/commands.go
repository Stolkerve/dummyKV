package internal

type Command = string

const (
	PingCommand Command = "PING"
	EchoCommand Command = "ECHO"
	GetCommand  Command = "GET"
	SetCommand  Command = "SET"
)
