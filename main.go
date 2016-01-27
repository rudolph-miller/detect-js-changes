package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "detect-js-changes"
	app.Usage = "detects JS changes"
	app.Version = "0.0.1"

	var env string

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "env, e",
			Usage:       "env for cofig file",
			EnvVar:      "ENV",
			Destination: &env,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "detect",
			Usage: "detects changes",
			Action: func(c *cli.Context) {
				println("detect changes")
			},
		},
		{
			Name:  "download",
			Usage: "downloads JS files",
			Action: func(c *cli.Context) {
				println("download JS files")
			},
		},
		{
			Name:  "reset",
			Usage: "resets downloaded JS files",
			Action: func(c *cli.Context) {
				println("reset downloaded JS files")
			},
		},
	}

	app.Run(os.Args)
}
