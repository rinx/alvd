package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/cli/server"
)

func main() {
	app := &cli.App{
		Name:  "alvd",
		Usage: "A Lightweight Vald",
		Action: func(c *cli.Context) error {
			fmt.Println("boom! I say!")
			return nil
		},
		Commands: []*cli.Command{
			server.NewCommand(),
			agent.NewCommand(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
