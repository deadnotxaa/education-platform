package v1

import (
    "net/http"

    "github.com/deadnotxaa/education-platform/backend/internal/controller/http/v1/request"
    "github.com/gofiber/fiber/v2"
)

// @Summary     Get User
// @Description Get User information by ID
// @ID          getUser
// @Tags  	    user
// @Accept      json
// @Produce     json
// @Success     200 {object} entity.User
// @Failure     400 {object} response.Error
// @Router      /user/getuser [get]
func (r *V1) getUser(ctx *fiber.Ctx) error {
    var body request.User

    if err := ctx.BodyParser(&body); err != nil {
        r.l.Error(err, "http - v1 - getUser")

        return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
    }

    if err := r.v.Struct(body); err != nil {
        r.l.Error(err, "http - v1 - getUser")

        return errorResponse(ctx, http.StatusBadRequest, "invalid request body")
    }

    user, err := r.p.GetUserById(ctx.UserContext(), body.ID)
    if err != nil {
        r.l.Error(err, "http - v1 - getUser")

        return errorResponse(ctx, http.StatusInternalServerError, "database problems")
    }

    return ctx.Status(http.StatusOK).JSON(user)
}
