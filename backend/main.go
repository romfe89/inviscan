package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/romfe89/inviscan/backend/handlers"
)

func main() {
	http.HandleFunc("/api/ping", handlers.PingHandler)
	http.HandleFunc("/api/scan", handlers.ScanHandler)
	http.HandleFunc("/api/results", handlers.ResultsHandler)	

	port := "8080"
	fmt.Printf("Servidor backend rodando em http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
