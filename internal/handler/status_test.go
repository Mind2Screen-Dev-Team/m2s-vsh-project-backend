package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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