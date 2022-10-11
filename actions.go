package goc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func GoogleSetup(c *cli.Context) {
	client := GetClient()
	calList, err := client.CalendarList.List().Do()
	if err != nil {
		log.Fatalf("Unable to retrieve calendar list: %v", err)
	}

	fmt.Println("\nCalendar list:\n--------------")
	for _, elem := range calList.Items {
		fmt.Println(elem.Summary, "  :  ", elem.Id)
	}

	// read user input
	fmt.Print("\nPaste the calendar ID you want to use: ")
	reader := bufio.NewReader(os.Stdin)
	calId, _ := reader.ReadString('\n')
	calId = strings.Replace(calId, "\n", "", -1)

	data := &FileData{
		CalendarId: calId,
	}
	writeToFile(data)

	fmt.Println("You are ready to start tracking!")
}

func StartTask(c *cli.Context) {
	name := "Working"

	if c.NArg() > 0 {
		name = c.Args()[0]
	}

	startTime := c.String("t")
	if startTime == "" {
		startTime = getTime()
	} else {
		startTime = stringToTime(startTime)
	}

	data := readFile()

	if data.CurrentTask.Name != "" {
		newEvent := createEvent(data, getTime())
		event := insertToCalendar(data.CalendarId, newEvent)
		fmt.Println("Previous task added to calendar:", event.HtmlLink)
	}

	data.CurrentTask.Name = name
	data.CurrentTask.Start = startTime
	writeToFile(data)

	fmt.Println("New task started: " + name)
}

func EndTask(c *cli.Context) {
	data := readFile()

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		return
	}

	endTime := c.String("t")
	if endTime == "" {
		endTime = getTime()
	} else {
		endTime = stringToTime(endTime)
	}

	newEvent := createEvent(data, endTime)
	event := insertToCalendar(data.CalendarId, newEvent)
	data.CurrentTask.Reset()
	writeToFile(data)

	fmt.Println("Task added to calendar:", event.HtmlLink)
}

func EditCurrentTask(c *cli.Context) {
	if c.NumFlags() == 0 {
		log.Fatal("Missing at least one flag")
	}

	data := readFile()
	name := c.String("n")
	start := c.String("t")

	if name != "" {
		data.CurrentTask.Name = name
		fmt.Println("New task name set: " + name)
	}
	if start != "" {
		start = stringToTime(start)
		data.CurrentTask.Start = start
		fmt.Println("New start time set: " + start)
	}

	writeToFile(data)
}

func TaskStatus(c *cli.Context) {
	data := readFile()

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		return
	}

	tdata := strings.Split(data.CurrentTask.Start, "T")
	t := fmt.Sprintf("%v %v", tdata[0], strings.Split(tdata[1], "+")[0][:5])
	fmt.Println("Task status:\n------------\nNavn: " + data.CurrentTask.Name + "\nStart: " + t)
}
