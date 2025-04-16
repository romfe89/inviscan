package scans

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func CompareWithPrevious(domain string, current []string, outputDir string) error {
	dir := "resultados"
	entries, err := os.ReadDir(dir)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao ler diretório de resultados: %v", err))
		return err
	}

	var lastScan string
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() && strings.HasPrefix(name, domain+"_") && filepath.Join(dir, name) != outputDir {
			if lastScan == "" || name > lastScan {
				lastScan = name
			}
		}
	}

	if lastScan == "" {
		utils.LogWarn("Nenhum scan anterior encontrado para comparação.")
		return nil
	}

	lastSubFile := filepath.Join(dir, lastScan, "subdomains.txt")
	lastSubs, err := readLines(lastSubFile)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao ler último subdomains.txt: %v", err))
		return nil
	}

	newSubs := diffSorted(lastSubs, current)
	utils.LogSuccess(fmt.Sprintf("Novos subdomínios desde %s: %d", lastScan, len(newSubs)))

	for _, s := range newSubs {
		fmt.Println("  +", s)
	}

	return nil
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines, sc.Err()
}

func diffSorted(old, current []string) []string {
	sort.Strings(old)
	sort.Strings(current)

	diff := []string{}
	seen := make(map[string]bool)
	for _, s := range old {
		seen[s] = true
	}
	for _, s := range current {
		if !seen[s] {
			diff = append(diff, s)
		}
	}
	return diff
}
