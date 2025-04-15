package scans

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func EnumerateSubdomains(domain string) ([]string, error) {
	var results []string

	fmt.Println("[+] Enumerando subdomínios com subfinder...")
	subfinderOut, err := runTool("subfinder", "-d", domain)
	if err != nil {
		return nil, fmt.Errorf("erro no subfinder: %v", err)
	}
	results = append(results, subfinderOut...)

	fmt.Println("[+] Enumerando subdomínios com assetfinder...")
	assetfinderOut, err := runTool("assetfinder", "--subs-only", domain)
	if err != nil {
		return nil, fmt.Errorf("erro no assetfinder: %v", err)
	}
	results = append(results, assetfinderOut...)

	fmt.Println("[+] Buscando subdomínios via crt.sh...")
	crtshOut, err := queryCRTSh(domain)
	if err != nil {
		return nil, fmt.Errorf("erro no crt.sh: %v", err)
	}
	results = append(results, crtshOut...)

	unique := removeDuplicates(results)
	return unique, nil
}

func runTool(name string, args ...string) ([]string, error) {
	cmd := exec.Command(name, args...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return parseLines(stdout.Bytes()), nil
}

func queryCRTSh(domain string) ([]string, error) {
	cmd := exec.Command("curl", "-s", fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain))
	jq := exec.Command("jq", "-r", ".[].name_value")
	sed := exec.Command("sed", "s/\\*\\.//g")

	// encadeia comandos curl | jq | sed
	curlOut, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	jq.Stdin = bytes.NewReader(curlOut)
	jqOut, err := jq.Output()
	if err != nil {
		return nil, err
	}

	sed.Stdin = bytes.NewReader(jqOut)
	sedOut, err := sed.Output()
	if err != nil {
		return nil, err
	}

	return parseLines(sedOut), nil
}

func parseLines(data []byte) []string {
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lines := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.Contains(line, "@") {
			lines = append(lines, line)
		}
	}
	return lines
}

func removeDuplicates(slice []string) []string {
	seen := map[string]bool{}
	unique := []string{}
	for _, val := range slice {
		if !seen[val] {
			seen[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}
