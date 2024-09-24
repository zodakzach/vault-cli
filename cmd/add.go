package cmd

import (
	db "data-manager/database"
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new sensitive data entry to the vault",
	Long:  `Add a new sensitive data entry to the vault with the specified service, identifier, and value.`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
		identifier, _ := cmd.Flags().GetString("identifier")
		value, _ := cmd.Flags().GetString("value")
		idType, _ := cmd.Flags().GetString("id-type")

        // Check if the vault is locked
        isLocked, err := db.GetVaultState()
        if err != nil {
            fmt.Println("Error retrieving vault state:", err)
            return
        }

        if isLocked {
            fmt.Println("Error: Vault is locked. Please unlock the vault using `unlock` before adding new entries.")
            return
        }

		if service == "" || identifier == "" {
			fmt.Println("Error: Service and identifier are required.")
			return
		}

		if err = db.CheckMasterPasswordSet(); err != nil {
			fmt.Println(err)
			return
		}

		// Automatically generate a random password if not provided
		if value == "" {
			value, err = GenerateRandomPassword(12) // Adjust the length as needed
			if err != nil {
				fmt.Println("Error generating password:", err)
				return
			}
			fmt.Printf("Generated password for %s: %s\n", service, value)
		}

		// Add the sensitive data to the vault
		err = db.AddSensitiveData(service, identifier, value, idType) // Assuming idType is username
		if err != nil {
			fmt.Println("Error adding Sensitive data entry:", err)
			return
		}

		fmt.Println("Sensitive data entry added successfully.")
	},
}

func init() {
	addCmd.Flags().StringP("service", "s", "", "Service name (required)")
	addCmd.Flags().StringP("identifier", "i", "", "Identifier (required)")
	addCmd.Flags().StringP("value", "v", "", "Value (optional)")
	addCmd.Flags().StringP("id-type", "t", "username", "Identifier type (e.g., username, email, api_key, secret_key)")

	addCmd.MarkFlagRequired("service")
	addCmd.MarkFlagRequired("identifier")
	addCmd.MarkFlagRequired("id-type") // Ensure id-type is also required
}
