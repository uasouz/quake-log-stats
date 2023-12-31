package main

import (
	"regexp"
	"strconv"
	"strings"
)

// getEventValue is a function that returns the value of an event based on its type
// and the remaining parts of the line
func getEventValue(eventType EventType, parts []string) (any, error) {
	switch eventType {
	case InitGame:
		return NewInitGameValue(parts)
	case Exit:
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

// Here you can find the implementation of the constructors of the event values
// The constructors are responsible for parsing the remaining parts of the line
// and returning the event value

type ClientUserinfoChangedValue struct {
	ClientID int
	UserInfo map[string]string
}

func NewClientUserinfoChangedValue(parts []string) (ClientUserinfoChangedValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientUserinfoChangedValue{}, err
	}

	userInfo := make(map[string]string)

	// the user info is a key-value pair separated by a backslash
	// so we need to split the remaining parts by a backslash and then iterate over the split parts by 2
	// to get the key-value pairs
	splitUserInfo := strings.Split(strings.Join(parts[1:], " "), "\\")
	for i := 0; i < len(splitUserInfo); i += 2 {
		userInfo[splitUserInfo[i]] = splitUserInfo[i+1]
	}

	return ClientUserinfoChangedValue{
		ClientID: clientID,
		UserInfo: userInfo,
	}, nil
}

type ClientDisconnectValue struct {
	ClientID int
}

func NewClientDisconnectValue(parts []string) (ClientDisconnectValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientDisconnectValue{}, err
	}

	return ClientDisconnectValue{
		ClientID: clientID,
	}, nil
}

type ClientBeginValue struct {
	ClientID int
}

func NewClientBeginValue(parts []string) (ClientBeginValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientBeginValue{}, err
	}

	return ClientBeginValue{
		ClientID: clientID,
	}, nil
}

type ShutdownGameValue struct {
}

func NewShutdownGameValue(parts []string) (ShutdownGameValue, error) {
	return ShutdownGameValue{}, nil
}

type ItemValue struct {
	ItemID   int
	ItemName string
}

func NewItemValue(parts []string) (ItemValue, error) {
	itemID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ItemValue{}, err
	}

	return ItemValue{
		ItemID:   itemID,
		ItemName: parts[1],
	}, nil
}

type ClientConnectValue struct {
	ClientID int
}

func NewClientConnectValue(parts []string) (ClientConnectValue, error) {
	clientID, err := strconv.Atoi(parts[0])
	if err != nil {
		return ClientConnectValue{}, err
	}

	return ClientConnectValue{
		ClientID: clientID,
	}, nil
}

type ExitGameValue struct {
	Reason string
}

func NewExitGameValue(parts []string) (ExitGameValue, error) {
	return ExitGameValue{
		Reason: strings.Join(parts, " "),
	}, nil
}

type InitGameValue struct {
	Settings map[string]string
}

func NewInitGameValue(parts []string) (InitGameValue, error) {
	// init game values are separated by a backslash.
	// each value is a key-value pair separated by a backslash
	configString := strings.TrimPrefix(strings.Join(parts, " "), "\\")

	settings := make(map[string]string)

	splitConfig := strings.Split(configString, "\\")

	// we must iterate over the split config by 2 because each key is followed by a value
	for i := 0; i < len(splitConfig); i += 2 {
		settings[splitConfig[i]] = splitConfig[i+1]
	}

	return InitGameValue{
		Settings: settings,
	}, nil
}

type KillValue struct {
	KillerID        int
	KillerName      string
	VictimID        int
	VictimName      string
	MeanOfDeathID   int
	MeanOfDeathName string
}

var killInfoRegex = regexp.MustCompile(`(.*) killed (.*) by (\w*)`)

func NewKillValue(parts []string) (KillValue, error) {
	// the first part of the kill value is the killer ID
	killerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return KillValue{}, err
	}

	// the second part of the kill value is the victim ID
	victimID, err := strconv.Atoi(parts[1])
	if err != nil {
		return KillValue{}, err
	}

	// the third part of the kill value is the mean of death ID
	meanOfDeathID, err := strconv.Atoi(strings.TrimSuffix(parts[2], ":"))
	if err != nil {
		return KillValue{}, err
	}

	// the rest of the parts are the kill info
	// in contrast to the other values, the kill info is not a single value
	// but a string that contains the killer name, the victim name and the mean of death name
	// and can suffer influence on how player inputs his name
	// so we need to parse it using a regex
	// We match the join parts of the kill info and get the capture groups to use as the kill info values
	killInfo := strings.Trim(strings.Join(parts[3:], " "), " ")

	killInfoParts := killInfoRegex.FindStringSubmatch(killInfo)

	return KillValue{
		KillerID:        killerID,
		KillerName:      killInfoParts[1],
		VictimID:        victimID,
		VictimName:      killInfoParts[2],
		MeanOfDeathID:   meanOfDeathID,
		MeanOfDeathName: killInfoParts[3],
	}, nil
}
