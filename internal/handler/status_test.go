package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestStatusHandler200OK menguji endpoint mengembalikan status code 200.
func TestStatusHandler200OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expect status code 200, got %d", rec.Code)
	}
}

// TestStatusHandlerJSONValid menguji endpoint mengembalikan JSON valid.
func TestStatusHandlerJSONValid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	var resp StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Errorf("response bukan JSON valid: %v", err)
	}
}

// TestStatusHandlerFieldTypes menguji tipe data field response.
func TestStatusHandlerFieldTypes(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	var resp StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal decode JSON: %v", err)
	}

	// Validasi field Status adalah "ok"
	if resp.Status != "ok" {
		t.Errorf("expect status 'ok', got '%s'", resp.Status)
	}

	// Validasi field Version sesuai format
	if resp.Version == "" {
		t.Error("version field kosong")
	}

	// Validasi field UptimeSeconds adalah number positif
	if resp.UptimeSeconds < 0 {
		t.Errorf("uptime_seconds negatif: %f", resp.UptimeSeconds)
	}
}

// TestStatusHandlerContentType menguji header Content-Type application/json.
func TestStatusHandlerContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	contentType := rec.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expect Content-Type 'application/json', got '%s'", contentType)
	}
}

// TestStatusHandlerCORSHeader menguji header CORS Access-Control-Allow-Origin.
func TestStatusHandlerCORSHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	origin := rec.Header().Get("Access-Control-Allow-Origin")
	if origin != "*" {
		t.Errorf("expect Access-Control-Allow-Origin '*', got '%s'", origin)
	}
}

// TestStatusHandlerServicesCount menguji response memuat 3 komponen.
func TestStatusHandlerServicesCount(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	var resp StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal decode JSON: %v", err)
	}
	if len(resp.Services) != 3 {
		t.Errorf("expect 3 services, got %d", len(resp.Services))
	}
}

// TestStatusHandlerServiceIDs menguji urutan dan id komponen.
func TestStatusHandlerServiceIDs(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	var resp StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal decode JSON: %v", err)
	}
	want := []string{"backend-api", "database", "auth"}
	for i, id := range want {
		if i >= len(resp.Services) || resp.Services[i].ID != id {
			t.Errorf("service[%d] expect id %q, got %q", i, id, resp.Services[i].ID)
		}
	}
}

// TestStatusHandlerServiceFields menguji field per komponen.
func TestStatusHandlerServiceFields(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	var resp StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal decode JSON: %v", err)
	}
	validStatus := map[string]bool{"ok": true, "degraded": true, "error": true, "unknown": true}
	for _, s := range resp.Services {
		if !validStatus[s.Status] {
			t.Errorf("service %s status tidak valid: %q", s.ID, s.Status)
		}
		if _, err := time.Parse(time.RFC3339, s.LastChecked); err != nil {
			t.Errorf("service %s last_checked bukan RFC 3339: %v", s.ID, err)
		}
		if s.LatencyMS < 0 {
			t.Errorf("service %s latency_ms negatif: %d", s.ID, s.LatencyMS)
		}
	}
}

// TestStatusHandlerAggregateStatus menguji status agregat = worst-of services.
func TestStatusHandlerAggregateStatus(t *testing.T) {
	cases := []struct {
		name  string
		comps []ComponentStatus
		want  string
	}{
		{"all ok", []ComponentStatus{{Status: "ok"}, {Status: "ok"}}, "ok"},
		{"one degraded", []ComponentStatus{{Status: "ok"}, {Status: "degraded"}}, "degraded"},
		{"one error", []ComponentStatus{{Status: "ok"}, {Status: "degraded"}, {Status: "error"}}, "error"},
		{"unknown only", []ComponentStatus{{Status: "unknown"}}, "unknown"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := aggregateStatus(tc.comps); got != tc.want {
				t.Errorf("aggregateStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}

// TestStatusHandlerBackwardCompat menguji field lama tetap ada.
func TestStatusHandlerBackwardCompat(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/status", nil)
	rec := httptest.NewRecorder()

	StatusHandler(rec, req)

	var resp StatusResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("gagal decode JSON: %v", err)
	}
	if resp.Status == "" || resp.Version == "" || resp.UptimeSeconds < 0 {
		t.Error("field lama status/version/uptime_seconds harus tetap ada")
	}
}
