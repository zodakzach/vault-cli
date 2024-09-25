package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	db "data-manager/database"
)

var updateCmd = &cobra.Command{
    Use:   "update",
    Short: "Update a sensitive data entry in the vault",
    Long:  `Update the value or identifier for a specific service in the vault.`,
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

		// Retrieve the existing sensitive data entry for this service and identifier
		existingEntry, err := db.GetSensitiveData(service, identifier)
		if err != nil {
			fmt.Printf("Error retrieving sensitive data: %v\n", err)
			return
		}

		// Prompt for new identifier (if any)
		newIdentifier := promptForInput(fmt.Sprintf("Enter new %s (leave empty to keep the current one): ", existingEntry.IdentifierType), existingEntry.Identifier)

		// Prompt for new value (if any)
		newValue := promptForInput("Enter new value (leave empty to keep the current one): ", existingEntry.Value)

        err = db.UpdateSensitiveData(service, identifier, newValue, newIdentifier)
        if err != nil {
            fmt.Printf("Error updating sensitive data: %v\n", err)
            return
        }
        fmt.Println("Sensitive data updated successfully.")
    },
}

func init() {
    // Define flags for the update command
    updateCmd.Flags().StringP("service", "s", "", "Service name (required)")
    updateCmd.Flags().StringP("identifier", "i", "", "Identifier (required)")
    
    updateCmd.MarkFlagRequired("service")
    updateCmd.MarkFlagRequired("identifier")
    
}

// promptForInput prompts the user for a new value, or keeps the existing one if the input is empty
func promptForInput(prompt string, existingValue string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Return the existing value if no new input is provided
	if input == "" {
		return existingValue
	}
	return input
}