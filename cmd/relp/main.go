package main

import (
	"bytes"
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

	// Enviar lista de argumentos al servidor
	argsMsg := internal.NewMessage(internal.MsgTypeArray, []internal.Message{
		internal.NewMessage(internal.MsgTypeString, "PING"),
	})
	argsMsgBuf, err := internal.EncodeMsg(argsMsg)
	if err != nil {
		log.Panicln(err.Error())
	}
	if _, err := conn.Write(argsMsgBuf); err != nil {
		log.Panicln(err.Error())
	}

	// Leer la respuesta del servidor
	payloadSizeBuf := []byte{0, 0}
	if _, err := conn.Read(payloadSizeBuf); err != nil {
		log.Panicln(err.Error())
	}
	payloadSize := uint16(payloadSizeBuf[0]) | uint16(payloadSizeBuf[1])<<8

	respMsgBuf := make([]byte, payloadSize)
	if _, err := conn.Read(respMsgBuf); err != nil {
		log.Panicln(err.Error())
	}
	respMsg, _ := internal.DecodeMsg(bytes.NewBuffer(respMsgBuf))

	// Imprimir la respuesta
	switch respMsg.Type {
	case internal.MsgTypeString:
		fmt.Println(respMsg.Value.(string))
	case internal.MsgTypeError:
		fmt.Printf("ERROR: %s\n", respMsg.Value.(string))
	}
}
