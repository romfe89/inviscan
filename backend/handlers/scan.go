package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/romfe89/inviscan/backend/scans"
)

type ScanRequest struct {
	URL string `json:"url"`
}

type ScanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func ScanHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	go func(url string) {
		err := scans.RunFullScan(url)
		if err != nil {
			fmt.Printf("[!] Erro ao executar scan para %s: %v\n", url, err)
		}
	}(req.URL)

	resp := ScanResponse{
		Status:  "ok",
		Message: fmt.Sprintf("Varredura iniciada para: %s", req.URL),
	}
	json.NewEncoder(w).Encode(resp)
}
