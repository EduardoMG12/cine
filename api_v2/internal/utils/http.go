package utils

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("Failed to encode JSON response", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func WriteError(w http.ResponseWriter, status int, err error, message string) {
	response := ErrorResponse{
		Error:   err.Error(),
		Message: message,
	}
	WriteJSON(w, status, response)
}

func WriteValidationError(w http.ResponseWriter, err error) {
	details := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			field := fieldError.Field()
			switch fieldError.Tag() {
			case "required":
				details[field] = "This field is required"
			case "email":
				details[field] = "Must be a valid email address"
			case "min":
				details[field] = "Must be at least " + fieldError.Param() + " characters"
			case "max":
				details[field] = "Must be no more than " + fieldError.Param() + " characters"
			case "alphanum":
				details[field] = "Must contain only letters and numbers"
			case "oneof":
				details[field] = "Must be one of: " + fieldError.Param()
			case "url":
				details[field] = "Must be a valid URL"
			default:
				details[field] = "Invalid value"
			}
		}
	}

	response := ErrorResponse{
		Error:   "validation_failed",
		Message: "Request validation failed",
		Details: details,
	}

	WriteJSON(w, http.StatusBadRequest, response)
}

func ParseJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (in case of proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return r.RemoteAddr
}
