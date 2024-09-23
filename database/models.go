package database

import "gorm.io/gorm"

type IdentifierType string

const (
	IdentifierTypeUsername IdentifierType = "username"
	IdentifierTypeEmail    IdentifierType = "email"
	IdentifierTypeAPIKey   IdentifierType = "api_key"
	IdentifierTypeSecret   IdentifierType = "secret_key"
)

type SensitiveData struct {
	gorm.Model
	Service        string         `gorm:"index:idx_service_identifier,unique"`
	Identifier     string         `gorm:"index:idx_service_identifier,unique"` // can be username, email, API key, etc.
	Value          string         // this could be the actual password, API key, or sensitive value
	IdentifierType IdentifierType // type of identifier (e.g., username, email, API key)
}

type MasterPassword struct {
	gorm.Model
	HashedPassword string `gorm:"uniqueIndex"` // Store the hashed password
}

// VaultState represents the state of the vault (locked or unlocked)
type VaultState struct {
    gorm.Model
    IsLocked bool `gorm:"default:true"` // Default to true (locked)
}
