package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/Stolkerve/dummyKV/internal"
	"github.com/urfave/cli/v2"
)

func main() {
	gob.Register([]internal.Message{})

	var port, address string

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
			port = ctx.String("port")
			address = ctx.String("address")

			internal.Listen(fmt.Sprintf("%s:%s", address, port))
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
