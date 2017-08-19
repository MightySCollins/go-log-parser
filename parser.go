package main

import (
	"github.com/hpcloud/tail"
	"regexp"
	"log"
	"encoding/csv"
	"os"
)

var regex = regexp.MustCompile(`publicvariable\.log:(?P<date>\d\d\.\d\d\.\d\d\d\d \d\d:\d\d:\d\d): (?P<name>.*?) \( \) (?P<guid>.*?) - #0 "(?P<data>.*?)" = any`)
var headers = []string{"full", "date", "name", "guid", "data"}

func main() {
	t, err := tail.TailFile("/home/scollins/GoglandProjects/log-parser/publicvariable.log", tail.Config{
		Follow: true,
		ReOpen: true,
	})
	checkError("Failed to read log file", err)

	file, err := os.Create("/home/scollins/GoglandProjects/log-parser/result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(headers)
	checkError("Cannot write headers to file", err)

	for line := range t.Lines {
		matches  := regex.FindStringSubmatch(line.Text)
		err := writer.Write(matches)
		checkError("Cannot write line to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}