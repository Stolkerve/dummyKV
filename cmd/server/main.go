package main

import (
	"encoding/gob"

	"github.com/Stolkerve/dummyKV/internal"
)

func main() {
	gob.Register([]internal.Message{})
	internal.Listen("0.0.0.0:8000")
}
