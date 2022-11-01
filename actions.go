package goc

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func GoogleSetup(c *cli.Context) error {
	client, _, err := GetClient()
	if err != nil {
		return err
	}

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

	data, err := readFile()
	if err != nil {
		return err
	}

	if calId == "" {
		fmt.Printf("Skipped, currently using: %v\n", data.CalendarId)
		return nil
	}

	data.CalendarId = calId
	err = writeToFile(data)
	if err != nil {
		return err
	}

	fmt.Println("Calendar ID added, you are ready to start tracking!")
	return nil
}

func StartTask(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data, err := readFile()
	if err != nil {
		return err
	}

	name := checkAndUseAlias(c.Args().Get(0), data)

	startTime := c.String("time")
	if startTime == "" {
		startTime = getTime()
	} else {
		startTime = stringToTime(startTime)
	}

	if data.CurrentTask.Name != "" {
		newEvent := createEvent(data, startTime)
		event, err := insertToCalendar(data.CalendarId, newEvent)
		if err != nil {
			return err
		}

		updatePrevTaskAlias(data)
		fmt.Println("Task added to calendar:", event.HtmlLink)
	}

	data.CurrentTask.Name = name
	data.CurrentTask.Start = startTime
	err = writeToFile(data)
	if err != nil {
		return err
	}

	fmt.Println("New task started: " + name)
	return nil
}

func EndTask(c *cli.Context) error {
	data, err := readFile()
	if err != nil {
		return err
	}

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		return nil
	}

	endTime := c.String("time")
	if endTime == "" {
		endTime = getTime()
	} else {
		endTime = stringToTime(endTime)
	}

	newEvent := createEvent(data, endTime)
	event, err := insertToCalendar(data.CalendarId, newEvent)
	if err != nil {
		return err
	}

	updatePrevTaskAlias(data)
	data.CurrentTask.Reset()

	err = writeToFile(data)
	if err != nil {
		return err
	}

	fmt.Println("Task added to calendar:", event.HtmlLink)
	return nil
}

func EditCurrentTask(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data, err := readFile()
	if err != nil {
		return err
	}

	name := checkAndUseAlias(c.String("name"), data)
	start := c.String("time")

	if name != "" {
		data.CurrentTask.Name = name
		fmt.Println("New task name set: " + name)
	}

	if start != "" {
		start = stringToTime(start)
		data.CurrentTask.Start = start
		fmt.Println("New start time set: " + formatTimeString(start))
	}

	err = writeToFile(data)
	if err != nil {
		return err
	}

	return nil
}

func InsertTask(c *cli.Context) error {
	if c.NArg() < 3 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	data, err := readFile()
	if err != nil {
		return err
	}

	name := checkAndUseAlias(c.Args().Get(0), data)
	startTime := stringToTime(c.Args().Get(1))
	endTime := stringToTime(c.Args().Get(2))

	data.CurrentTask.Name = name
	data.CurrentTask.Start = startTime

	newEvent := createEvent(data, endTime)
	event, err := insertToCalendar(data.CalendarId, newEvent)
	if err != nil {
		return err
	}

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

	data, err := readFile()
	if err != nil {
		return err
	}

	if data.TaskAlias == nil {
		data.TaskAlias = make(map[string]string)
	}

	data.TaskAlias[aliasName] = taskName
	err = writeToFile(data)
	if err != nil {
		return err
	}

	fmt.Println("Alias added:", aliasName+": "+taskName)
	return nil
}

func DelTaskAlias(c *cli.Context) error {
	if c.NArg() < 1 {
		cli.ShowSubcommandHelp(c)
		return nil
	}

	aliasName := c.Args().Get(0)
	data, err := readFile()
	if err != nil {
		return err
	}

	if data.TaskAlias[aliasName] == "" {
		fmt.Println("Alias does not exist")
		os.Exit(0)
	}

	delete(data.TaskAlias, aliasName)
	err = writeToFile(data)
	if err != nil {
		return err
	}

	fmt.Println("Alias deleted:", aliasName)
	return nil
}

func ShowAlias(c *cli.Context) error {
	data, err := readFile()
	if err != nil {
		return err
	}

	if len(data.TaskAlias) == 0 {
		fmt.Println("No alias exist at the moment...")
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
	data, err := readFile()
	if err != nil {
		return err
	}

	if data.CurrentTask.Name == "" {
		fmt.Println("No task exist at the moment...")
		return nil
	}

	duration, err := getTimeSince(data.CurrentTask.Start)
	if err != nil {
		return fmt.Errorf("unable to get time durtaion: %v", err)
	}

	if c.Bool("oneline") {
		fmt.Printf("%s (%v)\n", data.CurrentTask.Name, duration)
		return nil
	}

	startTime := formatTimeString(data.CurrentTask.Start)

	fmt.Println("Task status:\n------------")
	fmt.Println("Name:", data.CurrentTask.Name)
	fmt.Println("Start:", startTime)
	fmt.Println("Duration:", duration)
	return nil
}
