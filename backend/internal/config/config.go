package config

import (
	"time"

	cfgloader "github.com/makskozlov/linguist-hub/pkg/config"
)

type Config struct {
	Port                       string        `env:"PORT,default=8080"`
	DBDsn                      string        `env:"DB_DSN,required"`
	JWTSecret                  string        `env:"JWT_SECRET,required"`
	GeminiAPIKey               string        `env:"GEMINI_API_KEY,required"`
	LogLevel                   string        `env:"LOG_LEVEL,default=info"`
	AllowedOrigins             string        `env:"ALLOWED_ORIGINS,default=http://localhost:3000"`
	FeatureServerSTT           bool          `env:"FEATURE_SERVER_STT,default=true"`
	VoiceLessonStaleAfter      time.Duration `env:"VOICE_LESSON_STALE_AFTER,default=30m"`
	VoiceLessonCleanupInterval time.Duration `env:"VOICE_LESSON_CLEANUP_INTERVAL,default=5m"`
}

func MustLoad() *Config {
	cfg := &Config{}
	cfgloader.MustLoad(cfg)
	return cfg
}
