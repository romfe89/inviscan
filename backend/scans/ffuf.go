package scans

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/romfe89/inviscan/backend/config"
	"github.com/romfe89/inviscan/backend/utils"
)

type ffufResult struct {
	Results []struct {
		Input  map[string]string `json:"input"`
		Host   string            `json:"host"`
		URL    string            `json:"url"`
		Status int64             `json:"status"`
	} `json:"results"`
}

func RunFFUF(domain string, outputDir string) ([]string, error) {
	cfg := config.GetConfig()
	utils.LogInfo("Enumerando subdomínios com ffuf...")

	wordlist := cfg.Wordlists.Ffuf
	if _, err := os.Stat(wordlist); err != nil {
		utils.LogError(fmt.Sprintf("Wordlist ffuf não encontrada em %s (configurado em wordlists.ffuf).", wordlist))
		return nil, fmt.Errorf("wordlist ffuf não encontrada em %s", wordlist)
	}

	outputFile := filepath.Join(outputDir, "ffuf.json")
	baseDomain := strings.TrimPrefix(domain, "www.")

	cmd := exec.Command(cfg.Tools.Ffuf.Path,
		"-u", fmt.Sprintf("http://FUZZ.%s", baseDomain),
		"-w", wordlist,
		"-mc", cfg.Tools.Ffuf.MatchCodes,
		"-t", strconv.Itoa(cfg.Tools.Ffuf.Threads),
		"-of", "json",
		"-o", outputFile,
	)

	utils.LogInfo("Comando ffuf que será executado:")
	utils.LogInfo(strings.Join(cmd.Args, " "))

	outputBytes, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogWarn(fmt.Sprintf("ffuf falhou ou terminou com avisos. Código de saída pode não ser 0. Erro: %v", err))
		utils.LogWarn(fmt.Sprintf("Saída do ffuf (stdout/stderr):\n%s", string(outputBytes)))
	}

	data, err := os.ReadFile(outputFile)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao ler arquivo de resultado do ffuf %s: %v.", outputFile, err))
		return nil, nil
	}
	if len(data) == 0 {
		utils.LogInfo("Arquivo de resultado do ffuf está vazio (nenhum subdomínio encontrado ou ffuf falhou).")
		return nil, nil
	}

	var parsed ffufResult
	if err := json.Unmarshal(data, &parsed); err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao interpretar JSON do ffuf (%s): %v", outputFile, err))
		utils.LogWarn(fmt.Sprintf("Conteúdo do JSON (início): %s", string(data[:min(len(data), 200)])))
		return nil, nil
	}

	var found []string
	seen := make(map[string]bool)
	for _, r := range parsed.Results {
		var subdomain string
		if r.Host != "" {
			subdomain = r.Host
		} else if sub, ok := r.Input["FUZZ"]; ok && sub != "" {
			subdomain = fmt.Sprintf("%s.%s", sub, baseDomain)
		} else if r.URL != "" {
			urlParts := strings.SplitN(r.URL, "//", 2)
			if len(urlParts) == 2 {
				hostAndPath := strings.SplitN(urlParts[1], "/", 2)
				subdomain = hostAndPath[0]
			}
		}

		if subdomain != "" && strings.HasSuffix(subdomain, "."+baseDomain) && !seen[subdomain] {
			found = append(found, subdomain)
			seen[subdomain] = true
		}
	}

	if len(found) > 0 {
		utils.LogSuccess(fmt.Sprintf("ffuf encontrou %d subdomínios únicos (com status %s)", len(found), cfg.Tools.Ffuf.MatchCodes))
	} else {
		utils.LogInfo("ffuf não encontrou novos subdomínios com os critérios atuais.")
	}

	return found, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
