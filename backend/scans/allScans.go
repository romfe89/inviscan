package scans

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/romfe89/inviscan/backend/config"
	"github.com/romfe89/inviscan/backend/utils"
)

func RunFullScan(domain string) error {
	utils.LogInfo(fmt.Sprintf("Iniciando scan para: %s", domain))
	cfg := config.GetConfig()

	baseDomain := strings.TrimPrefix(domain, "www.")
	utils.LogInfo(fmt.Sprintf("Domínio base utilizado: %s", baseDomain))

	outputDirBase := cfg.OutputDirectoryBase
	outputDir := filepath.Join(outputDirBase, fmt.Sprintf("%s_%s", domain, utils.Timestamp()))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		utils.LogError(fmt.Sprintf("Erro crítico ao criar diretório de saída %s: %v", outputDir, err))
		return fmt.Errorf("falha ao criar diretório de saída: %v", err)
	}
	utils.LogInfo(fmt.Sprintf("Resultados serão salvos em: %s", outputDir))

	subdomains, err := EnumerateSubdomains(baseDomain)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro na enumeração inicial de subdomínios: %v", err))
	}
	utils.LogSuccess(fmt.Sprintf("Enumeração inicial encontrou: %d subdomínios", len(subdomains)))

	if ffufOut, err := RunFFUF(baseDomain, outputDir); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao executar ffuf: %v", err))
	} else if len(ffufOut) > 0 {
		originalCount := len(subdomains)
		subdomains = append(subdomains, ffufOut...)
		subdomains = removeDuplicates(subdomains)
		utils.LogSuccess(fmt.Sprintf("Adicionados %d subdomínios via ffuf (Total: %d)", len(subdomains)-originalCount, len(subdomains)))
	}

	subdomainsFile := filepath.Join(outputDir, "subdomains.txt")
	if err := os.WriteFile(subdomainsFile, []byte(strings.Join(subdomains, "\n")), 0644); err != nil {
		utils.LogWarn(fmt.Sprintf("Erro ao salvar subdomains.txt: %v", err))
	}

	active, err := ProbeActiveSites(subdomains, outputDir)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao verificar sites ativos: %v", err))
	} else {
		utils.LogSuccess(fmt.Sprintf("Sites ativos encontrados: %d", len(active)))
	}

	juicy := FilterJuicyTargets(active, outputDir)

	if err := CaptureScreenshots(active, outputDir); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao capturar screenshots: %v", err))
	}

	if err := CompareWithPrevious(domain, subdomains, outputDir); err != nil {
		utils.LogError(fmt.Sprintf("Erro na comparação com scan anterior: %v", err))
	}

	utils.LogSuccess(fmt.Sprintf("Scan concluído para %s: %d subdomínios | %d ativos | %d juicy", domain, len(subdomains), len(active), len(juicy)))
	return nil
}
