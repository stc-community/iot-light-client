package main

import (
	"os"

	"github.com/stc-community/iot-light-client/event"
	"github.com/stc-community/iot-light-client/script"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "iot-light-client"
	app.Version = "1.0.0"
	app.Author = "anonymous"
	//app.Before = Before
	app.Commands = script.Commands()
	app.Action = func(c *cli.Context) error {
		event.AcceptEvent()
		event.CatchCmdSignals()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		panic("app run error:" + err.Error())
	}
}
