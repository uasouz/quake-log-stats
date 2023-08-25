package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {

	filePath := flag.String("file", "", "path to the log file")

	flag.Parse()

	if filePath == nil {
		panic("you must provide a file path")
	}

	if filePath != nil && *filePath == "" {
		panic("you must provide a file path")
	}

	file, err := os.OpenFile(*filePath, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	events, err := ParseLog(file)
	if err != nil {
		panic(err)
	}

	reports := readAndAggregateEvents(events)

	// save reports to JSON file
	outputFile, err := os.OpenFile(fmt.Sprintf("reports_%d.json", time.Now().Unix()), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = saveReportsToJSONFile(reports, outputFile)

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Reports saved to file %s", outputFile.Name()))
	fmt.Println("Done!")
}

func saveReportsToJSONFile(reports map[string]MatchReport, outputFile *os.File) error {
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(reports)
}

func readAndAggregateEvents(events []Event) map[string]MatchReport {
	gameID := 0
	matches := make([]*Match, 0)

	reports := make(map[string]MatchReport)

	var match *Match
	for _, event := range events {
		switch event.tokenType {
		case InitGame:
			gameID++
			match = NewMatch(gameID)
			matches = append(matches, match)
		case ClientUserinfoChanged, Kill:
			match.AggregateEvent(event)
		case ShutdownGame:
			reports[fmt.Sprintf("game_%d", gameID)] = match.Report()
		}
	}

	return reports
}
