package goc

import (
	"encoding/json"
	"log"
	"os"
)

const CONFIG_PATH = ""
const SHARED_PATH = CONFIG_PATH + "/goc_cli/"
const FILE_NAME = "data.json"

func readFile() *FileData {
	path := SHARED_PATH + FILE_NAME

	createDataFileIfNotExists(path)

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer f.Close()
	data := &FileData{}
	err = json.NewDecoder(f).Decode(data)
	if err != nil {
		log.Fatalf("Unable to read from file: %v", err)
	}
	return data
}

func writeToFile(data *FileData) {
	path := SHARED_PATH + FILE_NAME
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Unable to open or create file: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		log.Fatalf("Unable to write to file: %v", err)
	}
}

func createDataFileIfNotExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		data := []byte("{}")
		err := os.WriteFile(path, data, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
