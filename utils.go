package goc

import (
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"google.golang.org/api/calendar/v3"
)

const TIME_FORMAT = time.RFC3339

func insertToCalendar(data *FileData, newEvent *calendar.Event) *calendar.Event {
	client, source := GetClient()

	event, err := client.Events.Insert(data.CalendarId, newEvent).Do()
	if err != nil {
		log.Fatalf("unable to add event to calendar: %v", err)
	}

	updateTotalDuration(client, data)
	updateToken(source)
	return event
}

func updateTotalDuration(client *calendar.Service, data *FileData) {
	listCall := client.Events.List(data.CalendarId)

	year, month, day := time.Now().Date()
	minTime := time.Date(year, month, day, 0, 0, 0, 0, time.Now().Location())

	listCall.TimeMin(minTime.Format(TIME_FORMAT))
	listCall.TimeMax(getTime())

	eventList, err := listCall.Do()
	if err != nil {
		log.Fatalf("unable to get calendar events: %v", err)
	}

	totalDuration := 0.0

	for _, evt := range eventList.Items {
		start, _ := time.Parse(TIME_FORMAT, evt.Start.DateTime)
		end, _ := time.Parse(TIME_FORMAT, evt.End.DateTime)
		totalDuration += end.Sub(start).Seconds()
	}

	data.DurationToday = time.Duration(math.Round(totalDuration)) * time.Second
	data.CurrentDate = CurDate{year, month, day}
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

func getTimeSince(start string) time.Duration {
	startTime, err := time.Parse(TIME_FORMAT, start)
	if err != nil {
		log.Fatalf("unable to parse time: %v", err)
	}
	duration := time.Since(startTime).Round(time.Second)
	return duration
}

func stringToTime(s string) string {
	now := time.Now()
	timezone, _ := now.Zone()
	fs := fmt.Sprintf("%d-%d-%d %v %v", now.Year(), now.Month(), now.Day(), s, timezone)
	t, err := time.Parse("2006-1-2 1504 MST", fs)
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
