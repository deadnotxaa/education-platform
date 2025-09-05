// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type PurchaseStatus int

const (
    PurchaseStatusPending   PurchaseStatus = iota // Pending purchase
    PurchaseStatusCompleted                       // Completed purchase
    PurchaseStatusCancelled                       // Canceled purchase
)

type (
    // CourseType - represents a type of course discount for purchase, e.g. "disabled person", "student", etc.
    CourseType struct {
        ID       int    `json:"id"           example:"1"`
        TypeName string `json:"type_name"    example:"Disabled Person discount"` // Type of discount
        Discount int    `json:"discount"     example:"10"`                       // Discount percentage
    }

    // Purchase - represents a purchase of a course by a user.
    Purchase struct {
        PurchaseID     int            `json:"id"                 example:"1"`
        UserID         int            `json:"user_id"            example:"1"`
        CourseID       int            `json:"course_id"          example:"1"`
        PurchaseDate   string         `json:"purchase_date"      example:"2022-01-02"`
        CourseTypeID   int            `json:"course_type_id"     example:"1"`      // ID of the course type (discount)
        TotalPrice     int            `json:"total_price"        example:"179.99"` // Total price after discount
        PurchaseStatus PurchaseStatus `json:"purchase_status"    example:"0"`
    }
)
