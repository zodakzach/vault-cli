package database

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"crypto/sha256"
)

// ParseIdentifierType attempts to convert a string to IdentifierType
func ParseIdentifierType(idType string) (IdentifierType, error) {
	switch idType {
	case "username":
		return IdentifierTypeUsername, nil
	case "email":
		return IdentifierTypeEmail, nil
	case "api_key":
		return IdentifierTypeAPIKey, nil
	case "secret_key":
		return IdentifierTypeSecret, nil
	default:
		return "", fmt.Errorf("invalid identifier type: %s", idType)
	}
}

// DeriveAESKey derives a 32-byte AES key from the bcrypt-hashed password using SHA-256
func DeriveAESKey(hashedPassword string) []byte {
	hash := sha256.Sum256([]byte(hashedPassword))
	return hash[:]
}

// Encrypt encrypts the given plaintext using the provided key
func Encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a nonce (number used once)
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCFBEncrypter(block, nonce)
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// Return the nonce + ciphertext as a hex string
	return hex.EncodeToString(nonce) + hex.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given ciphertext using the provided key
func Decrypt(ciphertextHex string, key []byte) (string, error) {
	// Decode the hex string
	data, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return "", err
	}

	if len(data) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Separate the nonce and the actual ciphertext
	nonce, ciphertext := data[:aes.BlockSize], data[aes.BlockSize:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Decrypt the data
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCFBDecrypter(block, nonce)
	stream.XORKeyStream(plaintext, ciphertext)

	return string(plaintext), nil
}
