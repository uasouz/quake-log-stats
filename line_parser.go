package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type EventType int

const (
	TimeStamp EventType = iota + 1
	InitGame
	ExitGame
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
		"ExitGame",
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
	case "ExitGame":
		return ExitGame
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

func ParseLine(line string) (Event, error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	parts := strings.Split(strings.Trim(line, " "), " ")

	eventTimeStamp, err := parseTimeStamp(parts[0])
	if err != nil {
		return Event{}, err
	}

	eventType := getEventType(parts[1])

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

func getEventValue(eventType EventType, parts []string) (any, error) {
	switch eventType {
	case InitGame:
		return NewInitGameValue(parts)
	case ExitGame:
		return NewExitGameValue(parts)
	case ShutdownGame:
		return NewShutdownGameValue(parts)
	case ClientConnect:
		return NewClientConnectValue(parts)
	case ClientUserinfoChanged:
		return NewClientUserinfoChangedValue(parts)
	case ClientDisconnect:
		return NewClientDisconnectValue(parts)
	case ClientBegin:
		return NewClientBeginValue(parts)
	case Item:
		return NewItemValue(parts)
	case Kill:
		return NewKillValue(parts)
	}
	return nil, nil
}

type ClientUserinfoChangedValue struct {
	clientID int
	userInfo map[string]string
}

func NewClientUserinfoChangedValue(parts []string) (ClientUserinfoChangedValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientUserinfoChangedValue{}, err
	}

	userInfo := make(map[string]string)

	splitUserInfo := strings.Split(strings.Join(parts[1:], " "), "\\")
	for i := 0; i < len(splitUserInfo); i += 2 {
		userInfo[splitUserInfo[i]] = splitUserInfo[i+1]
	}

	return ClientUserinfoChangedValue{
		clientID: clientID,
		userInfo: userInfo,
	}, nil
}

type ClientDisconnectValue struct {
	clientID int
}

func NewClientDisconnectValue(parts []string) (ClientDisconnectValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientDisconnectValue{}, err
	}

	return ClientDisconnectValue{
		clientID: clientID,
	}, nil
}

type ClientBeginValue struct {
	clientID int
}

func NewClientBeginValue(parts []string) (ClientBeginValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientBeginValue{}, err
	}

	return ClientBeginValue{
		clientID: clientID,
	}, nil
}

type ShutdownGameValue struct {
}

func NewShutdownGameValue(parts []string) (ShutdownGameValue, error) {
	return ShutdownGameValue{}, nil
}

type ItemValue struct {
	itemID   int
	itemName string
}

func NewItemValue(parts []string) (ItemValue, error) {
	itemID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ItemValue{}, err
	}

	return ItemValue{
		itemID:   itemID,
		itemName: parts[1],
	}, nil
}

type ClientConnectValue struct {
	clientID int
}

func NewClientConnectValue(parts []string) (ClientConnectValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientConnectValue{}, err
	}

	return ClientConnectValue{
		clientID: clientID,
	}, nil
}

type ExitGameValue struct {
	reason string
}

func NewExitGameValue(parts []string) (ExitGameValue, error) {
	return ExitGameValue{
		reason: strings.Join(parts, " "),
	}, nil
}

type InitGameValue struct {
	settings map[string]string
}

func NewInitGameValue(parts []string) (InitGameValue, error) {
	configString := strings.TrimPrefix(strings.Join(parts, " "), "\\")

	settings := make(map[string]string)

	splitConfig := strings.Split(configString, "\\")

	for i := 0; i < len(splitConfig); i += 2 {
		settings[splitConfig[i]] = splitConfig[i+1]
	}

	return InitGameValue{
		settings: settings,
	}, nil
}

type KillValue struct {
	killerID       int
	killerName     string
	victimID       int
	victimName     string
	deathCauseID   int
	deathCauseName string
}

func NewKillValue(parts []string) (KillValue, error) {
	killerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return KillValue{}, err
	}

	victimID, err := strconv.Atoi(parts[1])
	if err != nil {
		return KillValue{}, err
	}

	deathCauseID, err := strconv.Atoi(strings.TrimSuffix(parts[2], ":"))
	if err != nil {
		return KillValue{}, err
	}

	return KillValue{
		killerID:       killerID,
		killerName:     parts[3],
		victimID:       victimID,
		victimName:     parts[4],
		deathCauseID:   deathCauseID,
		deathCauseName: parts[5],
	}, nil
}
