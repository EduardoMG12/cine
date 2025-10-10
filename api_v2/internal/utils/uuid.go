package utils

import (
	"strings"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4
func GenerateUUID() string {
	return uuid.New().String()
}

// ParseUUID parses a UUID string and returns error if invalid
func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// NormalizeUUID ensures UUID is in lowercase format
func NormalizeUUID(str string) string {
	return strings.ToLower(strings.TrimSpace(str))
}

// EmptyUUID returns the zero UUID
func EmptyUUID() string {
	return uuid.Nil.String()
}

// IsEmptyUUID checks if UUID is empty/nil
func IsEmptyUUID(str string) bool {
	normalized := NormalizeUUID(str)
	return normalized == "" || normalized == EmptyUUID()
}
