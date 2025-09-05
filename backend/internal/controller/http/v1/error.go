package v1

import (
    "github.com/deadnotxaa/education-platform/backend/internal/controller/http/v1/response"
    "github.com/gofiber/fiber/v2"
)

func errorResponse(ctx *fiber.Ctx, code int, msg string) error {
    return ctx.Status(code).JSON(response.Error{Error: msg})
}