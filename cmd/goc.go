package main

import (
	"log"
	"os"

	"github.com/LeanderGangso/goc"
	"github.com/urfave/cli"
)

var app = cli.NewApp()

func main() {
	info()
	commands()

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	app.Name = "Google Calendar CLI"
	app.Usage = "A simple CLI for tracking hours into Google Calendar"
}

func commands() {
	app.Commands = []cli.Command{
		{
			Name:   "setup",
			Usage:  "Setup Google calendar credentials",
			Action: goc.GoogleAuth,
		},
		{
			Name:      "start",
			Aliases:   []string{"s"},
			Usage:     "Start tracking new task",
			ArgsUsage: "'Task name'",
			Action:    goc.StartTask,
		},
		{
			Name:    "end",
			Aliases: []string{"e"},
			Usage:   "End the currently tracked task",
			Action:  goc.EndTask,
		},
		{
			Name:    "status",
			Aliases: []string{"st"},
			Usage:   "Get current task status.",
			Action:  goc.TaskStatus,
		},
	}
}
