package vault

import (
	"os"
	"testing"

	db "vault-cli/database" // Your package for DB interaction
)

const testDBName = "test_vault.db"

// Setup the database for testing
func setup() error {
	return db.InitDB(testDBName)
}

// Cleanup the test database
func cleanup() {
	_ = os.Remove(testDBName) // Deletes the entire test database file
}

// TestUnlockVault tests the UnlockVault function
func TestUnlockVault(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	defer cleanup()

	// Mocking the initial vault state to locked
	if err := db.SetVaultState(true); err != nil {
		t.Fatalf("failed to set initial vault state: %v", err)
	}

	err := UnlockVault()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	isLocked, err := db.GetVaultState()
	if err != nil {
		t.Errorf("failed to get vault state: %v", err)
	}

	if isLocked {
		t.Error("expected vault to be unlocked, but it is still locked")
	}

	// Test unlocking an already unlocked vault
	err = UnlockVault()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	isLocked, err = db.GetVaultState()
	if err != nil {
		t.Errorf("failed to get vault state: %v", err)
	}

	if isLocked {
		t.Error("expected vault to remain unlocked, but it is locked")
	}
}

// TestLockVault tests the LockVault function
func TestLockVault(t *testing.T) {
	if err := setup(); err != nil {
		t.Fatalf("failed to initialize test database: %v", err)
	}
	defer cleanup()

	// Mocking the initial vault state to unlocked
	if err := db.SetVaultState(false); err != nil {
		t.Fatalf("failed to set initial vault state: %v", err)
	}

	err := LockVault()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	isLocked, err := db.GetVaultState()
	if err != nil {
		t.Errorf("failed to get vault state: %v", err)
	}

	if !isLocked {
		t.Error("expected vault to be locked, but it is still unlocked")
	}
}
