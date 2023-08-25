package main

import (
	"fmt"
	"strconv"
	"strings"
)

// parseTimeStamp parses a string in the format "MM:SS" and returns the
// corresponding number of seconds
func parseTimeStamp(value string) (int, error) {
	times := strings.Split(value, ":")

	minutes, err := strconv.Atoi(times[0])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(times[1])
	if err != nil {
		return 0, err
	}

	return minutes*60 + seconds, nil
}

// ParseLine parses a line of the log file and returns an Event
func ParseLine(line string) (Event, error) {
	// recover from panic in case of bad input
	// you can check a bad input example in the line 97 of the file "games.log"
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()

	// in order to parse the line we split it into parts
	parts := strings.Split(strings.Trim(line, " "), " ")

	// the first part is the timestamp
	eventTimeStamp, err := parseTimeStamp(parts[0])
	if err != nil {
		return Event{}, err
	}

	// the second part is the event type
	eventType := getEventType(parts[1])

	// the rest of the parts are the event value, which depends on the event type
	eventValue, err := getEventValue(eventType, parts[2:])

	if err != nil {
		return Event{}, err
	}

	return Event{
		timestamp: eventTimeStamp,
		tokenType: eventType,
		value:     eventValue,
	}, nil
}
