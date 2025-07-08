package v1

import (
    "net/http"

    "github.com/deadnotxaa/education-platform/backend/internal/controller/http/v1/request"
    "github.com/gofiber/fiber/v2"
)

// @Summary     Get Course
// @Description Get Course information by ID
// @ID          getCourse
// @Tags  	    course
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.TranslationHistory
// @Failure     400 {object} response.Error
// @Router      /course/getcourse [get]
func (r *V1) getCourse(ctx *fiber.Ctx) error {
    var body request.Course

    if err := ctx.BodyParser(&body); err != nil {
        r.l.Error(err, "http - v1 - getCourse")

        return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
    }

    if err := r.v.Struct(body); err != nil {
        r.l.Error(err, "http - v1 - getCourse")

        return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
    }

    course, err := r.p.GetCourseById(ctx.UserContext(), body.ID)
    if err != nil {
        r.l.Error(err, "http - v1 - getCourse")

        return errorResponse(ctx, http.StatusInternalServerError, "database problems")
    }

    return ctx.Status(http.StatusOK).JSON(course)
}