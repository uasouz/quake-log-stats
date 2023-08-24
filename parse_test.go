package main

import (
	"encoding/json"
	"testing"
)

func TestParseLine(t *testing.T) {
	linesToTest := []string{
		"1:47 InitGame: \\sv_floodProtect\\1\\sv_maxPing\\0\\sv_minPing\\0\\sv_maxRate\\10000\\sv_minRate\\0\\sv_hostname\\Code Miner Server\\g_gametype\\0\\sv_privateClients\\2\\sv_maxclients\\16\\sv_allowDownload\\0\\bot_minplayers\\0\\dmflags\\0\\fraglimit\\20\\timelimit\\15\\g_maxGameClients\\0\\capturelimit\\8\\version\\ioq3 1.36 linux-x86_64 Apr 12 2009\\protocol\\68\\mapname\\q3dm17\\gamename\\baseq3\\g_needpass\\0\n",
		"1:47 ClientConnect: 2",
		"3:49 ClientBegin: 5",
		"1:47 ClientUserinfoChanged: 2 n\\Dono da Bola\\t\\0\\model\\sarge\\hmodel\\sarge\\g_redteam\\\\g_blueteam\\\\c1\\4\\c2\\5\\hc\\95\\w\\0\\l\\0\\tt\\0\\tl\\0",
		"1:48 Item: 4 ammo_rockets",
		"2:04 Kill: 1022 2 19: <world> killed Dono da Bola by MOD_FALLIN",
		//"13:55 score: 6  ping: 19  client: 7 Mal",
		"14:11 ShutdownGame:",
		"14:11 ------------------------------------------------------------",
		"13:55 Exit: Fraglimit hit.",
		" 21:10 ClientDisconnect: 2",
	}

	expectedEvents := []Event{
		{timestamp: 107, tokenType: InitGame, value: InitGameValue{
			Settings: map[string]string{
				"sv_floodProtect":   "1",
				"sv_maxPing":        "0",
				"sv_minPing":        "0",
				"sv_maxRate":        "10000",
				"sv_minRate":        "0",
				"sv_hostname":       "Code Miner Server",
				"g_gametype":        "0",
				"sv_privateClients": "2",
				"sv_maxclients":     "16",
				"sv_allowDownload":  "0",
				"bot_minplayers":    "0",
				"dmflags":           "0",
				"fraglimit":         "20",
				"timelimit":         "15",
				"g_maxGameClients":  "0",
				"capturelimit":      "8",
				"version":           "ioq3 1.36 linux-x86_64 Apr 12 2009",
				"protocol":          "68",
				"mapname":           "q3dm17",
				"gamename":          "baseq3",
				"g_needpass":        "0",
			},
		}},
		{timestamp: 107, tokenType: ClientConnect, value: ClientConnectValue{ClientID: 2}},
		{timestamp: 229, tokenType: ClientBegin, value: ClientBeginValue{ClientID: 5}},
		{timestamp: 107, tokenType: ClientUserinfoChanged, value: ClientUserinfoChangedValue{
			ClientID: 2,
			UserInfo: map[string]string{
				"n":          "Dono da Bola",
				"t":          "0",
				"model":      "sarge",
				"hmodel":     "sarge",
				"g_redteam":  "",
				"g_blueteam": "",
				"c1":         "4",
				"c2":         "5",
				"hc":         "95",
				"w":          "0",
				"l":          "0",
				"tt":         "0",
				"tl":         "0",
			}}},
		{timestamp: 108, tokenType: Item, value: ItemValue{ItemID: 4, ItemName: "ammo_rockets"}},
		{timestamp: 124, tokenType: Kill, value: KillValue{KillerID: 1022, VictimID: 2, DeathCauseID: 19, DeathCauseName: "MOD_FALLIN"}},
		//{timestamp: 835, tokenType: Score, value: ScoreValue{score: 6, ping: 19, ClientID: 7, clientName: "Mal"}},
		{timestamp: 851, tokenType: ShutdownGame, value: ShutdownGameValue{}},
		{timestamp: 851, tokenType: divider, value: nil},
		{timestamp: 835, tokenType: Exit, value: ExitGameValue{Reason: "Fraglimit hit."}},
		{timestamp: 1270, tokenType: ClientDisconnect, value: ClientDisconnectValue{ClientID: 2}},
	}

	for i, line := range linesToTest {
		event, err := ParseLine(line)
		if err != nil {
			t.Errorf("Error while parsing line %d: %v", i, err)
		}

		if event.timestamp != expectedEvents[i].timestamp {
			t.Errorf("Timestamp %d does not match expected timestamp. Expected: %d, got: %d", i, expectedEvents[i].timestamp, event.timestamp)
		}

		if event.tokenType != expectedEvents[i].tokenType {
			t.Errorf("Token type %d does not match expected token type. Expected: %v, got: %v", i, expectedEvents[i].tokenType, event.tokenType)
		}

		eventJSON, err := json.Marshal(event)

		if err != nil {
			t.Errorf("Error while marshalling event %d: %v", i, err)
		}

		expectedEventJSON, err := json.Marshal(expectedEvents[i])

		if err != nil {
			t.Errorf("Error while marshalling expected event %d: %v", i, err)
		}

		if string(eventJSON) != string(expectedEventJSON) {
			t.Errorf("Event %d does not match expected event. Expected: %v, got: %v", i, expectedEvents[i], event)
		}
	}
}
