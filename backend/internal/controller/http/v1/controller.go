package v1

import (
    "github.com/deadnotxaa/education-platform/backend/internal/usecase"
    "github.com/deadnotxaa/education-platform/backend/pkg/logger"

    "github.com/go-playground/validator/v10"
)

type V1 struct {
    p usecase.Platform
    l logger.Interface
    v *validator.Validate
}
