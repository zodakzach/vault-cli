package main

import (
    "log"
    "data-manager/cmd"
    "data-manager/database"
)

func main() {
    // Initialize the database
    if err := database.InitDB(); err != nil {
        log.Fatalf("Could not initialize the database: %v", err)
    }

    // Execute the commands
    cmd.Execute()
}
