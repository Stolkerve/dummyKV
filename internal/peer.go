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
			break
		}

		payloadSize := uint16(payloadSizeBuf[0]) | uint16(payloadSizeBuf[0])<<8
		msgBuf := make([]byte, payloadSize)
		if _, err := p.conn.Read(msgBuf); err != nil {
			fmt.Printf("Conexion de peer %s cerrada\n", p.conn.RemoteAddr().String())
			break
		}

		msg, err := DecodeMsg(bytes.NewBuffer(msgBuf))
		if err != nil {
			p.WriteErr(fmt.Sprintf("Mensaje con formato invalido. %s", err.Error()))
			break
		}

		// Verificar si el contenido
		if msg.Type != MsgTypeArray {
			p.WriteErr("Los mensajes hacia el servidor deben ser tipo array siempre")
			break
		}

		args := msg.Value.([]Message)
		if len(args) == 0 {
			p.WriteErr("La lista de argumentos esta vacia")
			break
		}

		// Leer los comandos en los argumentos enviados por el cliente
		command := args[0]
		if command.Type != MsgTypeString {
			p.WriteErr("El comando debe ser un string")
			break
		}

		commandStr := strings.ToUpper(command.Value.(string))
		fmt.Printf("Comando %s recibido\n", commandStr)

		switch commandStr {
		case PingCommand:
			p.WriteMsg(NewMessage(MsgTypeString, "PONG"))
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
