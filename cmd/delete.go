package cmd

import (
	db "data-manager/database"
	"fmt"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a stored entry from the vault",
	Long:  `Delete a stored entry from the vault using the specified service and identifier.`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")
		identifier, _ := cmd.Flags().GetString("identifier")

		// Check if the vault is locked
		isLocked, err := db.GetVaultState()
		if err != nil {
			fmt.Println("Error getting vault state:", err)
			return
		}
		if isLocked {
			fmt.Println("Vault is locked. Please unlock it first.")
			return
		}

		if service == "" || identifier == "" {
			fmt.Println("Error: Both service and identifier are required.")
			return
		}

		// Attempt to delete the entry
		err = db.DeleteSensitiveData(service, identifier)
		if err != nil {
			fmt.Println("Error deleting entry:", err)
			return
		}

		fmt.Println("Entry deleted successfully.")
	},
}

func init() {
	// Add flags for service and identifier
	deleteCmd.Flags().StringP("service", "s", "", "Service name (required)")
	deleteCmd.Flags().StringP("identifier", "i", "", "Identifier (required)")
	deleteCmd.MarkFlagRequired("service")
	deleteCmd.MarkFlagRequired("identifier")

}
