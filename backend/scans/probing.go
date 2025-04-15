package scans

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func ProbeActiveSites(subdomains []string) ([]string, error) {
	utils.LogInfo("Verificando sites ativos com httprobe...")

	cleaned := filterValidDomains(subdomains)
	if len(cleaned) == 0 {
		utils.LogWarn("Nenhum domínio válido para testar com httprobe.")
		return nil, fmt.Errorf("nenhum domínio válido para testar")
	}

	cmd := exec.Command("httprobe", "-prefer-https", "-t", "2000")	
	var stdin bytes.Buffer
	for _, sub := range cleaned {
		stdin.WriteString(sub + "\n")
	}
	cmd.Stdin = &stdin

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao executar httprobe: %v", err))
		return nil, fmt.Errorf("erro ao executar httprobe: %v", err)
	}

	utils.LogSuccess(fmt.Sprintf("Sites ativos identificados: %d", len(parseLines(stdout.Bytes()))))
	return parseLines(stdout.Bytes()), nil
}

func filterValidDomains(domains []string) []string {
	valid := []string{}
	for _, d := range domains {
		if d != "" && !strings.HasPrefix(d, "http") && !strings.Contains(d, "/") && !strings.Contains(d, " ") {
			valid = append(valid, d)
		}
	}
	return valid
}
