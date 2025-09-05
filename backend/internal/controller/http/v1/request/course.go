package request

type Course struct {
    ID               int    `json:"id" validate:"required" example:"1"`
}
