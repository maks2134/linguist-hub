package httputil

import (
	"github.com/gofiber/fiber/v2"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func RespondJSON(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Success: true,
		Data:    data,
	})
}

func RespondError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*apierror.AppError); ok {
		return c.Status(appErr.StatusCode).JSON(Response{
			Success: false,
			Error: &ErrorInfo{
				Code:    appErr.Code,
				Message: appErr.Message,
			},
		})
	}

	return c.Status(500).JSON(Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		},
	})
}

func ParseBody(c *fiber.Ctx, v interface{}) error {
	if err := c.BodyParser(v); err != nil {
		return apierror.New(400, "INVALID_BODY", "Invalid request body")
	}
	return nil
}
