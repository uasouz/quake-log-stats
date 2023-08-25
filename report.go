package main

import (
	"embed"
	"fmt"
	"html/template"
	"os"
	"sort"
	"strings"
	"time"
)

func rankMapByValues(kills map[string]int) []string {
	ranking := make([]string, 0)
	for player := range kills {
		ranking = append(ranking, player)
	}

	sort.SliceStable(ranking, func(i, j int) bool {
		return kills[ranking[i]] > kills[ranking[j]]
	})

	return ranking
}

var funcMap = template.FuncMap{
	"rankMapByValues": rankMapByValues,
	"trimPrefix":      strings.TrimPrefix,
}

//go:embed report.gohtml
var fs embed.FS

func renderReport(data map[string]MatchData) (string, error) {
	reportTemplate, err := template.New("report").Funcs(funcMap).ParseFS(fs, "report.gohtml")

	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("report_%d.html", time.Now().Unix())
	reportFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return "", err
	}

	return fileName, reportTemplate.ExecuteTemplate(reportFile, "report.gohtml", data)
}

func printMatchesReport(data map[string]MatchData) {
	for matchID, matchData := range data {
		fmt.Println("------------------------------------------------")
		fmt.Println(fmt.Sprintf("Match %s", strings.TrimPrefix(matchID, "game_")))
		fmt.Println(fmt.Sprintf("Total kills: %d", matchData.TotalKills))
		fmt.Println("Players:")
		for _, player := range matchData.Players {
			fmt.Println(fmt.Sprintf("\t%s: %d", player, matchData.Kills[player]))
		}
		fmt.Println("Kills by means of death:")
		for meanOfDeath, kills := range matchData.KillsByMeansOfDeath {
			fmt.Println(fmt.Sprintf("\t%s: %d", meanOfDeath, kills))
		}
		fmt.Println("Ranking:")
		for i, player := range rankMapByValues(matchData.Kills) {
			fmt.Println(fmt.Sprintf("\t%d. %s: %d", i+1, player, matchData.Kills[player]))
		}
		fmt.Println()
	}
}
