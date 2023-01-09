package goc

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

const SHARED_PATH = "/goc_cli"
const FILE_NAME = "/data.json"

type FileData struct {
	CalendarId    string
	CurrentTask   DataTask
	TaskAlias     map[string]string
	DurationToday time.Duration
	CurrentDate   CurDate
	StatusOneline bool
	UpdateToken   TokenStatus
	Jira          *JiraAuth
}

func (f *FileData) GetDurationToday(force bool) time.Duration {
	if !force {
		year, month, day := time.Now().Date()
		date := f.CurrentDate
		if date.Year == year && date.Month == month && date.Day == day {
			return f.DurationToday // updated on new events to cal
		}
	}

	client, source := GetClient()
	updateTotalDuration(client, f)
	updateToken(source)
	writeToFile(f)

	return f.DurationToday
}

type JiraAuth struct {
	Username string
	Token    string
}

type TokenStatus struct {
	DayNumber int
	Done      bool
}

type CurDate struct {
	Year  int
	Month time.Month
	Day   int
}

type DataTask struct {
	Name  string
	Start string
}

func (f *DataTask) Reset() {
	f.Name = ""
	f.Start = ""
}

func getSharedPath() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("unable to get userConfigDir: %v", err)
	}

	sharedPath := configPath + SHARED_PATH

	return sharedPath
}

func getFilePath() string {
	commonPath := getSharedPath()
	fullFilePath := commonPath + FILE_NAME

	return fullFilePath
}

func readFile() *FileData {
	filepath := getFilePath()

	f, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("unable to read/create file: %v", err)
	}
	defer f.Close()

	data := &FileData{}
	err = json.NewDecoder(f).Decode(data)
	if err != nil && err != io.EOF {
		log.Fatalf("unable to decode data from file: %v", err)
	}

	return data
}

func writeToFile(data *FileData) {
	filepath := getFilePath()

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("unable to write/create file: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		log.Fatalf("unable to encode data to file: %v", err)
	}
}
