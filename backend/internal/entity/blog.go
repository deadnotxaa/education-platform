// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // BlogPost -.
    BlogPost struct {
        PostID             int    `json:"post_id"               example:"1"`
        AuthorID           int    `json:"author_id"             example:"1"`
        Title              string `json:"title"                 example:"Understanding Go Generics"`
        PublicationDate    string `json:"publication_date"      example:"2022-01-02"`
        Topic              string `json:"topic"                 example:"Go Programming"`
        ReadingTimeMinutes int    `json:"reading_time_minutes"  example:"5"`
        CoverImageUrl      string `json:"cover_image_url"       example:"https://example.com/cover.jpg"`
        Content            string `json:"content"               example:"This blog post explains the concept of ..."`
    }

    // Tag -.
    Tag struct {
        TagID int    `json:"id"     example:"1"`
        Name  string `json:"name"   example:"Go Generics"`
    }
)
