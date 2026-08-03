// Package main adalah entry point HTTP server.
package main

import (
	"log"
	"net/http"

	"github.com/Mind2Screen-Dev-Team/m2s-vsh-project-backend/internal/handler"
)

func main() {
	// Inisialisasi multiplexer HTTP
	mux := http.NewServeMux()

	// Registrasi route status endpoint
	mux.HandleFunc("GET /api/v1/status", handler.StatusHandler)

	// Jalankan server pada port 8080
	log.Println("Server berjalan pada port 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server gagal dimulai: %v", err)
	}
}