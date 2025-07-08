// Package entity defines main entities for business logic (services), database mapping, and
// HTTP response objects if suitable. Each logic group entity in its own file.
package entity

type (
    // CareerCenterStudent - represents a student in the career center.
    CareerCenterStudent struct {
        ID                 int    `json:"id"                        example:"1"`
        UserID             int    `json:"user_id"                   example:"1"`
        CourseID           int    `json:"course_id"                 example:"1"`
        CVUrl              string `json:"cv_url"                    example:"https://example.com/cv.pdf"`
        CareerSupportStart string `json:"career_support_start_date" example:"2022-01-01"`
        SupportPeriod      int    `json:"support_period"            example:"6"` // Support period in days
    }

    // PartnerCompany - represents a company that partners with the career center.
    PartnerCompany struct {
        CompanyID           int    `json:"company_id"              example:"1"`
        ShortName           string `json:"short_name"              example:"BigTech"`
        FullName            string `json:"full_name"               example:"BigTech Company LLC"`
        HiredGraduatesCount int    `json:"hired_graduates_count"   example:"10"`
        Requirements        string `json:"requirements"            example:"Good communication skills, knowledge of Go"`
        AgreementStatus     bool   `json:"agreement_status"        example:"active"` // Status of the partnership
    }

    // JobApplication - represents a job application submitted by a student to a partner company.
    JobApplication struct {
        ID              int    `json:"id"                        example:"1"`
        StudentID       int    `json:"student_id"                example:"1"` // ID of the student who applied
        ApplicationDate string `json:"application_date"          example:"2022-01-01"`
        Status          string `json:"status"                    example:"active"`
    }
)
