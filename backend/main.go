package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/romfe89/inviscan/backend/handlers"
	"github.com/romfe89/inviscan/backend/utils"
)

func main() {
	http.HandleFunc("/api/ping", handlers.PingHandler)
	http.HandleFunc("/api/scan", handlers.ScanHandler)
	http.HandleFunc("/api/results", handlers.ResultsHandler)

	port := "8080"
	utils.LogSuccess(fmt.Sprintf("Servidor backend rodando em http://localhost:%s", port))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
