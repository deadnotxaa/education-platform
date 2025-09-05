// Package usecase implements application business logic. Each logic group in own its file.
package usecase

import (
    "context"

    "github.com/deadnotxaa/education-platform/backend/internal/entity"
)

type (
    // Platform - specifies platform's use case interface.
    Platform interface {
        // GetCourseById retrieves a course by its ID.
        GetCourseById(ctx context.Context, courseID int) (entity.Course, error)

        // GetUserById retrieves some info about user by their ID.
        GetUserById(ctx context.Context, userID int) (entity.User, error)

        // GetTopCoursesReport retrieves a report of the top n courses.
        GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error)
    }
)
