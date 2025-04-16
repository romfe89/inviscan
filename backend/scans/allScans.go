package scans

import (
	"fmt"

	"github.com/romfe89/inviscan/backend/utils"
)

func RunFullScan(domain string) error {
	utils.LogInfo(fmt.Sprintf("Iniciando scan para: %s", domain))

	subdomains, err := EnumerateSubdomains(domain)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro na enumeração de subdomínios: %v", err))
		return err
	}
	utils.LogSuccess(fmt.Sprintf("Subdomínios encontrados: %d", len(subdomains)))

	active, err := ProbeActiveSites(subdomains)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao verificar sites ativos: %v", err))
		return err
	}
	utils.LogSuccess(fmt.Sprintf("Sites ativos encontrados: %d", len(active)))

	juicy := FilterJuicyTargets(active)
	utils.LogSuccess(fmt.Sprintf("Juicy targets encontrados: %d", len(juicy)))

	err = CaptureScreenshots(active)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro ao capturar screenshots: %v", err))
		return err
	}

	err = CompareWithPrevious(domain, subdomains)
	if err != nil {
		utils.LogWarn(fmt.Sprintf("Erro na comparação com o scan anterior: %v", err))
	}

	utils.LogSuccess(fmt.Sprintf("Scan concluído: %d subdomínios | %d ativos | %d juicy",
		len(subdomains), len(active), len(juicy),
	))
	return nil
}
