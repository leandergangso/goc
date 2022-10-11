package goc

import (
	"fmt"
	"log"
	"time"

	"google.golang.org/api/calendar/v3"
)

// magic reference: Mon Jan 2 15:04:05 MST 2006
var TIME_FORMAT = "2006-01-02 15:04 MST"

func insertToCalendar(calId string, newEvent *calendar.Event) *calendar.Event {
	client := GetClient()
	event, err := client.Events.Insert(calId, newEvent).Do()
	if err != nil {
		log.Fatalf("Unable to add event to calendar: %v", err)
	}
	return event
}

func createEvent(data *FileData, endTime string) *calendar.Event {
	return &calendar.Event{
		Summary: data.CurrentTask.Name,
		Start: &calendar.EventDateTime{
			DateTime: data.CurrentTask.Start,
			TimeZone: "Europe/Oslo",
		},
		End: &calendar.EventDateTime{
			DateTime: endTime,
			TimeZone: "Europe/Oslo",
		},
	}
}

func getTime() string {
	return time.Now().Format(time.RFC3339)
}

func stringToTime(s string) string {
	now := time.Now()
	fs := fmt.Sprintf("%d-%d-%d %v %v", now.Year(), now.Month(), now.Day(), s, "CEST")
	t, err := time.Parse(TIME_FORMAT, fs)
	if err != nil {
		log.Fatalf("Unable to parse time: %v", err)
	}
	return t.Format(time.RFC3339)
}
