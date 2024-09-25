package database

import (
	"testing"
	"os"
	"crypto/sha256"
)

// Test Setup and teardown
func setup(filename string) error {
	return InitDB(filename)
}

func teardown(filename string) {
	_ = os.Remove(filename) // Remove the database file if created
}

func TestInitDB(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	teardown(filename)
}

func TestInitializeVaultState(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	state, err := GetVaultState()
	if err != nil {
		t.Fatalf("Failed to get vault state: %v", err)
	}
	if state != true {
		t.Errorf("Expected vault state to be locked, got %v", state)
	}
}

func TestSetVaultState(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	if err := SetVaultState(true); err != nil {
		t.Fatalf("Failed to set vault state to true: %v", err)
	}

	state, err := GetVaultState()
	if err != nil {
		t.Fatalf("Failed to get vault state: %v", err)
	}
	if state != true {
		t.Errorf("Expected vault state to be locked, got %v", state)
	}

	if err := SetVaultState(false); err != nil {
		t.Fatalf("Failed to set vault state to false: %v", err)
	}

	state, err = GetVaultState()
	if err != nil {
		t.Fatalf("Failed to get vault state: %v", err)
	}
	if state != false {
		t.Errorf("Expected vault state to be unlocked, got %v", state)
	}
}

func TestSetMasterPassword(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	err := SetMasterPassword("mysecretpassword", false)
	if err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	// Check if the password is set correctly
	if err := CheckMasterPasswordSet(); err != nil {
		t.Fatalf("Master password should be set: %v", err)
	}
}

func TestVerifyMasterPassword(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	password := "mysecretpassword"
	if err := SetMasterPassword(password, false); err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	valid, err := VerifyMasterPassword(password)
	if err != nil || !valid {
		t.Errorf("Expected valid master password, got error: %v", err)
	}

	valid, err = VerifyMasterPassword("wrongpassword")
	if err == nil && valid {
		t.Error("Expected invalid master password, got valid")
	}
}

func TestAddAndGetSensitiveData(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	if err := SetMasterPassword("mysecretpassword", false); err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	if err := AddSensitiveData("example.com", "user@example.com", "mypassword", "email"); err != nil {
		t.Fatalf("Failed to add sensitive data: %v", err)
	}

	data, err := GetSensitiveData("example.com", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to get sensitive data: %v", err)
	}
	if data.Service != "example.com" || data.Identifier != "user@example.com" {
		t.Errorf("Expected service 'example.com' and identifier 'user@example.com', got %v and %v", data.Service, data.Identifier)
	}
}

func TestDeleteSensitiveData(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	if err := SetMasterPassword("mysecretpassword", false); err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	if err := AddSensitiveData("example.com", "user@example.com", "mypassword", "email"); err != nil {
		t.Fatalf("Failed to add sensitive data: %v", err)
	}

	if err := DeleteSensitiveData("example.com", "user@example.com"); err != nil {
		t.Fatalf("Failed to delete sensitive data: %v", err)
	}

	_, err := GetSensitiveData("example.com", "user@example.com")
	if err == nil {
		t.Error("Expected error getting deleted sensitive data, got none")
	}
}

func TestUpdateSensitiveData(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	if err := SetMasterPassword("mysecretpassword", false); err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	if err := AddSensitiveData("example.com", "user@example.com", "mypassword", "email"); err != nil {
		t.Fatalf("Failed to add sensitive data: %v", err)
	}

	if err := UpdateSensitiveData("example.com", "user@example.com", "newpassword", ""); err != nil {
		t.Fatalf("Failed to update sensitive data: %v", err)
	}

	data, err := GetSensitiveData("example.com", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to get sensitive data: %v", err)
	}
	if data.Value != "newpassword" {
		t.Errorf("Expected updated value 'newpassword', got %v", data.Value)
	}
}

func TestGetAllSensitiveData(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	if err := SetMasterPassword("mysecretpassword", false); err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	if err := AddSensitiveData("example.com", "user@example.com", "mypassword", "email"); err != nil {
		t.Fatalf("Failed to add sensitive data: %v", err)
	}

	allData, err := GetAllSensitiveData("")
	if err != nil {
		t.Fatalf("Failed to get all sensitive data: %v", err)
	}
	if len(allData) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(allData))
	}
}

func TestCheckMasterPasswordSet(t *testing.T) {
	filename := "test_vault.db"
	if err := setup(filename); err != nil {
		t.Fatalf("Failed to initialize DB: %v", err)
	}
	defer teardown(filename)

	err := CheckMasterPasswordSet()
	if err == nil {
		t.Error("Expected error for unset master password, got none")
	}

	if err := SetMasterPassword("mysecretpassword", false); err != nil {
		t.Fatalf("Failed to set master password: %v", err)
	}

	if err := CheckMasterPasswordSet(); err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestParseIdentifierType tests the ParseIdentifierType function
func TestParseIdentifierType(t *testing.T) {
	tests := []struct {
		input    string
		expected IdentifierType
		wantErr  bool
	}{
		{"username", IdentifierTypeUsername, false},
		{"email", IdentifierTypeEmail, false},
		{"api_key", IdentifierTypeAPIKey, false},
		{"secret_key", IdentifierTypeSecret, false},
		{"invalid", "", true}, // expect error for invalid input
	}

	for _, test := range tests {
		result, err := ParseIdentifierType(test.input)
		if (err != nil) != test.wantErr {
			t.Errorf("ParseIdentifierType(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
		}
		if !test.wantErr && result != test.expected {
			t.Errorf("ParseIdentifierType(%q) = %v, want %v", test.input, result, test.expected)
		}
	}
}

// TestDeriveAESKey tests the DeriveAESKey function
func TestDeriveAESKey(t *testing.T) {
	password := "mysecretpassword"
	expected := sha256.Sum256([]byte(password))

	result := DeriveAESKey(password)
	// Compare the results
	if !equal(result, expected[:]) {
		t.Errorf("DeriveAESKey(%q) = %x, want %x", password, result, expected[:])
	}
}

// Helper function to compare byte slices
func equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// TestEncryptDecrypt tests the Encrypt and Decrypt functions
func TestEncryptDecrypt(t *testing.T) {
	key := DeriveAESKey("mysecretpassword") // Deriving the key
	plaintext := "Hello, World!"

	// Encrypt the plaintext
	ciphertextHex, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Failed to encrypt: %v", err)
	}

	// Decrypt the ciphertext
	decryptedText, err := Decrypt(ciphertextHex, key)
	if err != nil {
		t.Fatalf("Failed to decrypt: %v", err)
	}

	if decryptedText != plaintext {
		t.Errorf("Decrypt() = %q, want %q", decryptedText, plaintext)
	}
}