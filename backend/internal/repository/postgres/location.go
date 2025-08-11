package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type LocationRepo struct {
	db *sqlx.DB
}

func NewLocationRepo(db *sqlx.DB) *LocationRepo {
	return &LocationRepo{
		db: db,
	}
}

type Location interface{}

func (r *LocationRepo) GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date_of_issue, date_of_receiving, status, need_confirmed,
		COALESCE(person_id::text, '') AS person_id, COALESCE(department_id::text, '') AS department_id
		FROM %s WHERE instrument_id=$1 ORDER BY date_of_issue DESC, created_at DESC LIMIT 1`,
		LocationTable,
	)
	location := &models.Location{}

	if err := r.db.GetContext(ctx, location, query, dto.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	return location, nil
}
