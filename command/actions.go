package command

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

func StartTask(c *cli.Context) {
	if c.NArg() == 0 {
		log.Fatal("Missing required task name")
	}

	data := readFile()

	if data.CurrentTask.Name != "" {
		newEvent := createEvent(data)
		insertToCalendar(data.CalendarId, newEvent)
	}

	name := c.Args()[0]
	data.CurrentTask.Name = name
	data.CurrentTask.Start = getTime()
	writeToFile(data)

	fmt.Println("New task started: " + name)
}

func EndTask(c *cli.Context) {
	data := readFile()

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		return
	}

	newEvent := createEvent(data)
	insertToCalendar(data.CalendarId, newEvent)
	fmt.Println("Old task ended: " + data.CurrentTask.Name)
}

func TaskStatus(c *cli.Context) {
	data := readFile()

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		return
	}

	fmt.Println("Task status:\n------------\nNavn: " + data.CurrentTask.Name + "\nStart: " + data.CurrentTask.Start)
}
