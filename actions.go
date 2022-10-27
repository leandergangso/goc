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

	data := readFile()

	if calId == "" {
		fmt.Printf("Skipped, currently using: %v", data.CalendarId)
		os.Exit(0)
	}

	data.CalendarId = calId
	writeToFile(data)

	fmt.Println("Calendar ID added, you are ready to start tracking!")
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
		newEvent := createEvent(data, startTime)
		event := insertToCalendar(data.CalendarId, newEvent)
		updatePrevTaskAlias(data)
		fmt.Println("Task added to calendar:", event.HtmlLink)
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
	updatePrevTaskAlias(data)
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

	aliasName := c.Args()[0]
	taskName := c.Args()[1]

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
		fmt.Println("Alias does not exist")
		os.Exit(0)
	}

	delete(data.TaskAlias, aliasName)
	writeToFile(data)

	fmt.Println("Alias deleted:", aliasName)
}

func ShowAlias(c *cli.Context) {
	data := readFile()

	if len(data.TaskAlias) == 0 {
		fmt.Println("No alias exist at the moment...")
		os.Exit(0)
	}

	fmt.Println("Alias list:\n-----------")

	prevTasks := make(map[string]string)

	for key, val := range data.TaskAlias {
		if strings.Contains(key, "prev") {
			prevTasks[key] = val
		} else {
			fmt.Println(key + ": " + val)
		}
	}
	fmt.Println()
	for key, val := range prevTasks {
		fmt.Println(key + ": " + val)
	}
}

func TaskStatus(c *cli.Context) {
	data := readFile()

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		os.Exit(0)
	}

	t := formatTimeString(data.CurrentTask.Start)
	fmt.Println("Task status:\n------------\nNavn: " + data.CurrentTask.Name + "\nStart: " + t)
}
