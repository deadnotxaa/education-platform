package v1

import (
    "net/http"

    "github.com/deadnotxaa/education-platform/backend/internal/controller/http/v1/request"
    "github.com/gofiber/fiber/v2"
)

// @Summary     Get TopCoursesReport
// @Description Get TopCoursesReport
// @ID          getTopCoursesReport
// @Tags  	    report
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.TopCoursesReport
// @Failure     400 {object} response.Error
// @Router      /report/get-top-courses-report [get]
func (r *V1) getTopCoursesReport(ctx *fiber.Ctx) error {
    var body request.TopCoursesReport

    if err := ctx.BodyParser(&body); err != nil {
        r.l.Error(err, "http - v1 - getTopCoursesReport")

        return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
    }

    if err := r.v.Struct(body); err != nil {
        r.l.Error(err, "http - v1 - getTopCoursesReport")

        return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
    }

    report, err := r.p.GetTopCoursesReport(ctx.UserContext(), body.LimitNumber)
    if err != nil {
        r.l.Error(err, "http - v1 - getTopCoursesReport")

        return errorResponse(ctx, http.StatusInternalServerError, "database problems")
    }

    return ctx.Status(http.StatusOK).JSON(report)
}
