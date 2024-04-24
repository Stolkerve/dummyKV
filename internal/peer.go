package internal

import (
	"bytes"
	"fmt"
	"net"
	"strings"
)

type Peer struct {
	conn  net.Conn
	state *ServerState
}

func (p *Peer) Start() {
	for {
		// Leer por mensajes del cliente
		payloadSizeBuf := []byte{0, 0}
		if _, err := p.conn.Read(payloadSizeBuf); err != nil {
			fmt.Printf("Conexion de peer %s cerrada\n", p.conn.RemoteAddr().String())
			return
		}

		payloadSize := uint16(payloadSizeBuf[0]) | uint16(payloadSizeBuf[1])<<8

		msgBuf := make([]byte, payloadSize)
		if _, err := p.conn.Read(msgBuf); err != nil {
			fmt.Printf("Conexion de peer %s cerrada\n", p.conn.RemoteAddr().String())
			return
		}

		msg, err := DecodeMsg(bytes.NewBuffer(msgBuf))
		if err != nil {
			p.WriteErr(fmt.Sprintf("Mensaje con formato invalido. %s", err.Error()))
			return
		}

		// Verificar si el contenido
		if msg.Type != MsgTypeArray {
			p.WriteErr("Los mensajes hacia el servidor deben ser tipo array siempre")
			return
		}

		// Leer los argumentos enviados por el cliente
		args := msg.Value.([]Message)
		argsIter := NewArgsIterator(args)
		command, err := argsIter.Next()
		if err != nil {
			p.WriteErr("La lista de argumentos esta vacia")
			return
		}
		if command.Type != MsgTypeString {
			p.WriteErr("El comando debe ser un string")
			return
		}

		commandStr := strings.ToUpper(command.Value.(string))
		fmt.Printf("Comando %s recibido\n", commandStr)

		switch commandStr {
		case PingCommand:
			p.WriteMsg(NewMessage(MsgTypeString, "PONG"))
		case EchoCommand:
			text, err := argsIter.Next()
			if err != nil {
				p.WriteErr("El comando espera un argumento")
				return
			}
			if text.Type != MsgTypeString {
				p.WriteErr("El argumento debe ser un string")
				return
			}
			p.WriteMsg(NewMessage(MsgTypeString, text.Value.(string)))
		case GetCommand:
			key, err := argsIter.Next()
			if err != nil {
				p.WriteErr("El comando espera un argumento")
				return
			}
			if key.Type != MsgTypeString {
				p.WriteErr("La llave debe ser un string")
				return
			}
			p.WriteMsg(NewMessage(MsgTypeNull, nil))

		case SetCommand:
			key, err := argsIter.Next()
			if err != nil {
				p.WriteErr("El comando espera un argumento")
				return
			}
			if key.Type != MsgTypeString {
				p.WriteErr("La llave debe ser un string")
				return
			}
			value, err := argsIter.Next()
			if err != nil {
				p.WriteErr("El segundo comando espera un argumento")
				return
			}
			fmt.Printf("Llave: %s. Valor: %v\n", key.Value.(string), value.Value)
			p.WriteMsg(NewMessage(MsgTypeString, "OK"))

		default:
			p.WriteErr(fmt.Sprintf("Comando %s desconocido", commandStr))
		}
	}
}

func (p *Peer) WriteMsg(msg Message) {
	respMsgBuf, _ := EncodeMsg(msg)
	p.conn.Write(respMsgBuf)
}

func (p *Peer) WriteErr(err string) {
	msg := NewMessage(MsgTypeError, err)
	msgBuf, _ := EncodeMsg(msg)
	p.conn.Write(msgBuf)
}

func NewPeer(conn net.Conn, state *ServerState) *Peer {
	peer := Peer{
		conn:  conn,
		state: state,
	}

	return &peer
}
