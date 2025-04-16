package scans

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/romfe89/inviscan/backend/utils"
)

type ffufResult struct {
	Results []struct {
		Input map[string]string `json:"input"`
	} `json:"results"`
}

func RunFFUF(domain string) ([]string, error) {
	utils.LogInfo("Enumerando subdomínios com ffuf...")

	wordlist := "/usr/share/wordlists/dirb/big.txt"
	cmd := exec.Command("ffuf",
		"-u", fmt.Sprintf("http://FUZZ.%s", domain),
		"-w", wordlist,
		"-mc", "200",
		"-o", "ffuf.json",
		"-of", "json",
		"-t", "50")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		utils.LogWarn(fmt.Sprintf("ffuf falhou: %s", stderr.String()))
		return nil, nil // não abortar
	}

	var parsed ffufResult
	if err := json.Unmarshal(stdout.Bytes(), &parsed); err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao interpretar saída do ffuf: %v", err))
		return nil, nil
	}

	var found []string
	for _, r := range parsed.Results {
		if sub, ok := r.Input["FUZZ"]; ok {
			found = append(found, fmt.Sprintf("%s.%s", sub, domain))
		}
	}

	return found, nil
}
