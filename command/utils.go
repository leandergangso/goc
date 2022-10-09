package command

import (
	"fmt"
	"log"
	"time"

	"github.com/LeanderGangso/goc/api"
	"google.golang.org/api/calendar/v3"
)

func insertToCalendar(calId string, newEvent *calendar.Event) {
	client := api.GetClient()
	event, err := client.Events.Insert(calId, newEvent).Do()
	if err != nil {
		log.Fatalf("Unable to add event to calendar: %v", err)
	}
	fmt.Println("Old task added to calendar:", event.HtmlLink)
}

func createEvent(data *FileData) *calendar.Event {
	return &calendar.Event{
		Summary: data.CurrentTask.Name,
		Start: &calendar.EventDateTime{
			DateTime: data.CurrentTask.Start,
		},
		End: &calendar.EventDateTime{
			DateTime: getTime(),
		},
	}
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}
