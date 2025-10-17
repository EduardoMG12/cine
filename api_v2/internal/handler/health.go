package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
)

type HealthHandler struct {
	localizer *i18n.Localizer
}

func NewHealthHandler(localizer *i18n.Localizer) *HealthHandler {
	return &HealthHandler{
		localizer: localizer,
	}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Message   string            `json:"message"`
}

func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "2.0.0",
		Services: map[string]string{
			"api":      "healthy",
			"database": "healthy",
		},
		Message: h.localizer.T(r.Context(), "health.status"),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
