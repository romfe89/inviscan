package scans

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/romfe89/inviscan/backend/config"
	"github.com/romfe89/inviscan/backend/utils"
)

func EnumerateSubdomains(domain string) ([]string, error) {
	cfg := config.GetConfig()
	var results []string
	baseDomain := strings.TrimPrefix(domain, "www.")
	var finalErr error

	utils.LogInfo("Enumerando subdomínios com subfinder...")
	if subfinderOut, err := runTool(cfg.Tools.Subfinder.Path, "-d", baseDomain, "-silent"); err != nil {
		utils.LogWarn(fmt.Sprintf("Falha no subfinder: %v", err))
		finalErr = err
	} else {
		results = append(results, subfinderOut...)
	}

	utils.LogInfo("Enumerando subdomínios com assetfinder...")
	if assetfinderOut, err := runTool(cfg.Tools.Assetfinder.Path, "--subs-only", baseDomain); err != nil {
		utils.LogWarn(fmt.Sprintf("Falha no assetfinder: %v", err))
		finalErr = err
	} else {
		results = append(results, assetfinderOut...)
	}

	utils.LogInfo("Buscando subdomínios via crt.sh...")
	if crtshOut, err := queryCRTSh(baseDomain); err != nil {
		utils.LogWarn(fmt.Sprintf("Falha no crt.sh: %v", err))
		finalErr = err
	} else {
		results = append(results, crtshOut...)
	}

	unique := removeDuplicates(results)
	return unique, finalErr
}

func runTool(name string, args ...string) ([]string, error) {
	cmd := exec.Command(name, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("comando '%s %s' falhou: %v\nStderr: %s", name, strings.Join(args, " "), err, stderr.String())
	}
	return parseLines(stdout.Bytes()), nil
}

func queryCRTSh(domain string) ([]string, error) {
	cfg := config.GetConfig()
	curlCmd := exec.Command(cfg.Tools.Curl.Path, "--compressed", "-s", fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain))
	jqCmd := exec.Command(cfg.Tools.Jq.Path, "-r", ".[].name_value")
	sedCmd := exec.Command(cfg.Tools.Sed.Path, "s/\\*\\.//g")

	utils.LogInfo("Executando pipeline crt.sh: curl | jq | sed")

	jqCmd.Stdin, _ = curlCmd.StdoutPipe()
	sedCmd.Stdin, _ = jqCmd.StdoutPipe()

	var finalOutput bytes.Buffer
	var stderrBuf bytes.Buffer
	sedCmd.Stdout = &finalOutput
	sedCmd.Stderr = &stderrBuf

	if err := sedCmd.Start(); err != nil {
		return nil, fmt.Errorf("falha ao iniciar sed: %v", err)
	}
	if err := jqCmd.Start(); err != nil {
		return nil, fmt.Errorf("falha ao iniciar jq: %v", err)
	}
	if err := curlCmd.Start(); err != nil {
		return nil, fmt.Errorf("falha ao iniciar curl: %v", err)
	}
	if err := curlCmd.Wait(); err != nil {
		return nil, fmt.Errorf("curl falhou: %v", err)
	}
	if err := jqCmd.Wait(); err != nil {
		return nil, fmt.Errorf("jq falhou: %v", err)
	}
	if err := sedCmd.Wait(); err != nil {
		stderrString := stderrBuf.String()
		if stderrString != "" {
			utils.LogWarn(fmt.Sprintf("Sed stderr: %s", stderrString))
		}
		return nil, fmt.Errorf("sed falhou: %v", err)
	}

	return parseLines(finalOutput.Bytes()), nil
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
		valLower := strings.ToLower(val)
		if !seen[valLower] {
			seen[valLower] = true
			unique = append(unique, val)
		}
	}
	return unique
}
