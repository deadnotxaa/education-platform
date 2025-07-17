package platform

import (
    "context"
    "fmt"
    "github.com/deadnotxaa/education-platform/backend/internal/entity"
    "github.com/deadnotxaa/education-platform/backend/internal/repo"
)

// UseCase - Platform use case
type UseCase struct {
    postgresRepo repo.PostgresRepo
    redisRepo    repo.RedisRepo
}

// New -.
func New(pgr repo.PostgresRepo, rr repo.RedisRepo) *UseCase {
    return &UseCase{
        postgresRepo: pgr,
        redisRepo:    rr,
    }
}

func (us *UseCase) GetCourseById(ctx context.Context, courseID int) (entity.Course, error) {
    course, err := us.postgresRepo.GetCourseById(ctx, courseID)
    if err != nil {
        return entity.Course{}, fmt.Errorf("platform - GetCourse - postgresRepo.GetCourse: %w", err)
    }

    return course, nil
}

func (us *UseCase) GetUserById(ctx context.Context, userID int) (entity.User, error) {
    user, err := us.postgresRepo.GetUserById(ctx, userID)
    if err != nil {
        return entity.User{}, fmt.Errorf("platform - GetUser - postgresRepo.GetUser: %w", err)
    }

    return user, nil
}

func (us *UseCase) GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error) {
    reports, err := us.postgresRepo.GetTopCoursesReport(ctx, limit)
    if err != nil {
        return nil, fmt.Errorf("platform - GetTopCoursesReport - postgresRepo.GetTopCoursesReport: %w", err)
    }

    return reports, nil
}
