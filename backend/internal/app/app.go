package app

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/deadnotxaa/education-platform/backend/config"
    "github.com/deadnotxaa/education-platform/backend/internal/controller/http"
    "github.com/deadnotxaa/education-platform/backend/internal/repo/cache"
    "github.com/deadnotxaa/education-platform/backend/internal/repo/persistent"
    "github.com/deadnotxaa/education-platform/backend/internal/usecase/platform"
    "github.com/deadnotxaa/education-platform/backend/pkg/httpserver"
    "github.com/deadnotxaa/education-platform/backend/pkg/logger"
    "github.com/deadnotxaa/education-platform/backend/pkg/postgres"
    "github.com/deadnotxaa/education-platform/backend/pkg/redis"
)

func Run(cfg *config.Config) {
    // Logger
    l := logger.New(cfg.Log.Level)

    // Postgres Repository
    postgresConnStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        cfg.Postgres.PostgresUser,
        cfg.Postgres.PostgresPassword,
        cfg.Postgres.PostgresHost,
        cfg.Postgres.PostgresPort,
        cfg.Postgres.PostgresDbName,
    )

    pg, err := postgres.New(postgresConnStr, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
    if err != nil {
        l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
    }
    defer pg.Close()

    // Redis Repository
    redisConnStr := fmt.Sprintf("redis://%s:%s@%s:%s/%s",
        cfg.Redis.RedisUser,
        cfg.Redis.RedisPassword,
        cfg.Redis.RedisHost,
        cfg.Redis.RedisPort,
        cfg.Redis.RedisDbName,
    )

    rdb, err := redis.New(redisConnStr)
    if err != nil {
        l.Fatal(fmt.Errorf("app - Run - redis.New: %w", err))
    }

    defer func(rdb *redis.Redis) {
        err = rdb.Close()
        if err != nil {
            l.Fatal(fmt.Errorf("app - Run - rdb.Close: %w", err))
        }
    }(rdb)

    // Use-Case
    rdbRepo := cache.New(rdb)
    pgRepo := persistent.New(pg, rdbRepo)

    platformUseCase := platform.New(
        pgRepo,
        rdbRepo,
    )

    // HTTP Server
    httpServer := httpserver.New(httpserver.Port(cfg.HTTP.Port), httpserver.Prefork(cfg.HTTP.UsePreforkMode))
    http.NewRouter(httpServer.App, cfg, platformUseCase, l)

    // Start servers
    httpServer.Start()

    // Waiting signal
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

    select {
    case s := <-interrupt:
        l.Info("app - Run - signal: %s", s.String())
    case err = <-httpServer.Notify():
        l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
    }

    // Shutdown
    err = httpServer.Shutdown()
    if err != nil {
        l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
    }
}
