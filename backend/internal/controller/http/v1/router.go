package v1

import (
    "github.com/deadnotxaa/education-platform/backend/internal/usecase"
    "github.com/deadnotxaa/education-platform/backend/pkg/logger"
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

// NewCourseRoutes -.
func NewCourseRoutes(apiV1Group fiber.Router, p usecase.Platform, l logger.Interface) {
    r := &V1{p: p, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

    courseGroup := apiV1Group.Group("/course")
    {
        courseGroup.Get("/getcourse", r.getCourse)
    }
}

func NewUserRoutes(apiV1Group fiber.Router, p usecase.Platform, l logger.Interface) {
    r := &V1{p: p, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

    userGroup := apiV1Group.Group("/user")
    {
        userGroup.Get("/getuser", r.getUser)
    }
}

func NewReportRoutes(apiV1Group fiber.Router, p usecase.Platform, l logger.Interface) {
    r := &V1{p: p, l: l, v: validator.New(validator.WithRequiredStructEnabled())}

    userGroup := apiV1Group.Group("/report")
    {
        userGroup.Get("/get-top-courses-report", r.getTopCoursesReport)
    }
}
