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
