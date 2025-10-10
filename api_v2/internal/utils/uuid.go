package utils

import (
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
