package persistent

import (
	"context"
	"fmt"
	"log"

	"github.com/deadnotxaa/education-platform/backend/internal/entity"
	"github.com/deadnotxaa/education-platform/backend/internal/repo"
	"github.com/deadnotxaa/education-platform/backend/pkg/postgres"
)

// PostgresRepo -.
type PostgresRepo struct {
    *postgres.Postgres
    rr repo.RedisRepo
}

// New -.
func New(pg *postgres.Postgres, rr repo.RedisRepo) *PostgresRepo {
    return &PostgresRepo{pg, rr}
}

// GetCourseById -.
func (r *PostgresRepo) GetCourseById(ctx context.Context, courseID int) (entity.Course, error) {
    sql, args, err := r.Builder.
        Select("course_id", "name", "description", "specialization_id", "duration", "price", "difficulty_level_id",
            "created_at", "updated_at").
        From("course").
        Where("course_id = ?", courseID).
        ToSql()

    if err != nil {
        return entity.Course{}, fmt.Errorf("PostgresRepo - GetCourse - r.Builder: %w", err)
    }

    row := r.Pool.QueryRow(ctx, sql, args...)

    ent := entity.Course{}
    var createdAt, updatedAt interface{}

    err = row.Scan(&ent.CourseID, &ent.Name, &ent.Description, &ent.SpecializationID, &ent.Duration, &ent.Price,
        &ent.DifficultyLevelID, &createdAt, &updatedAt)

    if err != nil {
        return entity.Course{}, fmt.Errorf("PostgresRepo - GetCourse - row.Scan: %w", err)
    }

    // Handle nullable timestamps TODO: change after redesign of the Postgres schema
    if createdAt != nil {
        ent.CreatedAt = createdAt.(string)
    }
    if updatedAt != nil {
        ent.UpdatedAt = updatedAt.(string)
    }

    return ent, nil
}

func (r *PostgresRepo) GetUserById(ctx context.Context, userID int) (entity.User, error) {
    sql, args, err := r.Builder.
        Select("account_id", "name", "surname", "email").
        From("users").
        Where("account_id = ?", userID).
        ToSql()

    if err != nil {
        return entity.User{}, fmt.Errorf("PostgresRepo - GetUserById - r.Builder: %w", err)
    }

    row := r.Pool.QueryRow(ctx, sql, args...)

    ent := entity.User{}
    err = row.Scan(&ent.AccountID, &ent.Name, &ent.Surname, &ent.Email)

    if err != nil {
        return entity.User{}, fmt.Errorf("PostgresRepo - GetUserById - row.Scan: %w", err)
    }

    return ent, nil
}

func (r *PostgresRepo) GetTopCoursesReport(ctx context.Context, limit uint32) ([]entity.TopCoursesReport, error) {
    // Try to get data from Redis first
    report, err := r.rr.GetTopCoursesReport(ctx, limit)
    if err == nil {
       log.Printf("reports found in Redis cache for limit %d", limit)
       return report, nil
    }

    rows, err := r.Pool.Query(ctx,
        `SELECT 
            c.name AS course_name,
            dl.name AS difficulty_level,
            c.duration,
            AVG(cr.rating) AS avg_rating,
            COUNT(cr.review_id)::int AS total_reviews,
            COALESCE(STRING_AGG(DISTINCT t.work_place, ', '), '') AS teachers_work_places
        FROM course c
        JOIN difficulty_level dl ON c.difficulty_level_id = dl.id
        LEFT JOIN course_review cr ON c.course_id = cr.course_id
        LEFT JOIN course_teacher ct ON c.course_id = ct.course_id
        LEFT JOIN teacher t ON ct.teacher_id = t.employee_id
        GROUP BY c.course_id, c.name, dl.name, c.duration
        ORDER BY avg_rating DESC NULLS LAST
        LIMIT $1;`,
        limit,
    )

    if err != nil {
        return nil, fmt.Errorf("PostgresRepo - GetTopCoursesReport - r.Pool.Query: %w", err)
    }
    defer rows.Close()

    entities := make([]entity.TopCoursesReport, 0, limit)

    for rows.Next() {
        e := entity.TopCoursesReport{}

        err = rows.Scan(&e.CourseName, &e.DifficultyLevel, &e.Duration, &e.AverageRating,
            &e.TotalReviews, &e.TeachersWorkPlaces)

        if err != nil {
            return nil, fmt.Errorf("PostgresRepo - GetTopCoursesReport - rows.Scan: %w", err)
        }

        entities = append(entities, e)
    }

    err = r.rr.SetTopCoursesReport(ctx, limit, entities)
    if err != nil {
        return nil, fmt.Errorf("PostgresRepo - GetTopCoursesReport - r.rr.SetTopCoursesReport: %w", err)
    }

    return entities, nil
}
