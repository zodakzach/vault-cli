package cmd

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomPassword generates a secure random password
func generateRandomPassword(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	// Use base64 encoding to generate a readable password
	return base64.RawStdEncoding.EncodeToString(bytes)[:length], nil
}

