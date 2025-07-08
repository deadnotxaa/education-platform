package persistent

import (
    "context"
    "fmt"

    "github.com/deadnotxaa/education-platform/backend/internal/entity"
    "github.com/deadnotxaa/education-platform/backend/pkg/postgres"
)

// BackendRepo -.
type BackendRepo struct {
    *postgres.Postgres
}

// New -.
func New(pg *postgres.Postgres) *BackendRepo {
    return &BackendRepo{pg}
}

// GetCourseById -.
func (r *BackendRepo) GetCourseById(ctx context.Context, courseID int) (entity.Course, error) {
    sql, args, err := r.Builder.
        Select("course_id", "name", "description", "specialization_id", "duration", "price", "difficulty_level_id",
            "created_at", "updated_at").
        From("course").
        Where("course_id = ?", courseID).
        ToSql()

    if err != nil {
        return entity.Course{}, fmt.Errorf("BackendRepo - GetCourse - r.Builder: %w", err)
    }

    row := r.Pool.QueryRow(ctx, sql, args...)

    ent := entity.Course{}
    var createdAt, updatedAt interface{}

    err = row.Scan(&ent.CourseID, &ent.Name, &ent.Description, &ent.SpecializationID, &ent.Duration, &ent.Price,
        &ent.DifficultyLevelID, &createdAt, &updatedAt)

    if err != nil {
        return entity.Course{}, fmt.Errorf("BackendRepo - GetCourse - row.Scan: %w", err)
    }

    // Handle nullable timestamps
    if createdAt != nil {
        ent.CreatedAt = createdAt.(string)
    }
    if updatedAt != nil {
        ent.UpdatedAt = updatedAt.(string)
    }

    return ent, nil
}
