package scans

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func CompareWithPrevious(domain string, current []string) error {
	dir := "resultados"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de resultados: %v", err)
	}

	var lastScan string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), domain+"_") {
			lastScan = entry.Name()
		}
	}

	if lastScan == "" {
		fmt.Println("[!] Nenhum scan anterior encontrado para comparação.")
		return nil
	}

	lastSubFile := filepath.Join(dir, lastScan, "subdomains.txt")
	lastSubs, err := readLines(lastSubFile)
	if err != nil {
		fmt.Printf("[!] Erro ao ler último subdomains.txt: %v\n", err)
		return nil
	}

	newSubs := diffSorted(lastSubs, current)
	fmt.Printf("[+] Novos subdomínios desde %s: %d\n", lastScan, len(newSubs))
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
	// Garantir ordenação
	sort.Strings(old)
	sort.Strings(current)

	diff := []string{}
	m := map[string]bool{}
	for _, s := range old {
		m[s] = true
	}
	for _, s := range current {
		if !m[s] {
			diff = append(diff, s)
		}
	}
	return diff
}
