package middleware

import (
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func Recover(log *slog.Logger) fiber.Handler {
	return recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c fiber.Ctx, e interface{}) {
			log.Error("panic recovered",
				slog.Any("error", e),
				slog.String("path", c.Path()),
			)
		},
	})
}
