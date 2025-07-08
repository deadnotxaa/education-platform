// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // DifficultyLevel -.
    DifficultyLevel struct {
        ID          int    `json:"id"           example:"3"`
        Name        string `json:"name"         example:"Advanced"`
        Description string `json:"description"  example:"Advanced course level description"`
    }

    // CourseSpecialization -.
    CourseSpecialization struct {
        ID          int    `json:"id"           example:"1"`
        Name        string `json:"name"         example:"GoLang Development"`
        Description string `json:"description"  example:"Software development using Go programming language"`
    }

    // Course -.
    Course struct {
        CourseID          int    `json:"id"                  example:"1"`
        Name              string `json:"name"                example:"Introduction to Go"`
        Description       string `json:"description"         example:"A beginner's course on Go programming language"`
        SpecializationID  int    `json:"specialization_id"   example:"1"`
        Duration          int    `json:"duration"            example:"30"` // Duration in hours
        Price             int    `json:"price"               example:"199.99"`
        DifficultyLevelID int    `json:"difficulty_level_id" example:"3"`
        CreatedAt         string `json:"created_at"          example:"2023-01-01T00:00:00Z"`
        UpdatedAt         string `json:"updated_at"          example:"2023-01-02T00:00:00Z"`
    }

    // CourseReview -.
    CourseReview struct {
        ReviewID   int    `json:"id"              example:"1"`
        CourseID   int    `json:"course_id"       example:"1"`
        UserID     int    `json:"user_id"         example:"42"`
        Rating     int    `json:"rating"          example:"5"` // Rating from 1 to 5
        Comment    string `json:"comment"         example:"Great course, learned a lot!"`
        ReviewDate string `json:"review_date"     example:"2023-01-15T00:00:00Z"`
    }

    // Certificate -.
    Certificate struct {
        CertificateID int    `json:"id"               example:"1"`
        UserID        int    `json:"user_id"          example:"42"`
        CourseID      int    `json:"course_id"        example:"1"`
        IssueDate     string `json:"issue_date"       example:"2023-01-20T00:00:00Z"`
    }

    // CourseCalendar -.
    CourseCalendar struct {
        ID              int    `json:"id"               example:"1"`
        CourseID        int    `json:"course_id"        example:"1"`
        StartDate       string `json:"start_date"       example:"2023-02-01T10:00:00Z"`
        EndSalesDate    string `json:"end_sales_date"   example:"2023-01-30T23:59:59Z"`
        RemainingPlaces int    `json:"remaining_places" example:"20"`
    }
)
