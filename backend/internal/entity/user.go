// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // User represents a user in the system with all necessary fields.
    User struct {
        AccountID         int    `json:"account_id"             example:"1"`
        Name              string `json:"name"                   example:"John"`
        Surname           string `json:"surname"                example:"Doe"`
        BirthDate         string `json:"birth_date"             example:"2022-01-01"`
        Email             string `json:"email"                  example:"mail@example.com"`
        HashedPassword    string `json:"hashed_password"        example:"$2a$10$EIXo1z..."`
        ProfilePictureUrl string `json:"profile_picture_url"    example:"https://example.com/profile.jpg"`
        PhoneNumber       string `json:"phone_number"           example:"+1234567890"`
        SnilsNumber       string `json:"snils_number"           example:"123-456-789 01"`
        CreatedAt         string `json:"created_at"             example:"2022-01-01"`
        UpdatedAt         string `json:"updated_at"             example:"2022-01-02"`
    }
)
