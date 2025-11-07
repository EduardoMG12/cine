package http

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/dto"
)

func sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, code, message string) {
	response := dto.APIResponse{
		Success: false,
		Error: &dto.APIError{
			Code:    code,
			Message: message,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
