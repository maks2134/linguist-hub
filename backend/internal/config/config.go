package config

import (
	"fmt"
	"os"
	"strings"

	cfgloader "github.com/makskozlov/linguist-hub/pkg/config"
)

type Config struct {
	Port             string `env:"PORT,default=8080"`
	DBDsn            string `env:"DB_DSN"`
	JWTSecret        string `env:"JWT_SECRET,required"`
	LLMAPIKey        string `env:"LLM_API_KEY"`
	LLMBaseURL       string `env:"LLM_BASE_URL,default=https://openrouter.ai/api/v1"`
	LLMModel         string `env:"LLM_MODEL,default=openai/gpt-4o-mini"`
	LogLevel         string `env:"LOG_LEVEL,default=info"`
	AllowedOrigins   string `env:"ALLOWED_ORIGINS,default=http://localhost:3000"`
	FeatureServerSTT bool   `env:"FEATURE_SERVER_STT,default=false"`
	EnableSwagger    bool   `env:"ENABLE_SWAGGER,default=false"`
	AppEnv           string `env:"APP_ENV,default=development"`
}

func MustLoad() *Config {
	cfg := &Config{}
	cfgloader.MustLoad(cfg)
	cfg.resolveDSN()
	cfg.resolveLLMKey()
	if strings.TrimSpace(cfg.JWTSecret) == "" {
		panic("config: JWT_SECRET is required")
	}
	if strings.TrimSpace(cfg.LLMAPIKey) == "" {
		panic("config: LLM_API_KEY or OPENROUTER_API_KEY is required")
	}
	return cfg
}

func (c *Config) resolveDSN() {
	if strings.TrimSpace(c.DBDsn) == "" {
		c.DBDsn = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	}
	if strings.TrimSpace(c.DBDsn) == "" {
		panic("config: DB_DSN or DATABASE_URL is required")
	}
	c.DBDsn = ensureSSLMode(c.DBDsn)
}

func (c *Config) resolveLLMKey() {
	if strings.TrimSpace(c.LLMAPIKey) != "" {
		return
	}
	// Aliases for convenience / migration from Gemini.
	for _, k := range []string{"OPENROUTER_API_KEY", "GEMINI_API_KEY"} {
		if v := strings.TrimSpace(os.Getenv(k)); v != "" {
			c.LLMAPIKey = v
			return
		}
	}
}

func ensureSSLMode(dsn string) string {
	lower := strings.ToLower(dsn)
	if strings.Contains(lower, "sslmode=") {
		return dsn
	}
	if strings.HasPrefix(lower, "postgres://") || strings.HasPrefix(lower, "postgresql://") {
		sep := "?"
		if strings.Contains(dsn, "?") {
			sep = "&"
		}
		return dsn + sep + "sslmode=require"
	}
	return dsn
}

func (c *Config) OriginList() []string {
	parts := strings.Split(c.AllowedOrigins, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"http://localhost:3000"}
	}
	return out
}

func (c *Config) SwaggerEnabled() bool {
	if c.EnableSwagger {
		return true
	}
	return !strings.EqualFold(c.AppEnv, "production")
}

func (c *Config) String() string {
	return fmt.Sprintf("port=%s app_env=%s swagger=%v model=%s", c.Port, c.AppEnv, c.SwaggerEnabled(), c.LLMModel)
}
