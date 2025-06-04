package server

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Healthy bool   `json:"healthy"`
	Message string `json:"message"`
}

func CheckOllamaHealth() (bool, string) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false, "Cannot connect to Ollama: " + err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "Ollama returned non-200 status: " + resp.Status
	}
	return true, "Ollama is healthy"
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	healthy, msg := CheckOllamaHealth()
	response := HealthResponse{
		Healthy: healthy,
		Message: msg + " | Timestamp: " + time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
