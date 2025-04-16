package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

type ScanResult struct {
	Domain       string `json:"domain"`
	Timestamp    string `json:"timestamp"`
	Subdomains   int    `json:"subdomains"`
	ActiveSites  int    `json:"activeSites"`
	JuicyTargets int    `json:"juicyTargets"`
	Path         string `json:"path"`
}

func ResultsHandler(w http.ResponseWriter, r *http.Request) {
	utils.LogInfo("Recebida requisição para /api/results")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	entries, err := os.ReadDir("resultados")
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao ler diretório de resultados: %v", err))
		http.Error(w, "Erro ao ler resultados", http.StatusInternalServerError)
		return
	}

	var results []ScanResult

	for _, entry := range entries {
		if !entry.IsDir() || strings.HasPrefix(entry.Name(), "gowitness_") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) != 2 {
			continue
		}
		domain := parts[0]
		timestamp := parts[1]

		base := filepath.Join("resultados", entry.Name())

		subcount := countLines(filepath.Join(base, "subdomains.txt"))
		actcount := countLines(filepath.Join(base, "active_sites.txt"))
		juicount := countLines(filepath.Join(base, "juicytargets.txt"))

		results = append(results, ScanResult{
			Domain:       domain,
			Timestamp:    timestamp,
			Subdomains:   subcount,
			ActiveSites:  actcount,
			JuicyTargets: juicount,
			Path:         entry.Name(),
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp > results[j].Timestamp
	})

	json.NewEncoder(w).Encode(results)
}

func countLines(path string) int {
	content, err := os.ReadFile(path)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Arquivo ausente ou erro: %s", path))
		return 0
	}
	return len(strings.Split(strings.TrimSpace(string(content)), "\n"))
}
