package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	db "data-manager/database"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import password entries from a file (CSV or JSON)",
	Long:  `Import password entries from a specified file in either CSV or JSON format.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the filename from the flags
		fileName, _ := cmd.Flags().GetString("file")

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

		// Validate file extension (only .json or .csv are accepted)
		ext := strings.ToLower(filepath.Ext(fileName))
		if ext != ".json" && ext != ".csv" {
			fmt.Println("Error: Invalid file type. Only .json or .csv files are accepted.")
			return
		}

		// Import based on file type
		if ext == ".json" {
			err = importFromJSON(fileName)
		} else if ext == ".csv" {
			err = importFromCSV(fileName)
		}

		if err != nil {
			fmt.Printf("Error importing data: %v\n", err)
			return
		}

		fmt.Println("Data imported successfully.")
	},
}

// init initializes the import command flags
func init() {
	importCmd.Flags().StringP("file", "f", "", "Filename to import data from (required)")
	importCmd.MarkFlagRequired("file")
}

// importFromJSON reads and parses a JSON file, then adds the entries to the vault
func importFromJSON(fileName string) error {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Decode JSON data
	var entries []db.SensitiveData
	err = json.NewDecoder(file).Decode(&entries)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	// Add each entry to the vault
	for _, entry := range entries {
		err = db.AddSensitiveData(entry.Service, entry.Identifier, entry.Value, string(entry.IdentifierType))
		if err != nil {
			return fmt.Errorf("failed to add entry for service %s: %v", entry.Service, err)
		}
	}

	return nil
}

// importFromCSV reads and parses a CSV file, then adds the entries to the vault
func importFromCSV(fileName string) error {
	// Open the file
	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read the CSV headers (should match the format used in export)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV headers: %v", err)
	}

	// Ensure headers match the expected format
	expectedHeaders := []string{"Service", "Identifier", "Identifier Type", "Value"}
	if !compareHeaders(headers, expectedHeaders) {
		return fmt.Errorf("CSV headers do not match the expected format")
	}

	// Read each row and add the entry to the vault
	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read CSV row: %v", err)
		}

		// Convert the identifier type (row[2]) from string to db.IdentifierType
		identifierType, err := db.ParseIdentifierType(row[2])
		if err != nil {
			return fmt.Errorf("invalid identifier type: %v", err)
		}

		// Create a new entry and add to the vault
		entry := db.SensitiveData{
			Service:        row[0],
			Identifier:     row[1],
			IdentifierType: identifierType,
			Value:          row[3],
		}

		err = db.AddSensitiveData(entry.Service, entry.Identifier, entry.Value, string(entry.IdentifierType))
		if err != nil {
			return fmt.Errorf("failed to add entry for service %s: %v", entry.Service, err)
		}
	}

	return nil
}

// compareHeaders checks if the CSV headers match the expected format
func compareHeaders(headers, expected []string) bool {
	if len(headers) != len(expected) {
		return false
	}

	for i, header := range headers {
		if strings.TrimSpace(header) != strings.TrimSpace(expected[i]) {
			return false
		}
	}

	return true
}
