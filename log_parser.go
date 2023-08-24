package main

import (
	"bufio"
	"fmt"
	"io"
)

func ParseLog(logfile io.Reader) ([]Event, error) {
	scanner := bufio.NewScanner(logfile)
	events := make([]Event, 0)
	for scanner.Scan() {
		line := scanner.Text()
		event, err := ParseLine(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		events = append(events, event)
	}

	return events, nil
}
