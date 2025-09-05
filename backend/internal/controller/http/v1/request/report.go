package request

type (
    TopCoursesReport struct {
        LimitNumber uint32 `json:"limit_number" validate:"required" example:"10"`
    }

    DetailedPurchaseReport struct {
        LimitNumber uint32 `json:"limit_number" validate:"required" example:"10"`
    }
)
