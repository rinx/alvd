package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"

	"github.com/rinx/alvd/pkg/alvd/cli/agent"
	"github.com/rinx/alvd/pkg/alvd/cli/server"
)

var (
	Version    = "unknown"
	GoVersion  = "unknown"
	NGTVersion = "unknown"
)

func main() {
	cli.VersionPrinter = versionPrinter

	app := &cli.App{
		Name:  "alvd",
		Usage: "A Lightweight Vald",
		Commands: []*cli.Command{
			server.NewCommand(),
			agent.NewCommand(),
		},
		Version: Version,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func versionPrinter(c *cli.Context) {
	fmt.Printf("alvd Version: %s\n", c.App.Version)
	fmt.Printf("Build Go Version: %s\n", GoVersion)
	fmt.Printf("NGT Version: %s\n", NGTVersion)
}
