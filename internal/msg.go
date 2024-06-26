package internal

import (
	"bytes"
	"encoding/gob"
	"log"
	"math"
)

type MsgType uint8

const (
	MsgTypeString MsgType = iota + 1
	MsgTypeError
	MsgTypeInt
	MsgTypeNull
	MsgTypeArray
)

func (m MsgType) String() string {
	switch m {
	case MsgTypeString:
		return "MsgTypeString"
	case MsgTypeError:
		return "MsgTypeError"
	case MsgTypeInt:
		return "MsgTypeInt"
	case MsgTypeArray:
		return "MsgTypeArray"
	case MsgTypeNull:
		return "MsgTypeNull"
	}
	return ""
}

type Message struct {
	Type  MsgType
	Value interface{}
}

func NewMessage(_type MsgType, value interface{}) Message {
	return Message{
		Type:  _type,
		Value: value,
	}
}

func EncodeMsg(m Message) ([]byte, error) {
	var payLoadBuf bytes.Buffer
	enc := gob.NewEncoder(&payLoadBuf)
	err := enc.Encode(m)
	payloadBytes := payLoadBuf.Bytes()

	if len(payloadBytes) > math.MaxUint16 {
		log.Panic("MSG demasiado largo")
	}

	size := uint16(len(payloadBytes))
	msgSizeBuf := []byte{
		byte(size & 0xFF), byte((size >> 8) & 0xFF),
	}

	// msgBuf
	// -> 0xff, 0xff,  0x01, ..., ..., ...
	// [  Longitud  ] [tipo] [    data    ]
	//                [      mensaje      ]
	var msgBuf bytes.Buffer
	msgBuf.Write(msgSizeBuf)
	msgBuf.Write(payloadBytes)

	return msgBuf.Bytes(), err
}

func DecodeMsg(msgBuf *bytes.Buffer) (*Message, error) {
	dec := gob.NewDecoder(msgBuf)

	msg := new(Message)
	err := dec.Decode(msg)

	return msg, err
}
