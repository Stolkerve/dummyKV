package internal

import (
	"fmt"
	"log"
	"net"

	"github.com/Stolkerve/dummyKV/cache"
)

type ServerState struct{}

func Listen(address string) {
	cache := cache.NewCache()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Panicln(err.Error())
	}

	fmt.Printf("Server iniciado. Escuchando en: %s\n", address)

	for {
		conn, err := listener.Accept()
		HandleNewConnections(conn, err, &cache)
	}
}

func HandleNewConnections(conn net.Conn, err error, cache *cache.Cache) {
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	peer := NewPeer(conn, cache)
	go peer.Start()
}
