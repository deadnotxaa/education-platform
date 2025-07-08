package config

import (
    "fmt"

    "github.com/caarlos0/env/v11"
)

type (
    Config struct {
        App App
        DB  DB
        Log Log
        Swagger Swagger
        Metrics Metrics
        HTTP HTTP
    }

    App struct {
        Name string `env:"APP_NAME,required"`
    }

    DB struct {
        PostgresUser     string `env:"PG_USER,required"`
        PostgresPassword string `env:"PG_PASS,required"`
        PostgresHost     string `env:"PG_HOST,required"`
        PostgresPort     string `env:"PG_PORT,required"`
        PostgresDbName   string `env:"PG_NAME,required"`
        PoolMax          int    `env:"PG_POOL_MAX,required"`
    }

    // Log -.
    Log struct {
        Level string `env:"LOG_LEVEL" envDefault:"error"`
    }

    // Swagger -.
    Swagger struct {
        Enabled bool   `env:"SWAGGER_ENABLED" envDefault:"true"`
    }

    // Metrics -.
    Metrics struct {
        Enabled bool   `env:"METRICS_ENABLED" envDefault:"true"`
    }

    // HTTP -.
    HTTP struct {
        Port           string `env:"HTTP_PORT,required"`
        UsePreforkMode bool   `env:"HTTP_USE_PREFORK_MODE" envDefault:"false"`
    }
)

// NewConfig initializes a new Config instance by parsing environment variables.
func NewConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("config error: %w", err)
    }

    return cfg, nil
}
