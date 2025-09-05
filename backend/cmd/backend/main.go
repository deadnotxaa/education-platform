package main

import (
    "log"

    "github.com/deadnotxaa/education-platform/backend/config"
    "github.com/deadnotxaa/education-platform/backend/internal/app"
)

func main() {
    // Configure the application
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Error loading configuration: %v", err)
    }

    // Initialize the application with the configuration
    app.Run(cfg)
}
