package goc

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"
)

func GoogleAuth(c *cli.Context) {
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
