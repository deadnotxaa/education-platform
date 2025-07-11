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

func (us *UseCase) GetUserById(ctx context.Context, userID int) (entity.User, error) {
    user, err := us.repo.GetUserById(ctx, userID)
    if err != nil {
        return entity.User{}, fmt.Errorf("platform - GetUser - repo.GetUser: %w", err)
    }

    return user, nil
}

func (us *UseCase) GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error) {
    report, err := us.repo.GetTopCoursesReport(ctx, limit)
    if err != nil {
        return nil, fmt.Errorf("platform - GetTopCoursesReport - repo.GetTopCoursesReport: %w", err)
    }

    return report, nil
}
