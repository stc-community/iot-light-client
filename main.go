package main

import (
	"os"

	"github.com/stc-community/iot-light-client/ignite"

	"github.com/urfave/cli"
	_ "go.uber.org/automaxprocs"
)

func main() {
	app := cli.NewApp()
	app.Name = "iot-light-client"
	app.Version = "1.0.0"
	app.Author = "anonymous"
	app.Action = func(c *cli.Context) error {
		ignite.Event()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic("app run error:" + err.Error())
	}
}
