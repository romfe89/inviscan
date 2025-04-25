package scans

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/romfe89/inviscan/backend/config"
	"github.com/romfe89/inviscan/backend/utils"
)

func CaptureScreenshots(sites []string, outputDir string) error {
	cfg := config.GetConfig()
	if len(sites) == 0 {
		utils.LogWarn("Nenhum site ativo para capturar com gowitness.")
		return nil
	}
	gowitnessParentDir := filepath.Join(cfg.OutputDirectoryBase, filepath.Base(outputDir))
	gowitnessDir := filepath.Join(gowitnessParentDir, "gowitness")
	if err := os.MkdirAll(gowitnessDir, 0755); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao criar diretório gowitness %s: %v", gowitnessDir, err))
		return err
	}

	utils.LogInfo("Preparando targets.txt para gowitness...")
	targetsPath := filepath.Join(gowitnessDir, "targets.txt")
	if err := os.WriteFile(targetsPath, []byte(strings.Join(sites, "\n")), 0644); err != nil {
		utils.LogError(fmt.Sprintf("Erro ao escrever targets.txt em %s: %v", targetsPath, err))
		return err
	}

	utils.LogInfo("Executando gowitness scan...")
	cmd := exec.Command(cfg.Tools.Gowitness.Path,
		"scan", "file", "-f", "targets.txt",
		"--threads", strconv.Itoa(cfg.Tools.Gowitness.Threads),
		"--skip-html",
	)
	cmd.Dir = gowitnessDir

	output, err := cmd.CombinedOutput()
	utils.LogInfo(fmt.Sprintf("Saída do Gowitness:\n%s", string(output)))

	if err != nil {
		utils.LogWarn(fmt.Sprintf("gowitness retornou erro: %v", err))
		return fmt.Errorf("falha ao executar gowitness scan: %v", err)
	}

	relativeGowitnessDir := filepath.Join(filepath.Base(outputDir), "gowitness")
	utils.LogSuccess(fmt.Sprintf("Capturas salvas no diretório: %s", relativeGowitnessDir))
	return nil
}
