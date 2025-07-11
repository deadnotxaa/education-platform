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

    // Handle nullable timestamps TODO: change after redesign of the DB schema
    if createdAt != nil {
        ent.CreatedAt = createdAt.(string)
    }
    if updatedAt != nil {
        ent.UpdatedAt = updatedAt.(string)
    }

    return ent, nil
}

func (r *BackendRepo) GetUserById(ctx context.Context, userID int) (entity.User, error) {
    sql, args, err := r.Builder.
        Select("account_id", "name", "surname", "email").
        From("users").
        Where("account_id = ?", userID).
        ToSql()

    if err != nil {
        return entity.User{}, fmt.Errorf("BackendRepo - GetUserById - r.Builder: %w", err)
    }

    row := r.Pool.QueryRow(ctx, sql, args...)

    ent := entity.User{}
    err = row.Scan(&ent.AccountID, &ent.Name, &ent.Surname, &ent.Email)

    if err != nil {
        return entity.User{}, fmt.Errorf("BackendRepo - GetUserById - row.Scan: %w", err)
    }

    return ent, nil
}

func (r *BackendRepo) GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error) {
    rows, err := r.Pool.Query(ctx,
        `SELECT 
            c.name AS course_name,
            dl.name AS difficulty_level,
            c.duration,
            AVG(cr.rating) AS avg_rating,
            COUNT(cr.review_id) AS total_reviews,
            STRING_AGG(DISTINCT t.work_place, ', ') AS teachers_work_places
        FROM course c
        JOIN difficulty_level dl ON c.difficulty_level_id = dl.id
        LEFT JOIN course_review cr ON c.course_id = cr.course_id
        LEFT JOIN course_teacher ct ON c.course_id = ct.course_id
        LEFT JOIN teacher t ON ct.teacher_id = t.employee_id
        GROUP BY c.course_id, dl.name, c.duration
        ORDER BY avg_rating DESC NULLS LAST
        LIMIT $1;`,
        limit,
    )

    if err != nil {
        return nil, fmt.Errorf("BackendRepo - GetTopCoursesReport - r.Pool.Query: %w", err)
    }
    defer rows.Close()

    entities := make([]entity.TopCoursesReport, 0, limit)

    for rows.Next() {
        e := entity.TopCoursesReport{}

        err = rows.Scan(&e.CourseName, &e.DifficultyLevel, &e.Duration, &e.AverageRating,
            e.TotalReviews, &e.TeachersWorkPlaces)

        if err != nil {
            return nil, fmt.Errorf("BackendRepo - GetTopCoursesReport - rows.Scan: %w", err)
        }

        entities = append(entities, e)
    }

    return entities, nil
}
