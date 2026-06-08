package config

import (
	"testing"
	"time"
)

func TestMustLoad_durationFields(t *testing.T) {
	t.Setenv("TEST_STALE_AFTER", "30m")
	t.Setenv("TEST_CLEANUP_INTERVAL", "5m")

	var cfg struct {
		StaleAfter      time.Duration `env:"TEST_STALE_AFTER,default=1h"`
		CleanupInterval time.Duration `env:"TEST_CLEANUP_INTERVAL,default=10m"`
	}

	MustLoad(&cfg)

	if cfg.StaleAfter != 30*time.Minute {
		t.Fatalf("StaleAfter = %v, want 30m", cfg.StaleAfter)
	}
	if cfg.CleanupInterval != 5*time.Minute {
		t.Fatalf("CleanupInterval = %v, want 5m", cfg.CleanupInterval)
	}
}

func TestMustLoad_durationDefault(t *testing.T) {
	var cfg struct {
		Interval time.Duration `env:"TEST_UNUSED_INTERVAL,default=15m"`
	}

	MustLoad(&cfg)

	if cfg.Interval != 15*time.Minute {
		t.Fatalf("Interval = %v, want 15m", cfg.Interval)
	}
}
