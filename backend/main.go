package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/romfe89/inviscan/backend/config"
	"github.com/romfe89/inviscan/backend/handlers"
	"github.com/romfe89/inviscan/backend/utils"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Falha ao carregar configuração: %v", err)
		os.Exit(1)
	}

	cfg := config.GetConfig()

	if err := config.CheckTools(cfg); err != nil {
		utils.LogError("Erro de validação das ferramentas:")
		log.Fatalf("Erro: %v", err)
	}

	http.HandleFunc("/api/ping", handlers.PingHandler)
	http.HandleFunc("/api/scan", handlers.ScanHandler)
	http.HandleFunc("/api/results", handlers.ResultsHandler)

	port := cfg.Server.Port
	utils.LogSuccess(fmt.Sprintf("Servidor backend rodando em http://localhost:%s", port))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
