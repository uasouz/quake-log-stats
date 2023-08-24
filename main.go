package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func main() {
	parser := NewQuakeLogParser()
	events, err := parser.Parse("games.log")
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
		case ClientUserinfoChanged:
			eventValue := event.value.(ClientUserinfoChangedValue)
			match.AddPlayer(eventValue.clientID, eventValue.userInfo["n"])
		case Kill:
			eventValue := event.value.(KillValue)
			match.AddKill(eventValue.killerID, eventValue.victimID)
		case ShutdownGame:
			reports[fmt.Sprintf("game_%d", gameID)] = match.Report()
		}
	}

	return reports
}
