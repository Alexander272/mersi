package migrate

import (
	"database/sql"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/migrate/postgres/migrations"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(&migrations.Content)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set dialect: %w", err)
	}

	logger.Info("migration up till last")
	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to migrate up: %w", err)
	}

	return nil
}
