package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/makskozlov/linguist-hub/internal/constants"
	"github.com/makskozlov/linguist-hub/pkg/apierror"
	"github.com/makskozlov/linguist-hub/pkg/httputil"
)

// VoiceRateLimitConfig configures per-user voice turn limits.
type VoiceRateLimitConfig struct {
	MaxTurnsPerMinute int
	Window            time.Duration
}

// DefaultVoiceRateLimitConfig returns sensible defaults.
func DefaultVoiceRateLimitConfig() VoiceRateLimitConfig {
	return VoiceRateLimitConfig{
		MaxTurnsPerMinute: 20,
		Window:            time.Minute,
	}
}

type voiceBucket struct {
	timestamps []time.Time
}

type voiceRateLimiter struct {
	mu      sync.Mutex
	buckets map[uint]*voiceBucket
	cfg     VoiceRateLimitConfig
}

func newVoiceRateLimiter(cfg VoiceRateLimitConfig) *voiceRateLimiter {
	return &voiceRateLimiter{
		buckets: make(map[uint]*voiceBucket),
		cfg:     cfg,
	}
}

func (l *voiceRateLimiter) allow(userID uint) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.cfg.Window)
	b, ok := l.buckets[userID]
	if !ok {
		b = &voiceBucket{}
		l.buckets[userID] = b
	}
	filtered := b.timestamps[:0]
	for _, ts := range b.timestamps {
		if ts.After(cutoff) {
			filtered = append(filtered, ts)
		}
	}
	b.timestamps = filtered
	if len(b.timestamps) >= l.cfg.MaxTurnsPerMinute {
		return false
	}
	b.timestamps = append(b.timestamps, now)
	return true
}

// VoiceRateLimit limits voice API turns per authenticated user.
func VoiceRateLimit(cfg VoiceRateLimitConfig) fiber.Handler {
	lim := newVoiceRateLimiter(cfg)
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return httputil.RespondError(c, apierror.ErrUnauthorized)
		}
		if !lim.allow(userID) {
			return httputil.RespondError(c, apierror.New(constants.StatusTooManyRequests, "VOICE_RATE_LIMIT", "Too many voice requests, try again later"))
		}
		return c.Next()
	}
}
