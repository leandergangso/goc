package goc

import (
	"encoding/json"
	"log"
	"os"
)

var SHARED_PATH = os.Getenv("HOME") + "/.goc_cli/"
var FILE_NAME = "data.json"

func readFile() *FileData {
	path := SHARED_PATH + FILE_NAME
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
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		log.Fatalf("Unable to write to file: %v", err)
	}
}
