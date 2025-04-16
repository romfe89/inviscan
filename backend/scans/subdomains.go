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

func EnumerateSubdomains(domain string, outputDir string) ([]string, error) {
	var results []string

	utils.LogInfo("Enumerando subdomínios com subfinder...")
	if subfinderOut, err := runTool("subfinder", "-d", domain); err != nil {
		utils.LogWarn(fmt.Sprintf("Falha no subfinder: %v", err))
	} else {
		results = append(results, subfinderOut...)
	}

	utils.LogInfo("Enumerando subdomínios com assetfinder...")
	if assetfinderOut, err := runTool("assetfinder", "--subs-only", domain); err != nil {
		utils.LogWarn(fmt.Sprintf("Falha no assetfinder: %v", err))
	} else {
		results = append(results, assetfinderOut...)
	}

	utils.LogInfo("Buscando subdomínios via crt.sh...")
	if crtshOut, err := queryCRTSh(domain); err != nil {
		utils.LogWarn(fmt.Sprintf("Falha no crt.sh: %v", err))
	} else {
		results = append(results, crtshOut...)
	}

	unique := removeDuplicates(results)

	// Salvar resultado no diretório do scan
	subFile := filepath.Join(outputDir, "subdomains.txt")
	if err := os.WriteFile(subFile, []byte(strings.Join(unique, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao salvar subdomains.txt: %v", err))
	}

	utils.LogSuccess(fmt.Sprintf("Total de subdomínios encontrados: %d", len(unique)))
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
	cmd := exec.Command("curl", "--compressed", "-s", fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain))
	jq := exec.Command("jq", "-r", ".[].name_value")
	sed := exec.Command("sed", "s/\\*\\.//g")

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
	var lines []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.Contains(line, "@") {
			lines = append(lines, line)
		}
	}
	return lines
}

func removeDuplicates(slice []string) []string {
	seen := make(map[string]bool)
	var unique []string
	for _, val := range slice {
		if !seen[val] {
			seen[val] = true
			unique = append(unique, val)
		}
	}
	return unique
}
