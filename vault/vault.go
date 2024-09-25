package vault

import (
	db "data-manager/database" // Your package for DB interaction
	"fmt"
	"sync"
)

var (
	lockMutex sync.Mutex
)

// UnlockVault unlocks the vault and updates its state in the database
func UnlockVault() error {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	// Check if the vault is already unlocked
	isLocked, err := db.GetVaultState()
	if err != nil {
		return fmt.Errorf("failed to get vault state: %v", err)
	}

	if !isLocked {
		fmt.Println("Vault is already unlocked.")
		return nil
	}

	// Unlock the vault
	err = db.SetVaultState(false)
	if err != nil {
		return fmt.Errorf("failed to unlock the vault: %v", err)
	}

	fmt.Println("Vault unlocked.")
	return nil
}

// LockVault locks the vault and updates its state in the database
func LockVault() error {
	lockMutex.Lock()
	defer lockMutex.Unlock()

	// Lock the vault manually
	err := db.SetVaultState(true)
	if err != nil {
		return fmt.Errorf("failed to lock the vault: %v", err)
	}

	fmt.Println("Vault locked.")
	return nil
}
