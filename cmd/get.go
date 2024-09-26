package cmd

import (
	db "vault-cli/database"
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve a sensitive data entry in the vault",
	Long:  `Retrieve a specific service and identifier from the vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
		identifier, _ := cmd.Flags().GetString("identifier")

		// Check if the vault is locked
		isLocked, err := db.GetVaultState()
		if err != nil {
			fmt.Println("Error retrieving vault state:", err)
			return
		}

		if isLocked {
			fmt.Println("Error: Vault is locked. Please unlock the vault using `unlock`.")
			return
		}

		if service == "" || identifier == "" {
			fmt.Println("Error: Both service and identifier are required.")
			return
		}

		// Retrieve the sensitive data based on service and identifier
		entry, err := db.GetSensitiveData(service, identifier)
		if err != nil {
			fmt.Println("Error retrieving data:", err)
			return
		}

		// Print the retrieved value
		fmt.Printf("Service: %s\n", entry.Service)
		fmt.Printf("%s: %s\n", cases.Title(language.Und).String(string(entry.IdentifierType)), entry.Identifier)
		fmt.Printf("Password: %s\n", entry.Value)
	},
}

func init() {
	getCmd.Flags().StringP("service", "s", "", "Service name (required)")
	getCmd.Flags().StringP("identifier", "i", "", "Identifier (required)")
	getCmd.MarkFlagRequired("service")
	getCmd.MarkFlagRequired("identifier")
}
