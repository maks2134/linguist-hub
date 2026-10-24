package middleware

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
	"github.com/makskozlov/linguist-hub/pkg/ratelimit"
)

// Default abuse limits for a public demo (in-memory; resets on deploy).
var (
	DefaultRegisterLimit = ratelimit.New(5, time.Hour)
	DefaultLoginLimit    = ratelimit.New(30, time.Hour)
	DefaultAIMinuteLimit = ratelimit.New(20, time.Minute)
	DefaultAIDailyQuota  = ratelimit.NewDailyQuota(80)
)

// IPRateLimit limits by client IP (use with TrustProxy on Render).
func IPRateLimit(lim *ratelimit.Limiter, code, message string) fiber.Handler {
	return func(c fiber.Ctx) error {
		key := c.IP()
		if key == "" {
			key = "unknown"
		}
		if !lim.Allow(key) {
			return httputil.RespondError(c, apierror.New(
				constants.StatusTooManyRequests, code, message,
			))
		}
		return c.Next()
	}
}

// UserRateLimit limits by authenticated userID in Locals.
func UserRateLimit(lim *ratelimit.Limiter, code, message string) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}
		if !lim.Allow(ratelimit.UserKey(userID)) {
			return httputil.RespondError(c, apierror.New(
				constants.StatusTooManyRequests, code, message,
			))
		}
		return c.Next()
	}
}

// AIProtect combines per-minute burst + daily LLM quota for authenticated users.
func AIProtect(minute *ratelimit.Limiter, daily *ratelimit.DailyQuota) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}
		key := ratelimit.UserKey(userID)
		if !minute.Allow(key) {
			return httputil.RespondError(c, apierror.New(
				constants.StatusTooManyRequests,
				"AI_RATE_LIMIT",
				"Too many AI requests, slow down",
			))
		}
		if !daily.Allow(key) {
			return httputil.RespondError(c, apierror.New(
				constants.StatusTooManyRequests,
				"AI_DAILY_QUOTA",
				fmt.Sprintf("Daily AI quota reached (%d). Try again tomorrow.", daily.Max()),
			))
		}
		return c.Next()
	}
}

// AIDailyQuota limits LLM-related routes by authenticated user (UTC day).
func AIDailyQuota(daily *ratelimit.DailyQuota) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}
		if !daily.Allow(ratelimit.UserKey(userID)) {
			return httputil.RespondError(c, apierror.New(
				constants.StatusTooManyRequests,
				"AI_DAILY_QUOTA",
				fmt.Sprintf("Daily AI quota reached (%d). Try again tomorrow.", daily.Max()),
			))
		}
		return c.Next()
	}
}
