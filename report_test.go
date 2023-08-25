package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"
)

func TestMatchReport(t *testing.T) {
	t.Run("Test MatchReport", func(t *testing.T) {
		match := NewMatch(1)

		match.AggregateEvent(Event{
			timestamp: 0,
			tokenType: InitGame,
			value: InitGameValue{
				Settings: map[string]string{
					"timelimit": "15",
					"fraglimit": "0",
				},
			},
		})

		for id, player := range playerNames {
			match.AggregateEvent(Event{
				timestamp: 0,
				tokenType: ClientUserinfoChanged,
				value: ClientUserinfoChangedValue{
					ClientID: id + 1,
					UserInfo: map[string]string{
						"n":          player,
						"t":          "0",
						"model":      "sarge",
						"hmodel":     "sarge",
						"g_redteam":  "",
						"g_blueteam": "",
						"c1":         "4",
					},
				},
			})
		}

		for i := 0; i < rand.Intn(45); i++ {
			match.AggregateEvent(generateKillEvent(i*60 + rand.Intn(60)))
		}

		matchReport := match.Report()

		matchReportJSON, err := json.Marshal(matchReport)

		if err != nil {
			t.Errorf("Error marshaling matchReport: %v", err)
		}

		fmt.Println(string(matchReportJSON))

	})
}

var playerNames = [16]string{
	"Mal",
	"Zeh",
	"Dono da Bola",
	"Assasinu Credi",
	"Isgalamido",
	"Oootsimo",
	"Oz",
	"Neuromancer",
	"Fatality",
	"Vovo Zilda",
	"Chessus",
	"Fatal Error",
	"Mr. Bean",
	"Chuck Norris",
	"Galactus",
	"Spawn",
}

func generateKillEvent(timestamp int) Event {

	killerID := rand.Intn(16)
	victimID := rand.Intn(15)

	if killerID == 16 {
		return Event{
			timestamp: timestamp,
			tokenType: Kill,
			value: KillValue{
				KillerName:      "<world>",
				KillerID:        1022,
				VictimName:      playerNames[victimID],
				VictimID:        victimID,
				MeanOfDeathID:   int(ModTriggerHurt),
				MeanOfDeathName: ModTriggerHurt.String(),
			},
		}
	}

	if killerID == victimID {
		return Event{
			timestamp: timestamp,
			tokenType: Kill,
			value: KillValue{
				KillerName:      playerNames[killerID],
				KillerID:        killerID,
				VictimName:      playerNames[victimID],
				VictimID:        victimID,
				MeanOfDeathID:   int(ModTriggerHurt),
				MeanOfDeathName: ModTriggerHurt.String(),
			},
		}
	}

	MeanOfDeathID := rand.Intn(26)

	return Event{
		timestamp: timestamp,
		tokenType: Kill,
		value: KillValue{
			KillerName:      playerNames[killerID],
			KillerID:        killerID,
			VictimName:      playerNames[victimID],
			VictimID:        victimID,
			MeanOfDeathID:   MeanOfDeathID,
			MeanOfDeathName: MeanOfDeath(MeanOfDeathID).String(),
		},
	}
}
