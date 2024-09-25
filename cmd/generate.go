package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
    Use:   "generate",
    Short: "Generate a random password",
    Long:  `Generate a secure random password of the specified length.`,
    Run: func(cmd *cobra.Command, args []string) {
        length, _ := cmd.Flags().GetInt("length")
        if length <= 0 {
            fmt.Println("Error: Length must be a positive integer.")
            return
        }

        password, err := GenerateRandomPassword(length)
        if err != nil {
            fmt.Println("Error generating password:", err)
            return
        }

        fmt.Println("Generated Password:", password)
    },
}

func init() {
    generateCmd.Flags().IntP("length", "l", 12, "Length of the password (default is 12)")
}
