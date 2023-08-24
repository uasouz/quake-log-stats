package main

import (
	"bufio"
	"fmt"
	"os"
)

type QuakeLogParser struct {
}

func (p *QuakeLogParser) Parse(filePath string) ([]Event, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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

func NewQuakeLogParser() *QuakeLogParser {
	return &QuakeLogParser{}
}
