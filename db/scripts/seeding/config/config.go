package config

import (
    "fmt"
    "github.com/caarlos0/env/v11"
    _ "github.com/caarlos0/env/v11"
)

type (
    Config struct {
        DBConfig   DBConfig
        SeedConfig SeedConfig
    }

    DBConfig struct {
        DbHost     string `env:"DB_HOST,required"`
        DbPort     string `env:"DB_PORT,required"`
        DbUser     string `env:"DB_USER,required"`
        DbPassword string `env:"DB_PASSWORD,required"`
        DbName     string `env:"DB_NAME,required"`
    }

    SeedConfig struct {
        SeedCount        int    `env:"SEED_COUNT"        envDefault:"10"`
        MigrationVersion string `env:"MIGRATION_VERSION" envDefault:"latest"`
        InsertBatchSize  int    `env:"INSERT_BATCH_SIZE" envDefault:"1000"`
    }
)

func NewConfig() (*Config, error) {
    cfg := &Config{}

    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("failed to parse environmental variables: %w", err)
    }

    return cfg, nil
}
