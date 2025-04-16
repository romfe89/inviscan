package scans

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func ProbeActiveSites(subdomains []string, outputDir string) ([]string, error) {
	utils.LogInfo("Verificando sites ativos com httprobe...")

	var active []string
	for _, sub := range filterValidDomains(subdomains) {
		cmd := exec.Command("httprobe", "-t", "5000")
		cmd.Stdin = strings.NewReader(sub + "\n")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			utils.LogWarn(fmt.Sprintf("httprobe falhou para %s: %s", sub, stderr.String()))
			continue
		}

		scanner := bufio.NewScanner(&stdout)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				active = append(active, line)
			}
		}
	}

	// Salvar os sites ativos no diretório de saída
	activeFile := filepath.Join(outputDir, "active_sites.txt")
	if err := os.WriteFile(activeFile, []byte(strings.Join(active, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao salvar active_sites.txt: %v", err))
	}

	utils.LogSuccess(fmt.Sprintf("Sites ativos encontrados: %d", len(active)))
	return active, nil
}

func filterValidDomains(domains []string) []string {
	valid := []string{}
	for _, d := range domains {
		// ignorar strings vazias, URLs com protocolo e espaços
		if d != "" && !strings.HasPrefix(d, "http") && !strings.Contains(d, "/") && !strings.Contains(d, " ") {
			valid = append(valid, d)
		}
	}
	return valid
}
