package main

import (
	"log"
	"os"

	"github.com/LeanderGangso/goc"
	"github.com/urfave/cli/v2"
)

const version = "v1.2.1"

func main() {
	app := &cli.App{
		Name:            "goc",
		Usage:           "A simple CLI for tracking hours into Google Calendar",
		Version:         version,
		Suggest:         true,
		HideHelpCommand: true,
		Commands:        commands,
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show help",
	}
	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the version",
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var commands = []*cli.Command{
	{
		Name:   "setup",
		Usage:  "Setup Google calendar",
		Action: goc.GoogleSetup,
	},
	{
		Name:      "start",
		Aliases:   []string{"s"},
		Usage:     "Start tracking new task",
		ArgsUsage: "'name'",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "Set start time for task (HH:MM)",
			},
		},
		Action: goc.StartTask,
	},
	{
		Name:    "end",
		Aliases: []string{"e"},
		Usage:   "End the currently tracked task",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "Set end time for task (HH:MM)",
			},
		},
		Action: goc.EndTask,
	},
	{
		Name:  "edit",
		Usage: "Edit the current task",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Set new task name",
			},
			&cli.StringFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "Set new task time (HH:MM)",
			},
		},
		Action: goc.EditCurrentTask,
	},
	{
		Name:      "insert",
		Aliases:   []string{"i"},
		Usage:     "Insert task directly to calendar",
		ArgsUsage: "'name' start(HH:MM) end(HH:MM)",
		Action:    goc.InsertTask,
	},
	{
		Name:      "alias",
		Aliases:   []string{"a"},
		Usage:     "Add alias to task",
		ArgsUsage: "'alias' 'task'",
		Action:    goc.AddTaskAlias,
	},
	{
		Name:    "remove",
		Aliases: []string{"r"},
		Usage:   "Remove an alias",
		Action:  goc.DelTaskAlias,
	},
	{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all aliases",
		Action:  goc.ShowAlias,
	},
	{
		Name:    "status",
		Aliases: []string{"st"},
		Usage:   "Get current task status",
		Action:  goc.TaskStatus,
	},
}
