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

// ComponentStatus adalah status satu komponen (CONTRACT-201).
type ComponentStatus struct {
	ID          string  `json:"id"`
	Status      string  `json:"status"` // ok|degraded|error|unknown
	Version     string  `json:"version"`
	Uptime      float64 `json:"uptime"`       // detik
	LastChecked string  `json:"last_checked"` // RFC 3339
	LatencyMS   int     `json:"latency_ms"`
	Message     string  `json:"message"`
}

// StatusResponse adalah struktur response untuk endpoint status.
// Backward-compat (CONTRACT-201): field lama dipertahankan, ditambah services.
type StatusResponse struct {
	Status        string            `json:"status"`
	Version       string            `json:"version"`
	UptimeSeconds float64           `json:"uptime_seconds"`
	Services      []ComponentStatus `json:"services"`
}

// probeBackendAPI melaporkan status komponen backend-api (diri sendiri).
func probeBackendAPI() ComponentStatus {
	return ComponentStatus{
		ID:          "backend-api",
		Status:      "ok",
		Version:     Version,
		Uptime:      time.Since(startTime).Seconds(),
		LastChecked: time.Now().Format(time.RFC3339),
		LatencyMS:   0,
		Message:     "",
	}
}

// probeDatabase melaporkan status komponen database via probe sintetik.
// Tanpa dependensi DB nyata; selalu ok sampai ada impl nyata.
func probeDatabase() ComponentStatus {
	start := time.Now()
	time.Sleep(2 * time.Millisecond)
	return ComponentStatus{
		ID:          "database",
		Status:      "ok",
		Version:     "",
		Uptime:      0,
		LastChecked: time.Now().Format(time.RFC3339),
		LatencyMS:   int(time.Since(start).Milliseconds()),
		Message:     "synthetic ping ok",
	}
}

// probeAuth melaporkan status komponen auth via probe sintetik.
// Tanpa service auth nyata; selalu ok sampai ada impl nyata.
func probeAuth() ComponentStatus {
	start := time.Now()
	time.Sleep(1 * time.Millisecond)
	return ComponentStatus{
		ID:          "auth",
		Status:      "ok",
		Version:     Version,
		Uptime:      time.Since(startTime).Seconds(),
		LastChecked: time.Now().Format(time.RFC3339),
		LatencyMS:   int(time.Since(start).Milliseconds()),
		Message:     "synthetic ok",
	}
}

// aggregateStatus menghitung status agregat = worst-of services.
// Urutan keparahan: error > degraded > ok > unknown.
func aggregateStatus(comps []ComponentStatus) string {
	severity := map[string]int{"error": 3, "degraded": 2, "ok": 1, "unknown": 0}
	worst := "unknown"
	for _, c := range comps {
		if severity[c.Status] > severity[worst] {
			worst = c.Status
		}
	}
	return worst
}

// StatusHandler menangani request GET /api/v1/status.
// Mengembalikan status agregat, versi, uptime, dan status per komponen.
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	services := []ComponentStatus{
		probeBackendAPI(),
		probeDatabase(),
		probeAuth(),
	}
	response := StatusResponse{
		Status:        aggregateStatus(services),
		Version:       Version,
		UptimeSeconds: time.Since(startTime).Seconds(),
		Services:      services,
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
