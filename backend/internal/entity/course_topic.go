// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // CourseTopic -.
    CourseTopic struct {
        ID                  int    `json:"id"                       example:"1"`
        Name                string `json:"name"                     example:"Introduction to Go"`
        Description         string `json:"description"              example:"An overview of Go programming language"`
        Technologies        string `json:"technologies"             example:"Go syntax, Go tools"`
        LaborIntensityHours int    `json:"labor_intensity_hours"    example:"10"`
        ProjectsNumber      int    `json:"projects_number"          example:"1"`
    }

    // Project -.
    Project struct {
        ProjectID   int    `json:"id"                      example:"1"`
        TopicID     int    `json:"topic_id"                example:"1"`
        Name        string `json:"name"                    example:"Go Basics Project"`
        Description string `json:"description"             example:"A project to apply basic Go concepts"`
    }
)
