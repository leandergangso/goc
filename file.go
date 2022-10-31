package goc

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const SHARED_PATH = "/goc_cli"
const FILE_NAME = "/data.json"

type FileData struct {
	CalendarId  string
	CurrentTask DataTask
	TaskAlias   map[string]string
}

type DataTask struct {
	Name  string
	Start string
}

func (f *DataTask) Reset() {
	f.Name = ""
	f.Start = ""
}

func getSharedPath() (string, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	sharedPath := configPath + SHARED_PATH
	return sharedPath, nil
}

func getFilePath() (string, error) {
	commonPath, err := getSharedPath()
	if err != nil {
		return "", err
	}
	fullFilePath := commonPath + FILE_NAME
	return fullFilePath, nil
}

func readFile() (*FileData, error) {
	filepath, err := getFilePath()
	if err != nil {
		return nil, fmt.Errorf("unable to get path: %v", err)
	}

	f, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("unable to read/create file: %v", err)
	}
	defer f.Close()

	data := &FileData{}
	err = json.NewDecoder(f).Decode(data)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("unable to decode data from file: %v", err)
	}
	return data, nil
}

func writeToFile(data *FileData) error {
	filepath, err := getFilePath()
	if err != nil {
		return fmt.Errorf("unable to get path: %v", err)
	}

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to write/create file: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return fmt.Errorf("unable to encode data to file: %v", err)
	}
	return nil
}
