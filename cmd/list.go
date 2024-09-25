package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	db "data-manager/database"
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

		// Group and display the entries by service
		fmt.Println("Stored services and their identifiers:")
		serviceMap := make(map[string][]db.SensitiveData)

		// Group entries by service
		for _, entry := range entries {
			serviceMap[entry.Service] = append(serviceMap[entry.Service], entry)
		}

		// Display the grouped services and identifiers
		for service, entries := range serviceMap {
			fmt.Printf("Service: %s\n", service)
			for _, entry := range entries {
				fmt.Printf("  %s: %s\n", entry.IdentifierType, entry.Identifier) // Print id_type and identifier
			}
		}
	},
}

func init() {
	// Add the id-type flag to filter by identifier type
	listCmd.Flags().StringP("id-type", "t", "", "Filter by identifier type (e.g., username, email, api_key)")
}
