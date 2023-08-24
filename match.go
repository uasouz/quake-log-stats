package main

import (
	"encoding/json"
)

// Match represents a Quake match
// It will aggregate all the kills and players
type Match struct {
	ID         int
	TotalKills int
	Players    map[int]string
	Kills      map[int]int
}

func (m *Match) AddKill(killerID, killedID int) {
	m.TotalKills++
	if killerID == 1022 {
		m.Kills[killedID]--
		return
	}

	m.Kills[killerID]++
}

func (m *Match) AddPlayer(playerID int, player string) {
	if player == "<world>" {
		return
	}

	m.Players[playerID] = player
	if _, ok := m.Kills[playerID]; !ok {
		m.Kills[playerID] = 0
	}
}

func NewMatch(id int) *Match {
	return &Match{
		ID:      id,
		Kills:   make(map[int]int),
		Players: map[int]string{},
	}
}

type MatchReport struct {
	TotalKills int `json:"total_kills"`
	Players    []string
	Kills      map[string]int
}

func (m *Match) Report() MatchReport {
	report := MatchReport{
		TotalKills: m.TotalKills,
		Players:    make([]string, 0),
		Kills:      make(map[string]int),
	}

	for _, player := range m.Players {
		report.Players = append(report.Players, player)
	}

	for playerID, kills := range m.Kills {
		report.Kills[m.Players[playerID]] = kills
	}

	return report
}

type Matches []*Match

func (m Matches) String() string {
	s, _ := json.Marshal(m)

	return string(s)
}
