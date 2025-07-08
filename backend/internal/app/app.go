package app

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"

    "github.com/deadnotxaa/education-platform/backend/config"
    "github.com/deadnotxaa/education-platform/backend/internal/controller/http"
    "github.com/deadnotxaa/education-platform/backend/internal/repo/persistent"
    "github.com/deadnotxaa/education-platform/backend/internal/usecase/platform"
    "github.com/deadnotxaa/education-platform/backend/pkg/httpserver"
    "github.com/deadnotxaa/education-platform/backend/pkg/logger"
    "github.com/deadnotxaa/education-platform/backend/pkg/postgres"
)

func Run(cfg *config.Config) {
    // Logger
    l := logger.New(cfg.Log.Level)

    // Repository
    connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        cfg.DB.PostgresUser,
        cfg.DB.PostgresPassword,
        cfg.DB.PostgresHost,
        cfg.DB.PostgresPort,
        cfg.DB.PostgresDbName,
    )

    pg, err := postgres.New(connStr, postgres.MaxPoolSize(cfg.DB.PoolMax))
    if err != nil {
        l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
    }
    defer pg.Close()

    // Use-Case
    platformUseCase := platform.New(
        persistent.New(pg),
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
