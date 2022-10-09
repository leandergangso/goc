package command

import (
	"encoding/json"
	"log"
	"os"
)

const fileName = "goc_data.json"

func readFile() *FileData {
	f, err := os.Open(fileName)
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
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to open file: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(data)
	if err != nil {
		log.Fatalf("Unable to write to file: %v", err)
	}
}
