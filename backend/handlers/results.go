package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type ScanResult struct {
	Domain        string   `json:"domain"`
	Timestamp     string   `json:"timestamp"`
	Subdomains    int      `json:"subdomains"`
	ActiveSites   int      `json:"activeSites"`
	JuicyTargets  int      `json:"juicyTargets"`
	Path          string   `json:"path"`
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	dirs, err := os.ReadDir("resultados")
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao ler diretórios: %v", err), 500)
		return
	}

	var results []ScanResult
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}

		parts := strings.SplitN(dir.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}
		domain := parts[0]
		timestamp := parts[1]
		base := filepath.Join("resultados", dir.Name())

		subcount := countLines(filepath.Join(base, "subdomains.txt"))
		actcount := countLines(filepath.Join(base, "active_sites.txt"))
		juicount := countLines(filepath.Join(base, "juicytargets.txt"))

		results = append(results, ScanResult{
			Domain:       domain,
			Timestamp:    timestamp,
			Subdomains:   subcount,
			ActiveSites:  actcount,
			JuicyTargets: juicount,
			Path:         dir.Name(),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp > results[j].Timestamp
	})

	json.NewEncoder(w).Encode(results)
}

func countLines(path string) int {
	f, err := os.ReadFile(path)
	if err != nil {
		return 0
	}
	return len(strings.Split(strings.TrimSpace(string(f)), "\n"))
}
