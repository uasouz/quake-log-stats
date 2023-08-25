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

	matchesData, err := readFileAndGenerateData(*filePath)

	if err != nil {
		panic(err)
	}

	// save matchesData to JSON file
	outputFile, err := os.OpenFile(fmt.Sprintf("matches_data_%d.json", time.Now().Unix()), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer outputFile.Close()

	err = saveMatchesDataToJSONFile(matchesData, outputFile)

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Matches data saved to file %s", outputFile.Name()))

	reportFileName, err := renderReport(matchesData)

	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("Report saved to file %s", reportFileName))

	fmt.Println("Done!")
}

func readFileAndGenerateData(filePath string) (map[string]MatchData, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	events, err := ParseLog(file)
	if err != nil {
		return nil, err
	}

	return readAndAggregateEvents(events), nil
}

func saveMatchesDataToJSONFile(data map[string]MatchData, outputFile *os.File) error {
	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func readAndAggregateEvents(events []Event) map[string]MatchData {
	gameID := 0
	matches := make([]*Match, 0)

	matchesData := make(map[string]MatchData)

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
			matchesData[fmt.Sprintf("game_%d", gameID)] = match.FormatData()
		}
	}

	return matchesData
}
