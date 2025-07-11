// Package repo implements application outer layer logic. Each logic group in its own file.
package repo

import (
    "context"

    "github.com/deadnotxaa/education-platform/backend/internal/entity"
)

type (
    // BackendRepo defines the methods for interacting with the backend repository.
    BackendRepo interface {
        // GetCourseById retrieves a course by its ID.
        GetCourseById(ctx context.Context, courseID int) (entity.Course, error)

        // GetUserById retrieves some info about user by their ID.
        GetUserById(ctx context.Context, userID int) (entity.User, error)

        // GetTopCoursesReport retrieves a report of the top n courses.
        GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error)
    }
)
