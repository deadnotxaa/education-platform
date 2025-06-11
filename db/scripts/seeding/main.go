package main

import (
    "log"
    "seeder/config"
    "seeder/internal/seed"
)

func main() {
    // Configuration
    cfg, err := config.NewConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Run seeding
    seed.RunSeeding(cfg)
}
