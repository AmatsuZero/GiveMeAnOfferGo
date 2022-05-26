package main

import (
	"github.com/AmatsuZero/mycli/src"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	h := src.FindAvailableHost()
	log.Println(h)

	app := &cli.App{
		Commands: []*cli.Command{},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
