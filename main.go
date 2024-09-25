package main

import (
    "log"
    "os"
    "data-manager/cmd"
    "data-manager/database"
)

func main() {
    // Get the database path from the environment variable
    dbPath := os.Getenv("SQLITE_DB")
    if dbPath == "" {
        log.Fatal("SQLITE_DB environment variable is not set")
    }

    // Initialize the database
    if err := database.InitDB(dbPath); err != nil {
        log.Fatalf("Could not initialize the database: %v", err)
    }

    // Execute the commands
    cmd.Execute()
}
