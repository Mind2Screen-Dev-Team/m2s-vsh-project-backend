// Package handler menyediakan HTTP handler untuk endpoints API.
package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Version adalah versi aplikasi.
const Version = "0.1.0"

// startTime menyimpan waktu mulai server untuk menghitung uptime.
var startTime = time.Now()

// StatusResponse adalah struktur response untuk endpoint status.
type StatusResponse struct {
	Status       string  `json:"status"`
	Version      string  `json:"version"`
	UptimeSeconds float64 `json:"uptime_seconds"`
}

// StatusHandler menangani request GET /api/v1/status.
// Mengembalikan status, versi, dan uptime server dalam format JSON.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		Status:        "ok",
		Version:       Version,
		UptimeSeconds: time.Since(startTime).Seconds(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}