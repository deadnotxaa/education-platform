// Package usecase implements application business logic. Each logic group in own its file.
package usecase

import (
    "context"

    "github.com/deadnotxaa/education-platform/backend/internal/entity"
)

type (
    // Platform -.
    Platform interface {
        GetCourseById(ctx context.Context, courseID int) (entity.Course, error)
    }
)
