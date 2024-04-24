package internal

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/Stolkerve/dummyKV/cache"
)

type Peer struct {
	conn  net.Conn
	cache *cache.Cache
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
			v, ok := p.cache.Get(key.Value.(string))
			if !ok {
				p.WriteMsg(NewMessage(MsgTypeNull, nil))
				return
			}
			p.WriteMsg(NewMessage(MsgTypeString, v))

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

			expiration := cache.NO_EXPIRATION
			if expValue, err := argsIter.Next(); err == nil {
				if expValue.Type == MsgTypeString {
					milli, err := strconv.ParseInt(expValue.Value.(string), 10, 64)
					if err != nil {
						p.WriteErr("Duracion no es numero valido")
						return
					}
					expiration = time.Duration(milli) * time.Millisecond
				}
			}

			p.cache.Set(key.Value.(string), value.Value, expiration)
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

func NewPeer(conn net.Conn, cache *cache.Cache) *Peer {
	peer := Peer{
		conn:  conn,
		cache: cache,
	}

	return &peer
}
