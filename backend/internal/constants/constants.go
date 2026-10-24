package constants

import "time"

const (
	DefaultPort           = "8080"
	DefaultLogLevel       = "info"
	DefaultAllowedOrigins = "http://localhost:3000"
	DBMaxRetries          = 5
	DBRetryDelay          = 2 * time.Second
	// Neon / serverless-friendly pool defaults (override later if needed).
	DBMaxIdleConns    = 5
	DBMaxOpenConns    = 20
	DBConnMaxLifetime = 30 * time.Minute
	JWTTokenExpiry        = time.Hour * 24 * 7
	JWTContextKey         = "user_id"
	AuthorizationHeader   = "Authorization"
	BearerPrefix          = "Bearer"
	MinPasswordLength     = 6
	GeminiModel           = "gemini-flash-latest"
	StatusOK              = 200
	StatusBadRequest      = 400
	StatusUnauthorized    = 401
	StatusForbidden       = 403
	StatusNotFound        = 404
	StatusConflict        = 409
	StatusTooManyRequests = 429
	StatusInternalError   = 500
)
