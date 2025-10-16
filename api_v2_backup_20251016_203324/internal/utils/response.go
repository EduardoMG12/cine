package utils

import (
	"encoding/json"
	"net/http"

	"github.com/EduardoMG12/cine/api_v2/internal/i18n"
	"github.com/EduardoMG12/cine/api_v2/internal/middleware"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type Meta struct {
	Page       int `json:"page,omitempty"`
	PageSize   int `json:"page_size,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
	TotalItems int `json:"total_items,omitempty"`
}

// WriteJSONResponse writes a successful JSON response with i18n support
func WriteJSONResponse(w http.ResponseWriter, r *http.Request, statusCode int, data interface{}) {
	response := APIResponse{
		Success: true,
		Data:    data,
	}

	writeResponse(w, statusCode, response)
}

// WriteErrorResponse writes an error JSON response with i18n support
func WriteErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, errorCode string, details ...string) {
	lang := middleware.GetLanguageFromContext(r.Context())
	message := i18n.T(errorCode, lang)

	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    errorCode,
			Message: message,
			Details: getDetails(details),
		},
	}

	writeResponse(w, statusCode, response)
}

// WriteValidationErrorResponse writes validation error with i18n support
func WriteValidationErrorResponse(w http.ResponseWriter, r *http.Request, validationErrors map[string]string) {
	lang := middleware.GetLanguageFromContext(r.Context())

	// Translate validation errors
	translatedErrors := make(map[string]string)
	for field, errorKey := range validationErrors {
		translatedErrors[field] = i18n.T(errorKey, lang)
	}

	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    "error.validation_failed",
			Message: i18n.T("error.validation_failed", lang),
			Details: mapToString(translatedErrors),
		},
	}

	writeResponse(w, http.StatusBadRequest, response)
}

// WriteSuccessMessage writes a success message response with i18n support
func WriteSuccessMessage(w http.ResponseWriter, r *http.Request, messageKey string, data ...interface{}) {
	lang := middleware.GetLanguageFromContext(r.Context())
	message := i18n.T(messageKey, lang)

	response := APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": message,
			"data":    getOptionalData(data),
		},
	}

	writeResponse(w, http.StatusOK, response)
}

// WritePaginatedResponse writes paginated response with meta information
func WritePaginatedResponse(w http.ResponseWriter, r *http.Request, data interface{}, page, pageSize, totalItems int) {
	totalPages := (totalItems + pageSize - 1) / pageSize

	response := APIResponse{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			TotalItems: totalItems,
		},
	}

	writeResponse(w, http.StatusOK, response)
}

func writeResponse(w http.ResponseWriter, statusCode int, response APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func getDetails(details []string) string {
	if len(details) > 0 {
		return details[0]
	}
	return ""
}

func getOptionalData(data []interface{}) interface{} {
	if len(data) > 0 {
		return data[0]
	}
	return nil
}

func mapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	bytes, _ := json.Marshal(m)
	return string(bytes)
}
