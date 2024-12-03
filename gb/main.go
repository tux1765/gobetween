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
				Usage:   "Address for the HTTP server listen on.",
				Value:   "0.0.0.0",
			},
			&cli.StringFlag{
				Name:        "int",
				Aliases:     []string{"i"},
				Usage:       "interface name to listen for network traffic on.",
				DefaultText: "System multicast interface",
			},
			&cli.IntFlag{
				Name:    "buffer",
				Aliases: []string{"b"},
				Usage:   "Buffer size for UDP packets.",
				Value:   1500,
			},
		},
		Action: func(c *cli.Context) error {
			port := c.String("port")
			address := c.String("address")
			listenInterface := c.String("int")
			packetBufferSize := c.Int("buffer")
			fmt.Println("Initializing GoBetween proxy.")
			udpxy.StartServer(address, port, listenInterface, packetBufferSize)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
