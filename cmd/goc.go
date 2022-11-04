package main

import (
	"log"
	"os"

	"github.com/LeanderGangso/goc"
	"github.com/urfave/cli/v2"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	app := &cli.App{
		Name:            "goc",
		Usage:           "A simple CLI for tracking hours into Google Calendar",
		Suggest:         true,
		HideHelpCommand: true,
		Commands:        commands,
	}
	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show help",
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
		ArgsUsage: "['name' | alias 'description']",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "Set start time for task (HHMM)",
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
				Usage:   "Set end time for task (HHMM)",
			},
		},
		Action: goc.EndTask,
	},
	{
		Name:    "update",
		Aliases: []string{"u"},
		Usage:   "Update the current task",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Set new task name",
			},
			&cli.StringFlag{
				Name:    "time",
				Aliases: []string{"t"},
				Usage:   "Set new task time (HHMM)",
			},
		},
		Action: goc.EditCurrentTask,
	},
	{
		Name:      "insert",
		Aliases:   []string{"i"},
		Usage:     "Insert task directly to calendar",
		ArgsUsage: "'name' start(HHMM) end(HHMM)",
		Action:    goc.InsertTask,
	},
	{
		Name:      "alias",
		Aliases:   []string{"a"},
		Usage:     "Add new task alias",
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
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "update",
				Usage: "Get the latest duration info",
			},
			&cli.BoolFlag{
				Name:  "oneline",
				Usage: "List status in oneline format",
			},
			&cli.BoolFlag{
				Name:  "toggle",
				Usage: "Toggle oneline by default",
			},
		},
		Action: goc.TaskStatus,
	},
}
