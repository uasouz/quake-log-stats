package main

import (
	"fmt"
	"regexp"
	"strings"
)

type EventType int

const (
	TimeStamp EventType = iota + 1
	InitGame
	Exit
	ShutdownGame
	ClientConnect
	ClientUserinfoChanged
	ClientDisconnect
	ClientBegin
	Item
	Kill
	score
	divider
)

func (t EventType) String() string {
	return [...]string{
		"TimeStamp",
		"InitGame",
		"Exit",
		"ShutdownGame",
		"ClientConnect",
		"ClientUserinfoChanged",
		"ClientDisconnect",
		"ClientBegin",
		"Item",
		"Kill",
		"score",
		"divider",
	}[t-1]
}

type Event struct {
	timestamp int
	tokenType EventType
	value     any
}

func (e Event) String() string {
	readableTimestamp := fmt.Sprintf("%02d:%02d", e.timestamp/60, e.timestamp%60)
	return fmt.Sprintf("Event{timestamp: %d | %s , tokenType: %s, value: %v}",
		e.timestamp, readableTimestamp, e.tokenType, e.value)
}

func getEventType(value string) EventType {
	switch strings.TrimSuffix(value, ":") {
	case "InitGame":
		return InitGame
	case "Exit":
		return Exit
	case "ShutdownGame":
		return ShutdownGame
	case "ClientConnect":
		return ClientConnect
	case "ClientUserinfoChanged":
		return ClientUserinfoChanged
	case "ClientDisconnect":
		return ClientDisconnect
	case "ClientBegin":
		return ClientBegin
	case "Item":
		return Item
	case "Kill":
		return Kill
	case "score":
		return score
	}
	if match, _ := regexp.MatchString("^[0-9]{2}:[0-9]{2}$", value); match {
		return TimeStamp
	}
	return divider
}
