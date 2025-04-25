package scans

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/romfe89/inviscan/backend/config"
	"github.com/romfe89/inviscan/backend/utils"
)

func ProbeActiveSites(subdomains []string, outputDir string) ([]string, error) {
	cfg := config.GetConfig()
	utils.LogInfo("Verificando sites ativos com httpx...")

	validSubdomains := filterValidDomains(subdomains)
	if len(validSubdomains) == 0 {
		utils.LogWarn("Nenhum subdomínio válido para sondar com httpx.")
		activeFile := filepath.Join(outputDir, "active_sites.txt")
		os.WriteFile(activeFile, []byte(""), 0644)
		return []string{}, nil
	}

	tempFile, err := os.CreateTemp("", "httpx-input-*.txt")
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao criar arquivo temporário para httpx: %v", err))
		return nil, fmt.Errorf("falha ao criar arquivo temporário: %v", err)
	}
	defer os.Remove(tempFile.Name())

	writer := bufio.NewWriter(tempFile)
	for _, sub := range validSubdomains {
		_, _ = writer.WriteString(sub + "\n")
	}
	writer.Flush()
	tempFile.Close()

	utils.LogInfo(fmt.Sprintf("Executando httpx com a lista: %s", tempFile.Name()))

	httpxArgs := []string{
		"-l", tempFile.Name(),
		"-silent",
		"-threads", strconv.Itoa(cfg.Tools.Httpx.Threads),
		"-prefer-https",
		"-status-code",
		"-mc", cfg.Tools.Httpx.MatchCodes,
	}
	if cfg.Tools.Httpx.Timeout > 0 {
		httpxArgs = append(httpxArgs, "-timeout", strconv.Itoa(cfg.Tools.Httpx.Timeout/1000)) // httpx timeout é em segundos
	}

	cmd := exec.Command(cfg.Tools.Httpx.Path, httpxArgs...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run() // Execução sem timeout
	if err != nil {
		utils.LogWarn(fmt.Sprintf("httpx finalizou com erro (pode ser normal): %v", err))
		errMsg := stderr.String()
		if errMsg != "" {
			utils.LogWarn(fmt.Sprintf("httpx stderr: %s", errMsg))
		}
	}

	var active []string
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			parts := strings.Split(line, " ")
			if len(parts) > 0 && (strings.HasPrefix(parts[0], "http://") || strings.HasPrefix(parts[0], "https://")) {
				active = append(active, parts[0])
			} else if strings.HasPrefix(line, "http://") || strings.HasPrefix(line, "https://") {
				active = append(active, line)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao ler saída do httpx: %v", err))
	}

	activeFile := filepath.Join(outputDir, "active_sites.txt")
	if err := os.WriteFile(activeFile, []byte(strings.Join(active, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao salvar active_sites.txt: %v", err))
	} else {
		utils.LogSuccess(fmt.Sprintf("Salvo %d sites ativos em %s", len(active), activeFile))
	}

	return active, nil
}

func filterValidDomains(domains []string) []string {
	valid := []string{}
	seen := make(map[string]bool)
	for _, d := range domains {
		d = strings.TrimSpace(d)
		if d != "" && !strings.HasPrefix(d, "http") && !strings.Contains(d, "/") && !strings.Contains(d, " ") && !seen[d] {
			if !strings.ContainsAny(d, "0123456789") || strings.ContainsAny(d, "abcdefghijklmnopqrstuvwxyz") {
				valid = append(valid, d)
				seen[d] = true
			}
		}
	}
	return valid
}
