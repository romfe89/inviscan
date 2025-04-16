package scans

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func CaptureScreenshots(sites []string, outputDir string) error {
	if len(sites) == 0 {
		utils.LogWarn("Nenhum site ativo para capturar com gowitness.")
		return nil
	}

	// Cria subpasta dentro do diretório do scan
	gowitnessDir := filepath.Join(outputDir, "gowitness")
	if err := os.MkdirAll(gowitnessDir, 0755); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao criar diretório gowitness: %v", err))
		return err
	}

	// Escreve targets.txt
	utils.LogInfo("Preparando targets.txt para gowitness...")
	targetsPath := filepath.Join(gowitnessDir, "targets.txt")
	if err := os.WriteFile(targetsPath, []byte(strings.Join(sites, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao escrever targets.txt: %v", err))
		return err
	}

	// Executa gowitness (v3+)
	utils.LogInfo("Executando gowitness scan...")
	cmd := exec.Command("gowitness", "scan", "file", "-f", "targets.txt", "--threads", "4", "--skip-html")
	cmd.Dir = gowitnessDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		utils.LogWarn(fmt.Sprintf("gowitness retornou erro: %s", string(output)))
		return fmt.Errorf("falha ao executar gowitness scan: %v", err)
	}

	utils.LogSuccess(fmt.Sprintf("Capturas salvas em: %s", gowitnessDir))
	return nil
}
