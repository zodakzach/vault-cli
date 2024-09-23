package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
    db "data-manager/database"
)

// setMasterCmd represents the set-master command
var setMasterCmd = &cobra.Command{
    Use:   "set-master",
    Short: "Set or update the master password",
    Long:  `Set or update the master password for accessing the password vault.`,
    Run: func(cmd *cobra.Command, args []string) {
        masterPassword, _ := cmd.Flags().GetString("password")

        // Handle the logic to set the master password
        err := db.SetMasterPassword(masterPassword)
        if err != nil {
            fmt.Println("Error setting master password:", err)
            return
        }

        fmt.Println("Master password set successfully.")
    },
}

func init() {
    setMasterCmd.Flags().StringP("password", "p", "", "Master password (required)")
    setMasterCmd.MarkFlagRequired("password") // Make password flag required
}