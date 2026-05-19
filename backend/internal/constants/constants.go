package constants

import "time"

const (
	DefaultPort           = "8080"
	DefaultLogLevel       = "info"
	DefaultAllowedOrigins = "http://localhost:3000"
	DBMaxRetries          = 5
	DBRetryDelay          = 2 * time.Second
	DBMaxIdleConns        = 10
	DBMaxOpenConns        = 100
	DBConnMaxLifetime     = time.Hour
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
	StatusInternalError   = 500
)
