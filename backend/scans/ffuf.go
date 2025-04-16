package scans

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

type ffufResult struct {
	Results []struct {
		Input map[string]string `json:"input"`
		Host  string            `json:"host"`
	} `json:"results"`
}

func RunFFUF(domain string, outputDir string) ([]string, error) {
	utils.LogInfo("Enumerando subdomínios com ffuf...")

	wordlist := "/snap/seclists/900/Discovery/DNS/subdomains-top1million-5000.txt"
	if _, err := os.Stat(wordlist); err != nil {
		utils.LogWarn(fmt.Sprintf("Wordlist não encontrada: %s", wordlist))
		return nil, nil
	}

	outputFile := filepath.Join(outputDir, "ffuf.json")
	baseDomain := strings.TrimPrefix(domain, "www.")

	cmd := exec.Command("ffuf",
		"-u", fmt.Sprintf("http://FUZZ.%s", baseDomain),
		"-w", wordlist,
		"-mc", "200",
		"-t", "40",
		"-of", "json",
		"-o", outputFile,
	)

	utils.LogInfo("Comando que será executado:")
	utils.LogInfo(strings.Join(cmd.Args, " "))

	if err := cmd.Run(); err != nil {
		utils.LogWarn(fmt.Sprintf("ffuf falhou: %v", err))
		return nil, nil
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao ler %s: %v", outputFile, err))
		return nil, nil
	}
	if len(data) == 0 {
		utils.LogWarn("ffuf não retornou resultados.")
		return nil, nil
	}

	var parsed ffufResult
	if err := json.Unmarshal(data, &parsed); err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao interpretar JSON do ffuf: %v", err))
		return nil, nil
	}

	var found []string
	for _, r := range parsed.Results {
		if r.Host != "" {
			found = append(found, r.Host)
		} else if sub, ok := r.Input["FUZZ"]; ok {
			found = append(found, fmt.Sprintf("%s.%s", sub, baseDomain))
		}
	}

	utils.LogSuccess(fmt.Sprintf("ffuf encontrou %d subdomínios", len(found)))
	return found, nil
}
