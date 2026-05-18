package postgres

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/makskozlov/linguist-hub/internal/constants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	maxRetries := constants.DBMaxRetries
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err == nil {
			sqlDB, err := db.DB()
			if err != nil {
				return nil, fmt.Errorf("failed to get database instance: %w", err)
			}

			if err := sqlDB.Ping(); err == nil {
				sqlDB.SetMaxIdleConns(constants.DBMaxIdleConns)
				sqlDB.SetMaxOpenConns(constants.DBMaxOpenConns)
				sqlDB.SetConnMaxLifetime(constants.DBConnMaxLifetime)
				slog.Info("successfully connected to database")
				return db, nil
			}
		}

		slog.Error("failed to connect to database", "attempt", i+1, "max_retries", maxRetries, "error", err)
		time.Sleep(constants.DBRetryDelay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
