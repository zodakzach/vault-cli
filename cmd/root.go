package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "vault", // The name of your command
    Short: "A secure password manager", // Short description
    Long:  `Vault is a secure password manager for storing and retrieving your passwords from the terminal.`,
    Run: func(cmd *cobra.Command, args []string) {
        // Default action when no subcommands are provided
        cmd.Help() // Show help if no subcommand is given
    },
}


// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        // Exit the program or handle error appropriately
    }
}

func init() {
    // Here you can define flags and configuration settings.
    // For example, adding a global flag:
    // rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.vault.yaml)")

    // Add subcommands to rootCmd
    rootCmd.AddCommand(setMasterCmd)
    rootCmd.AddCommand(unlockCmd)
}
