package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type PingResponse struct {
	Message string `json:"message"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(PingResponse{Message: "Online"})
}

func main() {
	http.HandleFunc("/api/ping", pingHandler)

	port := "8080"
	fmt.Printf("Servidor backend rodando em http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
