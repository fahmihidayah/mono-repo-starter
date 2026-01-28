package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// TruncateString truncates a string to a maximum length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength]
}

// IsEmpty checks if a string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// SanitizeEmail converts email to lowercase and trims whitespace
func SanitizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// GenerateVerificationToken generates a cryptographically secure email verification token
func GenerateVerificationToken() (string, error) {
	return GenerateSecureToken(32)
}

func GenerateSlug(text string) string {
	cleanText := strings.TrimSpace(text)
	result := strings.ToLower(
		strings.ReplaceAll(cleanText, " ", "-"),
	)
	return result
}
