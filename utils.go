package goc

import (
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

// magic reference: Mon Jan 2 15:04:05 MST 2006
const TIME_FORMAT = time.RFC3339

func insertToCalendar(calId string, newEvent *calendar.Event) (*calendar.Event, error) {
	client, source, err := GetClient()
	if err != nil {
		return nil, err
	}

	event, err := client.Events.Insert(calId, newEvent).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to add event to calendar: %v", err)
	}

	updateToken(source)
	return event, nil
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
	return time.Now().Format(TIME_FORMAT)
}

func getTimeSince(start string) (string, error) {
	startTime, err := time.Parse(TIME_FORMAT, start)
	if err != nil {
		return "", err
	}
	duration := time.Since(startTime).Round(time.Second)
	return duration.String(), nil
}

func stringToTime(s string) string {
	now := time.Now()
	timezone, _ := now.Zone()
	fs := fmt.Sprintf("%d-%d-%d %v %v", now.Year(), now.Month(), now.Day(), s, timezone)
	t, err := time.Parse("2006-1-2 15:04 MST", fs)
	if err != nil {
		log.Fatalf("unable to parse time: %v", err)
	}
	return t.Format(TIME_FORMAT)
}

func formatTimeString(s string) string {
	data := strings.Split(s, "T")
	return fmt.Sprintf("%v %v", data[0], strings.Split(data[1], "+")[0][:5])
}

func checkAndUseAlias(name string, data *FileData) string {
	aliasName := data.TaskAlias[name]
	if aliasName != "" {
		return aliasName
	}
	return name
}

func updatePrevTaskAlias(data *FileData) {

	if _, ok := data.TaskAlias["prev2"]; ok {
		data.TaskAlias["prev3"] = data.TaskAlias["prev2"]
	}

	if _, ok := data.TaskAlias["prev"]; ok {
		data.TaskAlias["prev2"] = data.TaskAlias["prev"]
	}

	if data.TaskAlias == nil {
		data.TaskAlias = make(map[string]string)
	}

	data.TaskAlias["prev"] = data.CurrentTask.Name
}
