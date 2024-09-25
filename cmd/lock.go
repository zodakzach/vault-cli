package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    "data-manager/vault" 
)

var lockCmd = &cobra.Command{
    Use:   "lock",
    Short: "Lock the vault",
    Long:  `Lock the vault, preventing access to sensitive data until it is unlocked again.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Lock the vault 
		err := vault.LockVault()
        if err != nil {
            fmt.Println("Error locking the vault:", err)
            return
        }
    },
}

func init() {
    // For now, locking the vault does not require any flags
    // lockCmd.Flags().StringP("some-flag", "s", "", "Some description")
    // lockCmd.MarkFlagRequired("some-flag") // Uncomment if needed
}

