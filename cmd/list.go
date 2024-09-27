package cmd

import (
	"fmt"
	"strings"
	db "vault-cli/database"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stored services and identifiers",
	Long:  `List all services stored in the vault and their associated identifiers. You can filter by identifier type (e.g., username, email, api_key).`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the id_type flag from the command
		idType, _ := cmd.Flags().GetString("id-type")

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

		// Fetch all sensitive data from the database, potentially filtering by id_type
		entries, err := db.GetAllSensitiveData(idType)
		if err != nil {
			fmt.Println("Error fetching sensitive data:", err)
			return
		}

		if len(entries) == 0 {
			fmt.Println("No data found in the vault.")
			return
		}
		fmt.Println("Stored Services and Identifiers:")
		fmt.Println("--------------------------------")
	
		serviceMap := make(map[string][]db.SensitiveData)
	
		// Group entries by service
		for _, entry := range entries {
			serviceMap[entry.Service] = append(serviceMap[entry.Service], entry)
		}
	
		// Header with color
		fmt.Printf("\033[1;37m%-20s | %-10s | %-30s\033[0m\n", "Service", "Type", "Identifier")
		fmt.Println(strings.Repeat("-", 65))
	
		// Display grouped services with alternating colors
		alternate := false
		for service, entries := range serviceMap {
			for _, entry := range entries {
				// Switch row colors: Light Gray for odd rows, Normal for even rows
				if alternate {
					fmt.Printf("\033[0;37m%-20s | %-10s | %-30s\033[0m\n", service, entry.IdentifierType, entry.Identifier)
				} else {
					fmt.Printf("%-20s | %-10s | %-30s\n", service, entry.IdentifierType, entry.Identifier)
				}
				alternate = !alternate
			}
		}
	},
}

func init() {
	// Add the id-type flag to filter by identifier type
	listCmd.Flags().StringP("id-type", "t", "", "Filter by identifier type (e.g., username, email, api_key)")
}
