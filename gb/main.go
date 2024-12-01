package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tux1765/gobetween/udpxy"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gb",
		Usage: "Proxy UDP MPEGTS packets through HTTP",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Port to run the http server on.",
				Value:   "4001",
			},
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Address to listen on.",
				Value:   "0.0.0.0",
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			address := c.String("address")
			fmt.Println("Initializing GoBetween proxy.")
			udpxy.StartServer(address, port)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
