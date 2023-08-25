package main

import (
	"encoding/json"
)

// Match represents a Quake match
// It will aggregate all the kills and players
type Match struct {
	ID                  int
	TotalKills          int
	Players             map[int]string
	Kills               map[int]int
	KillsByMeansOfDeath map[string]int
}

func (m *Match) AggregateEvent(event Event) {
	switch event.tokenType {
	case ClientUserinfoChanged:
		value := event.value.(ClientUserinfoChangedValue)
		m.AddPlayer(value.ClientID, value.UserInfo["n"])
	case Kill:
		value := event.value.(KillValue)
		m.AddKill(value)

	}
}

func (m *Match) AddKill(kill KillValue) {
	m.TotalKills++
	m.KillsByMeansOfDeath[kill.MeanOfDeathName]++
	if kill.KillerID == 1022 {
		m.Kills[kill.VictimID]--
		return
	}

	m.Kills[kill.KillerID]++
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
		ID:                  id,
		Kills:               make(map[int]int),
		Players:             map[int]string{},
		KillsByMeansOfDeath: map[string]int{},
	}
}

type MatchReport struct {
	TotalKills          int            `json:"total_kills"`
	Players             []string       `json:"players"`
	Kills               map[string]int `json:"kills"`
	KillsByMeansOfDeath map[string]int `json:"kills_by_means"`
}

func (m *Match) Report() MatchReport {
	report := MatchReport{
		TotalKills:          m.TotalKills,
		Players:             make([]string, 0),
		Kills:               make(map[string]int),
		KillsByMeansOfDeath: m.KillsByMeansOfDeath,
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

type MeanOfDeath int

func (m MeanOfDeath) String() string {
	return []string{
		"MOD_UNKNOWN",
		"MOD_SHOTGUN",
		"MOD_GAUNTLET",
		"MOD_MACHINEGUN",
		"MOD_GRENADE",
		"MOD_GRENADE_SPLASH",
		"MOD_ROCKET",
		"MOD_ROCKET_SPLASH",
		"MOD_PLASMA",
		"MOD_PLASMA_SPLASH",
		"MOD_RAILGUN",
		"MOD_LIGHTNING",
		"MOD_BFG",
		"MOD_BFG_SPLASH",
		"MOD_WATER",
		"MOD_SLIME",
		"MOD_LAVA",
		"MOD_CRUSH",
		"MOD_TELEFRAG",
		"MOD_FALLING",
		"MOD_SUICIDE",
		"MOD_TARGET_LASER",
		"MOD_TRIGGER_HURT",
		"MOD_NAIL",
		"MOD_CHAINGUN",
		"MOD_PROXIMITY_MINE",
		"MOD_KAMIKAZE",
		"MOD_JUICED",
		"MOD_GRAPPLE",
	}[m]
}

const (
	ModUnknown MeanOfDeath = iota
	ModShotgun
	ModGauntlet
	ModMachinegun
	ModGrenade
	ModGrenadeSplash
	ModRocket
	ModRocketSplash
	ModPlasma
	ModPlasmaSplash
	ModRailgun
	ModLightning
	ModBfg
	ModBfgSplash
	ModWater
	ModSlime
	ModLava
	ModCrush
	ModTelefrag
	ModFalling
	ModSuicide
	ModTargetLaser
	ModTriggerHurt
	ModNail
	ModChaingun
	ModProximityMine
	ModKamikaze
	ModJuiced
	ModGrapple
)
