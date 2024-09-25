package cmd

import (
	db "data-manager/database"
	"data-manager/vault" // Make sure to import the vault package for timer control
	"fmt"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock the vault",
	Long:  `Unlock the vault by providing the master password.`,
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")

		if err := db.CheckMasterPasswordSet(); err != nil {
			fmt.Println(err)
			return
		}

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

		err = vault.UnlockVault()
		if err != nil {
			fmt.Println("Error unlocking the vault:", err)
			return
		}
	},
}

func init() {
	unlockCmd.Flags().StringP("password", "p", "", "Master password (required)")
	unlockCmd.MarkFlagRequired("password")
}
