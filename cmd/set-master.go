package cmd

import (
	db "vault-cli/database"
	"fmt"
	"github.com/spf13/cobra"
)

// setMasterCmd represents the set-master command
var setMasterCmd = &cobra.Command{
	Use:   "set-master",
	Short: "Set or update the master password",
	Long:  `Set or update the master password for accessing the password vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		masterPassword, _ := cmd.Flags().GetString("password")
		isMasterPasswordSet := false

		if err := db.CheckMasterPasswordSet(); err == nil {
			oldMasterPassword, _ := cmd.Flags().GetString("old-password")
			if oldMasterPassword == "" {
				fmt.Println("Error: Master password already set. Old master password is required to change the master password. Use --old-password")
				return
			}
			valid, err := db.VerifyMasterPassword(oldMasterPassword)
			if err != nil {
				fmt.Println("Error verifying old master password:", err)
				return
			}
			if !valid {
				fmt.Println("Invalid old master password. Please try again.")
				return
			}
			isMasterPasswordSet = true
		}

		// Handle the logic to set the master password
		err := db.SetMasterPassword(masterPassword, isMasterPasswordSet)
		if err != nil {
			fmt.Println("Error setting master password:", err)
			return
		}

		fmt.Println("Master password set successfully.")
	},
}

func init() {
	setMasterCmd.Flags().StringP("password", "p", "", "New master password (required)")
	setMasterCmd.Flags().StringP("old-password", "o", "", "Old master password (required if changing)")
	setMasterCmd.MarkFlagRequired("password") // Make new password flag required
}
