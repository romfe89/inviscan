package scans

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

type ffufResult struct {
	Results []struct {
		Input map[string]string `json:"input"`
	} `json:"results"`
}

func RunFFUF(domain string) ([]string, error) {
	utils.LogInfo("Enumerando subdomínios com ffuf...")

	wordlist := "/snap/seclists/900/Discovery/DNS/subdomains-top1million-5000.txt"
	outputFile := fmt.Sprintf("ffuf_%s.json", strings.ReplaceAll(domain, ".", "_"))

	if _, err := os.Stat(wordlist); err != nil {
		utils.LogWarn(fmt.Sprintf("Wordlist não encontrada: %s", wordlist))
		return nil, nil
	}

	cmd := exec.Command("ffuf",
		"-u", fmt.Sprintf("http://FUZZ.%s", domain),
		"-w", wordlist,
		"-mc", "200",
		"-of", "json",
		"-o", outputFile,
	)

	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		utils.LogWarn(fmt.Sprintf("ffuf falhou: %v", err))
		return nil, nil
	}

	// Verificar se o arquivo de saída foi criado
	data, err := os.ReadFile(outputFile)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Não foi possível ler a saída do ffuf (%s): %v", outputFile, err))
		return nil, nil
	}

	if len(data) == 0 {
		utils.LogWarn("ffuf não retornou resultados.")
		return nil, nil
	}

	var parsed ffufResult
	if err := json.Unmarshal(data, &parsed); err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao interpretar saída JSON do ffuf: %v", err))
		return nil, nil
	}

	var found []string
	for _, r := range parsed.Results {
		if sub, ok := r.Input["FUZZ"]; ok {
			found = append(found, fmt.Sprintf("%s.%s", sub, domain))
		}
	}

	_ = os.Remove(outputFile) // opcional: remover o arquivo gerado

	return found, nil
}
