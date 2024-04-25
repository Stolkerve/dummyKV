package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/Stolkerve/dummyKV/internal"
	"github.com/urfave/cli/v2"
)

func SendArgs(conn net.Conn, cliArgs []string) {
	// Enviar lista de argumentos al servidor
	args := make([]internal.Message, len(cliArgs))
	for i := 0; i < len(cliArgs); i++ {
		args[i] = internal.NewMessage(internal.MsgTypeString, cliArgs[i])
	}

	argsMsg := internal.NewMessage(internal.MsgTypeArray, args)
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
	case internal.MsgTypeNull:
		fmt.Println("null")
	}
}

func main() {
	gob.Register([]internal.Message{})

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Value:   "8000",
				Usage:   "Set the server port",
				Aliases: []string{"p"},
			},
			&cli.StringFlag{
				Name:    "address",
				Value:   "0.0.0.0",
				Usage:   "Set the server address",
				Aliases: []string{"addr"},
			},
		},
		Action: func(ctx *cli.Context) error {
			var port, address string
			port = ctx.String("port")
			address = ctx.String("address")

			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
			if err != nil {
				log.Panicln(err)
			}

			cliArgs := os.Args[1:]

			// RELP mode
			if len(cliArgs) == 0 {
				scanner := bufio.NewScanner(os.Stdin)
				for {
					fmt.Print("> ")
					scanner.Scan()
					SendArgs(conn, strings.Fields(scanner.Text()))
				}
			} else {
				SendArgs(conn, cliArgs)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
