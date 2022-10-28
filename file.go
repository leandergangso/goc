package goc

import (
	"encoding/json"
	"fmt"
	"os"
)

// func setup() {
// 	if err != nil {
// 		log.Fatalf("Unable to resolve config dir: %v", err)
// 	}
// 	fmt.Println("Using path:", CONFIG_PATH)
// }

const CONFIG_PATH = ""
const SHARED_PATH = CONFIG_PATH + "/goc_cli/"
const FILE_NAME = "data.json"

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

func readFile() (*FileData, error) {
	path := SHARED_PATH + FILE_NAME

	err := createDataFileIfNotExists(path)
	if err != nil {
		return nil, fmt.Errorf("unable to create file: %v", err)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %v", err)
	}
	defer f.Close()

	data := &FileData{}
	err = json.NewDecoder(f).Decode(data)
	if err != nil {
		return nil, fmt.Errorf("unable to decode data from file: %v", err)
	}
	return data, nil
}

func writeToFile(data *FileData) error {
	path := SHARED_PATH + FILE_NAME
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to open/create file: %v", err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		return fmt.Errorf("unable to encode data to file: %v", err)
	}
	return nil
}

func createDataFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		data := []byte("{}")
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			return fmt.Errorf("unable to write to file: %v", err)
		}
	}
	return nil
}
