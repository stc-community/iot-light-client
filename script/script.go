package script

import "github.com/urfave/cli"

func Commands() []cli.Command {
	return []cli.Command{
		{
			Name:  "generate",
			Usage: "generate a device",
			Subcommands: []cli.Command{
				{
					Name:   "thermometer",
					Usage:  "generate a thermometer",
					Action: GenerateThermometer,
				},
				{
					Name:   "agv",
					Usage:  "generate a agv",
					Action: GenerateAgv,
				},
			},
		},
	}
}
