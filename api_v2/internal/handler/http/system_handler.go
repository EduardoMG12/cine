package http

import (
	"net/http"
)

type SystemHandler struct{}

func NewSystemHandler() *SystemHandler {
	return &SystemHandler{}
}

// Root godoc
// @Summary API Root
// @Description Welcome endpoint with documentation and health check links
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "Welcome message with links"
// @Router / [get]
func (h *SystemHandler) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Welcome to CineVerse API","documentation":"http://localhost:8080/swagger/index.html","health":"/health"}`))
}

// HealthCheck godoc
// @Summary Health Check
// @Description Check if the API is running and healthy
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string "Health status"
// @Router /health [get]
func (h *SystemHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"cineverse-api"}`))
}
