package v1

import (
    "github.com/deadnotxaa/education-platform/backend/internal/usecase"
    "github.com/deadnotxaa/education-platform/backend/pkg/logger"
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

// NewTranslationRoutes -.
func NewTranslationRoutes(apiV1Group fiber.Router, p usecase.Platform, l logger.Interface) {
    r := &V1{p: p, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

    translationGroup := apiV1Group.Group("/course")

    {
        translationGroup.Get("/getcourse", r.getCourse)
    }
}