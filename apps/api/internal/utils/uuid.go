package utils

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID v4 string
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateUUIDv7 generates a new UUID v7 string (time-ordered)
// UUID v7 is better for database primary keys as it's sortable by creation time
func GenerateUUIDv7() string {
	return uuid.Must(uuid.NewV7()).String()
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// ParseUUID parses a UUID string and returns a UUID object
func ParseUUID(u string) (uuid.UUID, error) {
	return uuid.Parse(u)
}
