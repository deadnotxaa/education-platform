package request

type User struct {
    ID               int    `json:"id" validate:"required" example:"1"`
}
