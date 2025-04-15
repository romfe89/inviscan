package scans

import (
	"bytes"
	"fmt"
	"os/exec"
)

func CaptureScreenshots(sites []string) error {
	fmt.Println("[+] Inicializando gowitness...")

	// Limpar banco de dados anterior
	exec.Command("rm", "-f", "gowitness.sqlite3").Run()

	// Inicializar novo banco de dados
	if err := exec.Command("gowitness", "init").Run(); err != nil {
		return fmt.Errorf("erro ao inicializar gowitness: %v", err)
	}

	for _, url := range sites {
		fmt.Printf("  ↳ Capturando: %s\n", url)

		cmd := exec.Command("gowitness", "scan", "single", "--url", url, "--write-db")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("    [!] Falha ao capturar %s: %s\n", url, stderr.String())
		}
	}

	fmt.Println("[+] Capturas finalizadas. Use `gowitness server` para visualizar.")
	return nil
}
