package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/romfe89/inviscan/backend/scans"
	"github.com/romfe89/inviscan/backend/utils"
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
		utils.LogWarn("Método HTTP não permitido")
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var req ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.LogError("Erro ao decodificar JSON da requisição")
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	parsed, err := url.Parse(req.URL)
	target := req.URL
	if err == nil && parsed.Host != "" {
		target = parsed.Host
	}

	utils.LogInfo(fmt.Sprintf("Iniciando scan para: %s", target))

	err = scans.RunFullScan(target)
	if err != nil {
		utils.LogError(fmt.Sprintf("Erro durante o scan de %s: %v", target, err))
		http.Error(w, fmt.Sprintf("Erro ao executar scan: %v", err), 500)
		return
	}

	utils.LogSuccess(fmt.Sprintf("Scan concluído para: %s", target))

	resp := ScanResponse{
		Status:  "ok",
		Message: fmt.Sprintf("Varredura concluída para: %s", target),
	}
	json.NewEncoder(w).Encode(resp)
}
