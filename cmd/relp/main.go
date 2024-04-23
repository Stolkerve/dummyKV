package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"log"
	"net"

	"github.com/Stolkerve/dummyKV/internal"
)

func main() {
	gob.Register([]internal.Message{})

	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Panicln(err)
	}

	msg := internal.NewMessage(internal.MsgTypeArray, []internal.Message{
		internal.NewMessage(internal.MsgTypeString, "PING"),
	})

	msgBuf, err := internal.EncodeMsg(msg)
	if err != nil {
		log.Panicln(err.Error())
	}
	if _, err := conn.Write(msgBuf); err != nil {
		log.Panicln(err.Error())
	}

	payloadSizeBuf := []byte{0, 0}
	if _, err := conn.Read(payloadSizeBuf); err != nil {
		log.Panicln(err.Error())
	}
	payloadSize := binary.LittleEndian.Uint16(payloadSizeBuf)
	respMsgBuf := make([]byte, payloadSize)
	if _, err := conn.Read(respMsgBuf); err != nil {
		log.Panicln(err.Error())
	}
	respMsg, _ := internal.DecodeMsg(bytes.NewBuffer(respMsgBuf))

	fmt.Println(respMsg)
}
