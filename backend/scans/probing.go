package scans

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func ProbeActiveSites(subdomains []string) ([]string, error) {
	utils.LogInfo("Verificando sites ativos com httprobe (individual)...")

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

	utils.LogSuccess(fmt.Sprintf("Sites ativos identificados: %d", len(active)))
	return active, nil
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
