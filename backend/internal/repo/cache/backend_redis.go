package cache

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/deadnotxaa/education-platform/backend/internal/entity"
    "github.com/deadnotxaa/education-platform/backend/pkg/redis"
    "time"
)

// RedisRepo -.
type RedisRepo struct {
    *redis.Redis
}

// New -.
func New(rdb *redis.Redis) *RedisRepo {
    return &RedisRepo{
        Redis:  rdb,
    }
}

func (rr *RedisRepo) GetCourseById(ctx context.Context, courseID int) (entity.Course, error) {
    return entity.Course{}, nil
}

func (rr *RedisRepo) GetUserById(ctx context.Context, userID int) (entity.User, error) {
    return entity.User{}, nil
}

func (rr *RedisRepo) GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error) {
    key := fmt.Sprintf("top_courses_report:%d", limit)

    cachedData, err := rr.Client.Get(ctx, key).Result()
    if err == nil {
        var reports []entity.TopCoursesReport
        if err := json.Unmarshal([]byte(cachedData), &reports); err != nil {
            return nil, fmt.Errorf("RedisRepo - GetTopCoursesReport - json.Unmarshal: %w", err)
        }
        return reports, nil
    }

    return nil, fmt.Errorf("RedisRepo - GetTopCoursesReport - cache miss or error: %w", err)
}

func (rr *RedisRepo) SetTopCoursesReport(ctx context.Context, limit uint32, reports []entity.TopCoursesReport) error {
    key := fmt.Sprintf("top_courses_report:%d", limit)

    data, err := json.Marshal(reports)
    if err != nil {
        return fmt.Errorf("RedisRepo - SetTopCoursesReport - json.Marshal: %w", err)
    }

    if err := rr.Client.Set(ctx, key, data, 5 * time.Minute).Err(); err != nil {
        return fmt.Errorf("RedisRepo - SetTopCoursesReport - Client.Set: %w", err)
    }

    return nil
}
