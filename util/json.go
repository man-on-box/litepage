package util

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func ParseJSONFile(file string, data any) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Could not open file: %v", err)
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		log.Fatalf("Could not read file: %v", err)
	}

	err = json.Unmarshal(contents, &data)
	if err != nil {
		log.Fatalf("Could not parse JSON: %v", err)
	}
}
