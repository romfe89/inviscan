package scans

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/romfe89/inviscan/backend/utils"
)

func RunFullScan(domain string) error {
	utils.LogInfo(fmt.Sprintf("Iniciando scan para: %s", domain))

	// Cria diretório de saída
	timestamp := time.Now().Format("20060102_150405")
	outputDir := filepath.Join("resultados", fmt.Sprintf("%s_%s", domain, timestamp))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao criar diretório de saída: %v", err))
		return err
	}

	// Executa etapas
	subdomains, err := EnumerateSubdomains(domain, outputDir)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro na enumeração: %v", err))
		return err
	}
	
	if ffufOut, err := RunFFUF(domain, outputDir); err == nil {
		subdomains = append(subdomains, ffufOut...)
		subdomains = removeDuplicates(subdomains)
	}
	
	utils.LogSuccess(fmt.Sprintf("Subdomínios encontrados: %d", len(subdomains)))

	active, err := ProbeActiveSites(subdomains, outputDir)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao verificar sites ativos: %v", err))
		return err
	}
	utils.LogSuccess(fmt.Sprintf("Sites ativos encontrados: %d", len(active)))

	juicy := FilterJuicyTargets(active, outputDir)
	utils.LogSuccess(fmt.Sprintf("Juicy targets encontrados: %d", len(juicy)))

	err = CaptureScreenshots(active, outputDir)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao capturar screenshots: %v", err))
		return err
	}

	err = CompareWithPrevious(domain, subdomains, outputDir)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Erro na comparação com o scan anterior: %v", err))
	}

	utils.LogSuccess(fmt.Sprintf("Scan concluído: %d subdomínios | %d ativos | %d juicy",
		len(subdomains), len(active), len(juicy)))

	return nil
}
