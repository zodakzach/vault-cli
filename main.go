package main

import (
	"data-manager/cmd"
	"data-manager/database"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Determine the database path (defaulting to user's home directory)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Could not get user home directory: %v", err)
	}
	dbPath := filepath.Join(homeDir, "vault.db") // Place it in the home directory

	// Initialize the database
	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Could not initialize the database: %v", err)
	}

	// Execute the commands
	cmd.Execute()
}
