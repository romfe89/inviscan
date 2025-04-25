package config

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/romfe89/inviscan/backend/utils"
)

func CheckTools(cfg *Config) error {
	utils.LogInfo("Verificando existência das ferramentas configuradas...")

	toolsToCheck := map[string]string{
		"subfinder":   cfg.Tools.Subfinder.Path,
		"assetfinder": cfg.Tools.Assetfinder.Path,
		"curl":        cfg.Tools.Curl.Path,
		"jq":          cfg.Tools.Jq.Path,
		"sed":         cfg.Tools.Sed.Path,
		"ffuf":        cfg.Tools.Ffuf.Path,
		"httpx":       cfg.Tools.Httpx.Path,
		"gowitness":   cfg.Tools.Gowitness.Path,
	}

	var missingTools []string

	for name, path := range toolsToCheck {
		if path == "" {
			missingTools = append(missingTools, fmt.Sprintf("%s (path não configurado!)", name))
			continue
		}
		_, err := exec.LookPath(path)
		if err != nil {
			utils.LogWarn(fmt.Sprintf("Ferramenta '%s' NÃO encontrada no path configurado ('%s'): %v", name, path, err))
			missingTools = append(missingTools, fmt.Sprintf("%s (em '%s')", name, path))
		} else {
			
		}
	}

	if len(missingTools) > 0 {
		return fmt.Errorf("as seguintes ferramentas essenciais não foram encontradas ou não são executáveis: %s. Verifique a instalação e as configurações 'tools.<tool>.path' no seu config.yaml ou variáveis de ambiente", strings.Join(missingTools, ", "))
	}

	utils.LogSuccess("Todas as ferramentas essenciais foram encontradas.")
	return nil
}
