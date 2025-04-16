package scans

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/romfe89/inviscan/backend/utils"
)

func CaptureScreenshots(sites []string) error {
	if len(sites) == 0 {
		utils.LogWarn("Nenhum site ativo para capturar com gowitness.")
		return nil
	}

	// Cria pasta para capturas
	timestamp := time.Now().Format("20060102_150405")
	outputDir := filepath.Join("resultados", fmt.Sprintf("gowitness_%s", timestamp))
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao criar diretório de capturas: %v", err))
		return err
	}

	// Escreve arquivo de alvos
	utils.LogInfo("Preparando targets.txt para gowitness...")
	targetsPath := filepath.Join(outputDir, "targets.txt")
	if err := os.WriteFile(targetsPath, []byte(strings.Join(sites, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao escrever targets.txt: %v", err))
		return err
	}

	// Executa gowitness scan usando subcomando `file`
	utils.LogInfo("Executando gowitness scan...")
	cmd := exec.Command("gowitness", "scan", "file", "-f", "targets.txt", "--threads", "4", "--skip-html")
	cmd.Dir = outputDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogWarn(fmt.Sprintf("gowitness retornou erro: %s", string(output)))
		return fmt.Errorf("falha ao executar gowitness scan: %v", err)
	}

	utils.LogSuccess(fmt.Sprintf("Capturas salvas em: %s", outputDir))
	return nil
}
