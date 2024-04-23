package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Panicln(err.Error())
	}

	fmt.Println("Server iniciado. Escuchando en: 0.0.0.0:8000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicln(err.Error())
			continue
		}
		readBuff := make([]byte, 1024)
		conn.Read(readBuff)
		fmt.Println(string(readBuff))

		conn.Write([]byte("pong"))
	}
}
