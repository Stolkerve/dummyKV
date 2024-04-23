package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		log.Panicln(err)
	}

	conn.Write([]byte("ping"))

	readBuff := make([]byte, 1024)
	conn.Read(readBuff)

	fmt.Println(string(readBuff))
}
