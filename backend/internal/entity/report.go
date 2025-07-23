// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // TopCoursesReport - represents a report of top courses with their details.
    TopCoursesReport struct {
        CourseName         string  `json:"name"                 example:"Introduction to Go"`
        DifficultyLevel    string  `json:"difficulty_level"     example:"Advanced"`
        Duration           int     `json:"duration"             example:"30"` // Duration in days
        AverageRating      float64 `json:"average_rating"       example:"4.5"`
        TotalReviews       int     `json:"total_reviews"        example:"150"`
        TeachersWorkPlaces string  `json:"teachers_work_places" example:"Ozon, Tindex, U-Bank"`
    }

    // DetailedPurchaseReport - represents a detailed report of a user's course purchase.
    DetailedPurchaseReport struct {
        UserName           string  `json:"user_name"            example:"John"`
        UserSurname        string  `json:"user_surname"         example:"Doe"`
        CourseName         string  `json:"course_name"          example:"Introduction to Go"`
        SpecializationName string  `json:"specialization_name"  example:"GoLang Development"`
        CourseType         string  `json:"course_type"          example:"Disabled Person discount"`
        TotalPrice         float64 `json:"total_price"          example:"179.99"`
        PurchaseDate       string  `json:"purchase_date"        example:"2022-01-02"`
        TeacherWorkPlace   string  `json:"teacher_work_place"   example:"Ozon"`
    }
)
