package scans

import (
	"fmt"
	"os"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func RunFullScan(domain string) error {
	utils.LogInfo(fmt.Sprintf("Iniciando scan para: %s", domain))

	baseDomain := strings.TrimPrefix(domain, "www.")
	utils.LogInfo(fmt.Sprintf("Domínio base utilizado: %s", baseDomain))

	outputDir := fmt.Sprintf("backend/resultados/%s_%s", domain, utils.Timestamp())
	os.MkdirAll(outputDir, 0755)

	subdomains, err := EnumerateSubdomains(baseDomain)
	if err != nil {
		return err
	}
	utils.LogSuccess(fmt.Sprintf("Total de subdomínios encontrados: %d", len(subdomains)))

	if ffufOut, err := RunFFUF(baseDomain, outputDir); err == nil {
		subdomains = append(subdomains, ffufOut...)
		subdomains = removeDuplicates(subdomains)
		utils.LogSuccess(fmt.Sprintf("Subdomínios encontrados: %d", len(subdomains)))
	}

	active, err := ProbeActiveSites(subdomains, outputDir)
	if err != nil {
		return fmt.Errorf("erro ao verificar sites ativos: %v", err)
	}
	utils.LogSuccess(fmt.Sprintf("Sites ativos encontrados: %d", len(active)))

	juicy := FilterJuicyTargets(active, outputDir)
	utils.LogSuccess(fmt.Sprintf("Juicy targets encontrados: %d", len(juicy)))

	if err := CaptureScreenshots(active, outputDir); err != nil {
		return fmt.Errorf("erro ao capturar screenshots: %v", err)
	}

	if err := CompareWithPrevious(domain, subdomains, outputDir); err != nil {
		return fmt.Errorf("erro na comparação: %v", err)
	}

	utils.LogSuccess(fmt.Sprintf("Scan concluído: %d subdomínios | %d ativos | %d juicy", len(subdomains), len(active), len(juicy)))
	utils.LogSuccess(fmt.Sprintf("Scan concluído para: %s", domain))
	return nil
}
