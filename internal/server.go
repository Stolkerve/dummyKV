package internal

import (
	"fmt"
	"log"
	"net"
)

type ServerState struct{}

func Listen(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicln(err.Error())
	}

	fmt.Printf("Server iniciado. Escuchando en: %s\n", address)

	for {
		HandleNewConnections(listener.Accept())
	}
}

func HandleNewConnections(conn net.Conn, err error) {
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	peer := NewPeer(conn, &ServerState{})
	go peer.Start()
}
