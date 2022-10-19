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
	if c.NArg() < 1 {
		log.Fatal("Missing required argument")
	}

	data := readFile()
	name := checkAndUseAlias(c.Args()[0], data)

	startTime := c.String("t")
	if startTime == "" {
		startTime = getTime()
	} else {
		startTime = stringToTime(startTime)
	}

	if data.CurrentTask.Name != "" {
		newEvent := createEvent(data, getTime())
		event := insertToCalendar(data.CalendarId, newEvent)
		updatePrevTaskAlias(data.CurrentTask.Name, data)
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
	updatePrevTaskAlias(data.CurrentTask.Name, data)
	data.CurrentTask.Reset()
	writeToFile(data)

	fmt.Println("Task added to calendar:", event.HtmlLink)
}

func EditCurrentTask(c *cli.Context) {
	if c.NumFlags() == 0 {
		log.Fatal("Missing at least one flag")
	}

	data := readFile()
	name := checkAndUseAlias(c.String("n"), data)
	start := c.String("t")

	if name != "" {
		data.CurrentTask.Name = name
		fmt.Println("New task name set: " + name)
	}
	if start != "" {
		start = stringToTime(start)
		data.CurrentTask.Start = start
		fmt.Println("New start time set: " + formatTimeString(start))
	}

	writeToFile(data)
}

func InsertTask(c *cli.Context) {
	if c.NArg() < 3 {
		log.Fatal("Missing required arguments")
	}

	data := readFile()
	name := checkAndUseAlias(c.Args()[0], data)
	startTime := stringToTime(c.Args()[1])
	endTime := stringToTime(c.Args()[2])

	data.CurrentTask.Name = name
	data.CurrentTask.Start = startTime

	newEvent := createEvent(data, endTime)
	event := insertToCalendar(data.CalendarId, newEvent)

	fmt.Println("Task addded to calendar:", event.HtmlLink)
}

func AddTaskAlias(c *cli.Context) {
	if c.NArg() < 2 {
		log.Fatal("Missing required arguments")
	}

	taskName := c.Args()[0]
	aliasName := c.Args()[1]

	data := readFile()
	if data.TaskAlias == nil {
		data.TaskAlias = make(map[string]string)
	}
	data.TaskAlias[aliasName] = taskName
	writeToFile(data)

	fmt.Println("Alias added:", aliasName+": "+taskName)
}

func DelTaskAlias(c *cli.Context) {
	if c.NArg() < 1 {
		log.Fatal("Missing required argument")
	}

	aliasName := c.Args()[0]
	data := readFile()

	if data.TaskAlias[aliasName] == "" {
		log.Fatal("Alias does not exist")
	}

	delete(data.TaskAlias, aliasName)
	writeToFile(data)

	fmt.Println("Alias deleted:", aliasName)
}

func ShowAlias(c *cli.Context) {
	data := readFile()

	if len(data.TaskAlias) == 0 {
		log.Fatal("No alias exist at the moment...")
	}

	fmt.Println("Alias list:\n-----------")

	for key, val := range data.TaskAlias {
		fmt.Println(key + ": " + val)
	}
}

func TaskStatus(c *cli.Context) {
	data := readFile()

	if data.CurrentTask.Name == "" {
		log.Fatal("No task exist at the moment...")
	}

	t := formatTimeString(data.CurrentTask.Start)
	fmt.Println("Task status:\n------------\nNavn: " + data.CurrentTask.Name + "\nStart: " + t)
}
