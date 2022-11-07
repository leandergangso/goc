package goc

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

func GoogleSetup(c *cli.Context) error {
	client, _ := GetClient()
	calList, err := client.CalendarList.List().Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve calendar list: %v", err)
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
		fmt.Printf("Skipped, currently using: %v\n", data.CalendarId)
		return nil
	}

	data.CalendarId = calId
	writeToFile(data)

	fmt.Println("Calendar ID added, you are ready to start tracking!")
	return nil
}

func StartTask(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data := readFile()
	name := checkAndUseAlias(c.Args().Get(0), data)
	desc := c.Args().Get(1)

	startTime := c.String("time")
	if startTime == "" {
		startTime = getTime()
	} else {
		startTime = stringToTime(startTime)
	}

	if data.CurrentTask.Name != "" {
		newEvent := createEvent(data, startTime)
		event := insertToCalendar(data, newEvent)

		updatePrevTaskAlias(data)
		fmt.Println("Task added to calendar:", event.HtmlLink)
	}

	taskName := name
	if desc != "" {
		taskName += " " + desc
	}

	data.CurrentTask.Name = taskName
	data.CurrentTask.Start = startTime
	writeToFile(data)

	fmt.Println("New task started: " + name)
	return nil
}

func EndTask(c *cli.Context) error {
	data := readFile()
	if data.CurrentTask.Name == "" {
		fmt.Println("No current task to end")
		return nil
	}

	endTime := c.String("time")
	if endTime == "" {
		endTime = getTime()
	} else {
		endTime = stringToTime(endTime)
	}

	newEvent := createEvent(data, endTime)
	event := insertToCalendar(data, newEvent)

	updatePrevTaskAlias(data)
	data.CurrentTask.Reset()
	writeToFile(data)

	fmt.Println("Task added to calendar:", event.HtmlLink)
	return nil
}

func EditCurrentTask(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data := readFile()
	name := checkAndUseAlias(c.String("name"), data)
	start := c.String("time")

	if name != "" {
		data.CurrentTask.Name = name
		fmt.Println("Task name set: " + name)
	}
	if start != "" {
		start = stringToTime(start)
		data.CurrentTask.Start = start
		fmt.Println("Start time set: " + formatTimeString(start))
	}

	writeToFile(data)
	return nil
}

func InsertTask(c *cli.Context) error {
	if c.NArg() < 3 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data := readFile()
	name := checkAndUseAlias(c.Args().Get(0), data)
	startTime := stringToTime(c.Args().Get(1))
	endTime := stringToTime(c.Args().Get(2))

	data.CurrentTask.Name = name
	data.CurrentTask.Start = startTime

	newEvent := createEvent(data, endTime)
	event := insertToCalendar(data, newEvent)

	data.CurrentTask.Reset()
	writeToFile(data)
	fmt.Println("Task added to calendar:", event.HtmlLink)
	return nil
}

func AddTaskAlias(c *cli.Context) error {
	if c.NArg() < 2 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	aliasName := c.Args().Get(0)
	taskName := c.Args().Get(1)

	data := readFile()
	if data.TaskAlias == nil {
		data.TaskAlias = make(map[string]string)
	}

	data.TaskAlias[aliasName] = taskName
	writeToFile(data)
	fmt.Println("Alias added:", aliasName+": "+taskName)
	return nil
}

func DelTaskAlias(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	aliasName := c.Args().Get(0)
	data := readFile()

	if data.TaskAlias[aliasName] == "" {
		fmt.Println("Alias does not exist")
		os.Exit(0)
	}

	delete(data.TaskAlias, aliasName)
	writeToFile(data)
	fmt.Println("Alias deleted:", aliasName)
	return nil
}

func ShowAlias(c *cli.Context) error {
	data := readFile()
	if len(data.TaskAlias) == 0 {
		fmt.Println("No alias exists")
		return nil
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
	fmt.Println("prev" + ": " + data.TaskAlias["prev"])
	fmt.Println("prev2" + ": " + data.TaskAlias["prev2"])
	fmt.Println("prev3" + ": " + data.TaskAlias["prev3"])
	return nil
}

func TaskStatus(c *cli.Context) error {
	data := readFile()

	if c.Bool("toggle") {
		data.StatusOneline = !data.StatusOneline
		writeToFile(data)
		fmt.Printf("Oneline set to: %v\n", data.StatusOneline)
		return nil
	}

	if c.Bool("list") {
		client, _ := GetClient()
		eventList := getTodaysCalendarEvents(client, data)

		if len(eventList.Items) == 0 {
			fmt.Println("No task for today")
			return nil
		}

		customFormat := "15:04"

		fmt.Println("Todays tasks:")
		for _, evt := range eventList.Items {
			start, _ := time.Parse(TIME_FORMAT, evt.Start.DateTime)
			end, _ := time.Parse(TIME_FORMAT, evt.End.DateTime)
			duration := end.Sub(start)
			fmt.Printf("- %v: %v %v (%v)\n", evt.Summary, start.Format(customFormat), end.Format(customFormat), duration)
		}
		return nil
	}

	durationToday := data.GetDurationToday(c.Bool("update"))

	if data.CurrentTask.Name == "" {
		if data.StatusOneline || c.Bool("oneline") {
			fmt.Printf("No current task (%v)\n", durationToday)
		} else {
			fmt.Println("No current task")
			fmt.Println("Duration today:", durationToday)
		}
		return nil
	}

	taskDuration := getTimeSince(data.CurrentTask.Start)

	totalDuration := taskDuration + durationToday

	if data.StatusOneline || c.Bool("oneline") {
		fmt.Printf("%s (%v) (%v)\n", data.CurrentTask.Name, taskDuration, totalDuration)
		return nil
	}

	startTime := formatTimeString(data.CurrentTask.Start)

	fmt.Println("Task status:\n------------")
	fmt.Println("Name:", data.CurrentTask.Name)
	fmt.Println("Start:", startTime)
	fmt.Println("Duration:", taskDuration)
	fmt.Println("Duration today:", totalDuration)
	return nil
}
