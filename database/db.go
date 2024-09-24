package database

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("vault.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	// Migrate the schema
	if err := DB.AutoMigrate(&SensitiveData{}, &MasterPassword{}, &VaultState{}); err != nil {
		return err
	}

    // Initialize the vault state in your database
    if err := InitializeVaultState(); err != nil {
        return err
    }

	return nil
}

// Initialize the vault state in your database
func InitializeVaultState() error {
    // Check if there's an existing vault state
    var state VaultState
    err := DB.First(&state).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // If no record exists, create the initial vault state
            state = VaultState{IsLocked: true} // Start with the vault locked
            if err := DB.Create(&state).Error; err != nil {
                return err
            }
        } else {
            return err // Return other errors
        }
    }
    return nil
}

// GetVaultState retrieves the current vault state from the database
func GetVaultState() (bool, error) {
    var state VaultState
    // Query the last (or the first) vault state record from the database
    err := DB.First(&state).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // If no record is found, assume the vault is locked
            return true, nil // Assume locked if no state is found
        }
        return false, err // Return other errors
    }
    return state.IsLocked, nil // Return the vault state (locked or unlocked)
}

// SetVaultState updates the vault's locked state in the database
func SetVaultState(isLocked bool) error {
    var state VaultState

    // Try to find the existing vault state
    err := DB.First(&state).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // If no record exists, create a new state
            state = VaultState{IsLocked: isLocked} // Set the state based on the input
            return DB.Create(&state).Error
        }
        return err // Return other errors
    }

    // If the record exists, update the state
    state.IsLocked = isLocked
    return DB.Save(&state).Error
}

func SetMasterPassword(password string, isMasterPasswordSet bool) error {
    if isMasterPasswordSet {
        // A master password exists, delete the old one
        if err := DB.Delete(&MasterPassword{}).Error; err != nil {
            return fmt.Errorf("failed to delete old master password: %w", err)
        }
    }

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	masterPassword := MasterPassword{HashedPassword: string(hashedPassword)}
	return DB.Create(&masterPassword).Error
}

func VerifyMasterPassword(inputPassword string) (bool, error) {
	var masterPassword MasterPassword
	err := DB.Last(&masterPassword).Error
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(masterPassword.HashedPassword), []byte(inputPassword))
	return err == nil, nil
}

func CheckMasterPasswordSet() error {
    var masterPassword MasterPassword
    if err := DB.First(&masterPassword).Error; err != nil {
        // If we can't find the master password, it means it hasn't been set
        return fmt.Errorf("please set the master password first using 'set-master'")
    }
    return nil
}

func AddSensitiveData(service, identifier, value, idType string) error {
	// Use the utility function to validate and convert idType
	identifierType, err := ParseIdentifierType(idType)
	if err != nil {
		return err
	}

	// Retrieve the hashed master password from the database
	var masterPassword MasterPassword
	if err := DB.First(&masterPassword).Error; err != nil {
		return fmt.Errorf("could not retrieve master password: %v", err)
	}

    key := DeriveAESKey(masterPassword.HashedPassword)

	// Encrypt the value using the hashed master password
	encryptedValue, err := Encrypt(value, key)
	if err != nil {
		return fmt.Errorf("error encrypting sensitive data: %v", err)
	}

	sensitiveData := SensitiveData{
		Service:        service,
		Identifier:     identifier,
		Value:          encryptedValue,
		IdentifierType: identifierType,
	}
	return DB.Create(&sensitiveData).Error
}

func GetSensitiveData(service, identifier string) ([]SensitiveData, error) {
	var sensitiveData []SensitiveData
	query := DB.Where("service = ?", service)

	// If identifier is provided, filter by it
	if identifier != "" {
		query = query.Where("identifier = ?", identifier)
	}

	err := query.Find(&sensitiveData).Error
	if err != nil {
		return nil, err // Return nil slice and the error
	}

	// Retrieve the hashed master password
	var masterPassword MasterPassword
	if err := DB.First(&masterPassword).Error; err != nil {
		return nil, fmt.Errorf("could not retrieve master password: %v", err)
	}

    key := DeriveAESKey(masterPassword.HashedPassword)

	// Decrypt the sensitive data values
	for i, data := range sensitiveData {
		decryptedValue, err := Decrypt(data.Value, key)
		if err != nil {
			return nil, fmt.Errorf("error decrypting sensitive data: %v", err)
		}
		sensitiveData[i].Value = decryptedValue // Replace the encrypted value with the decrypted one
	}

	return sensitiveData, nil
}

func ListSensitiveData(idType string) ([]SensitiveData, error) {
	var entries []SensitiveData
	query := DB

	// If idType is not empty, validate and filter by it
	if idType != "" {
		identifierType, err := ParseIdentifierType(idType)
		if err != nil {
			return nil, err // Return nil slice and the error
		}
		query = query.Where("identifier_type = ?", identifierType)
	}

	err := query.Find(&entries).Error
	if err != nil {
		return nil, err // Return nil slice and the error
	}

	// Retrieve the hashed master password
	var masterPassword MasterPassword
	if err := DB.First(&masterPassword).Error; err != nil {
		return nil, fmt.Errorf("could not retrieve master password: %v", err)
	}

    key := DeriveAESKey(masterPassword.HashedPassword)

	// Decrypt the sensitive data values
	for i, entry := range entries {
		decryptedValue, err := Decrypt(entry.Value, key)
		if err != nil {
			return nil, fmt.Errorf("error decrypting sensitive data: %v", err)
		}
		entries[i].Value = decryptedValue // Replace the encrypted value with the decrypted one
	}

	return entries, nil
}