package cmd

import (
	"bufio"
	db "data-manager/database"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new sensitive data entry to the vault",
	Long:  `Add a new sensitive data entry to the vault with the specified service, identifier, and value.`,
	Run: func(cmd *cobra.Command, args []string) {
		service, _ := cmd.Flags().GetString("service")

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

		// Prompt for identifier type
		idType := promptIdentifierType()

		// Prompt for identifier based on selected identifier type
		identifier := promptForIdentifier(idType)

		value := promptPassword("Enter data value for the service (or press Enter to auto-generate): ")

		// Automatically generate a random password if not provided
		if value == "" {
			// Auto-generate a password if none provided
			fmt.Println("No password entered. Generating a random password...")
			value, err = generateRandomPassword(12) // Adjust the length as needed
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

	addCmd.MarkFlagRequired("service")
}

// promptIdentifierType prompts the user to select an identifier type using arrow keys
func promptIdentifierType() string {
	// Available identifier types
	idTypes := []string{"username", "email", "api_key", "secret_key"}

	// Prompt UI for selecting the identifier type
	prompt := promptui.Select{
		Label: "Select an identifier type",
		Items: idTypes,
	}

	// Run the prompt and get the selected index
	_, selectedType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return ""
	}

	return selectedType
}

// promptForIdentifier prompts the user for an identifier value based on the identifier type
func promptForIdentifier(idType string) string {
	var prompt string
	switch idType {
	case "username":
		prompt = "Enter the username: "
	case "email":
		prompt = "Enter the email address: "
	case "api_key":
		prompt = "Enter the API key: "
	case "secret_key":
		prompt = "Enter the secret key: "
	default:
		prompt = "Enter the identifier: "
	}

	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	identifier, _ := reader.ReadString('\n')

	return strings.TrimSpace(identifier)
}

// promptPassword prompts the user for a password and hides input
func promptPassword(prompt string) string {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Move to the next line after password input
	if err != nil {
		fmt.Println("Error reading password:", err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(bytePassword))
}
