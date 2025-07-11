// Package http v1 implements routing paths. Each service in its own file.
package http

import (
    v1 "github.com/deadnotxaa/education-platform/backend/internal/controller/http/v1"
    "net/http"

    "github.com/deadnotxaa/education-platform/backend/config"
    "github.com/deadnotxaa/education-platform/backend/internal/controller/http/middleware"
    "github.com/deadnotxaa/education-platform/backend/internal/usecase"
    "github.com/deadnotxaa/education-platform/backend/pkg/logger"

    "github.com/ansrivas/fiberprometheus/v2"
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/swagger"
)

// NewRouter -.
// Swagger spec:
// @title       Educational Platform API
// @description Description of the Educational Platform API
// @version     1.0
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(app *fiber.App, cfg *config.Config, t usecase.Platform, l logger.Interface) {
    // Options
    app.Use(middleware.Logger(l))
    app.Use(middleware.Recovery(l))

    // Prometheus metrics
    if cfg.Metrics.Enabled {
        prometheus := fiberprometheus.New("my-service-name")
        prometheus.RegisterAt(app, "/metrics")
        app.Use(prometheus.Middleware)
    }

    // Swagger
    if cfg.Swagger.Enabled {
        app.Get("/swagger/*", swagger.HandlerDefault)
    }

    // K8s probe
    app.Get("/healthz", func(ctx *fiber.Ctx) error { return ctx.SendStatus(http.StatusOK) })

    // Routers
    apiV1Group := app.Group("/v1")
    {
        v1.NewCourseRoutes(apiV1Group, t, l)
        v1.NewUserRoutes(apiV1Group, t, l)
        v1.NewReportRoutes(apiV1Group, t, l)
    }
}
