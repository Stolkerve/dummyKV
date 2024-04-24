package internal

import (
	"bytes"
	"fmt"
	"net"
)

type Peer struct {
	conn  net.Conn
	state *ServerState
}

func (p *Peer) Start() {
	for {
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
			msg := NewMessage(MsgTypeError, err.Error())
			msgBuf, _ := EncodeMsg(msg)
			p.conn.Write(msgBuf)
			break
		}

		if msg.Type != MsgTypeArray {
			msg := NewMessage(MsgTypeError, "Los mensajes hacia el servidor deben ser tipo array siempre")
			msgBuf, _ := EncodeMsg(msg)
			p.conn.Write(msgBuf)
		}
		fmt.Println(msg.Value.([]Message))

		respMsg := NewMessage(MsgTypeString, "PONG")
		respMsgBuf, _ := EncodeMsg(respMsg)
		if _, err := p.conn.Write(respMsgBuf); err != nil {
			fmt.Printf("Conexion de peer %s cerrada\n", p.conn.RemoteAddr().String())
			break
		}
	}
}

func NewPeer(conn net.Conn, state *ServerState) *Peer {
	peer := Peer{
		conn:  conn,
		state: state,
	}

	return &peer
}
