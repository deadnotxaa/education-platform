package platform

import (
    "context"
    "fmt"
    "github.com/deadnotxaa/education-platform/backend/internal/entity"
    "github.com/deadnotxaa/education-platform/backend/internal/repo"
)

// UseCase - Platform use case
type UseCase struct {
    repo repo.BackendRepo
}

// New -.
func New(r repo.BackendRepo) *UseCase {
    return &UseCase{
        repo: r,
    }
}

func (us *UseCase) GetCourseById(ctx context.Context, courseID int) (entity.Course, error) {
    course, err := us.repo.GetCourseById(ctx, courseID)
    if err != nil {
        return entity.Course{}, fmt.Errorf("platform - GetCourse - repo.GetCourse: %w", err)
    }

    return course, nil
}
