package scans

import (
	"bytes"
	"fmt"
	"os/exec"
)

func ProbeActiveSites(subdomains []string) ([]string, error) {
	fmt.Println("[+] Verificando sites ativos com httprobe...")

	cmd := exec.Command("httprobe", "-prefer-https")
	var stdin bytes.Buffer
	for _, sub := range subdomains {
		stdin.WriteString(sub + "\n")
	}
	cmd.Stdin = &stdin

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("erro ao executar httprobe: %v", err)
	}

	active := parseLines(stdout.Bytes())
	return active, nil
}
