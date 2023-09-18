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
					Name:   "avg",
					Usage:  "generate a avg",
					Action: GenerateAvg,
				},
			},
		},
	}
}
