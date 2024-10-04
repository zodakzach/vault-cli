package cmd

import (
	db "vault-cli/database" 
	"vault-cli/vault"
	"fmt"
	"log"
	"golang.org/x/term" // This package allows for hidden password input
	"syscall"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock the vault",
	Long:  `Unlock the vault by providing the master password.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the master password is set
		if err := db.CheckMasterPasswordSet(); err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Prompt for the password and hide input
		fmt.Print("Enter master password: ")
		passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
		fmt.Println() // To move to the next line after password input

		if err != nil {
			log.Fatal("Error reading password:", err)
		}

		password := string(passwordBytes)

		// Verify the master password
		valid, err := db.VerifyMasterPassword(password)
		if err != nil {
			fmt.Println("Error verifying master password:", err)
			return
		}
		if !valid {
			fmt.Println("Invalid master password. Please try again.")
			return
		}

		// Unlock the vault
		err = vault.UnlockVault()
		if err != nil {
			fmt.Println("Error unlocking the vault:", err)
			return
		}

		fmt.Println("Vault unlocked successfully!")
	},
}

func init() {
	// No need for a password flag anymore since we're prompting interactively
}
