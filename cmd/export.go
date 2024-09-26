package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	db "vault-cli/database"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export all sensitive data entries to a file (CSV or JSON)",
	Long:  `Export all stored sensitive data entries to a specified file in either CSV or JSON format.`,
	Run: func(cmd *cobra.Command, args []string) {
		fileName, _ := cmd.Flags().GetString("file")
		format, _ := cmd.Flags().GetString("format")

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

		// Automatically append the appropriate file extension
		filePath := appendFileExtension(fileName, format)

		// Retrieve all sensitive data from the vault
		entries, err := db.GetAllSensitiveData("")
		if err != nil {
			fmt.Printf("Error retrieving sensitive data: %v\n", err)
			return
		}

		// Export based on format
		switch format {
		case "json":
			err = exportToJSON(filePath, entries)
		case "csv":
			err = exportToCSV(filePath, entries)
		}

		if err != nil {
			fmt.Printf("Error exporting data: %v\n", err)
			return
		}

		fmt.Println("Sensitive data exported successfully.")
	},
}

func init() {
	exportCmd.Flags().StringP("file", "f", "", "File path to export data (required)")
	exportCmd.Flags().StringP("format", "t", "json", "Export format (json or csv)")
	exportCmd.MarkFlagRequired("file")
	exportCmd.MarkFlagRequired("format")

}

// appendFileExtension appends the correct file extension based on the format
func appendFileExtension(fileName, format string) string {
	// Get the file extension from the fileName if it exists
	ext := filepath.Ext(fileName)

	// If the file already has the correct extension, return the file name as is
	if (format == "json" && ext == ".json") || (format == "csv" && ext == ".csv") {
		return fileName
	}

	// Strip the existing extension if it exists, and append the appropriate one
	baseName := strings.TrimSuffix(fileName, ext)
	if format == "json" {
		return baseName + ".json"
	} else {
		return baseName + ".csv"
	}
}

// exportToJSON exports sensitive data to a JSON file
func exportToJSON(filePath string, entries []db.SensitiveData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty-print JSON
	err = encoder.Encode(entries)
	if err != nil {
		return fmt.Errorf("failed to encode data to JSON: %v", err)
	}

	return nil
}

// exportToCSV exports sensitive data to a CSV file
func exportToCSV(filePath string, entries []db.SensitiveData) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV headers
	headers := []string{"Service", "Identifier", "Identifier Type", "Value"}
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("failed to write CSV headers: %v", err)
	}

	// Write CSV rows for each entry
	for _, entry := range entries {
		row := []string{entry.Service, entry.Identifier, string(entry.IdentifierType), entry.Value}
		err = writer.Write(row)
		if err != nil {
			return fmt.Errorf("failed to write CSV row: %v", err)
		}
	}

	return nil
}
