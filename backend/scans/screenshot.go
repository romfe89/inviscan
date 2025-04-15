package scans

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/romfe89/inviscan/backend/utils"
)

func CaptureScreenshots(sites []string) error {
	utils.LogInfo("Inicializando gowitness...")

	// Limpar banco de dados anterior
	exec.Command("rm", "-f", "gowitness.sqlite3").Run()

	// Inicializar novo banco de dados
	if err := exec.Command("gowitness", "init").Run(); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao inicializar gowitness: %v", err))
		return fmt.Errorf("erro ao inicializar gowitness: %v", err)
	}

	for _, url := range sites {
		utils.LogInfo(fmt.Sprintf("↳ Capturando: %s", url))

		cmd := exec.Command("gowitness", "scan", "single", "--url", url, "--write-db")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			utils.LogWarn(fmt.Sprintf("Falha ao capturar %s: %s", url, stderr.String()))
		}
	}

	utils.LogSuccess("Capturas finalizadas. Use `gowitness server` para visualizar.")
	return nil
}
