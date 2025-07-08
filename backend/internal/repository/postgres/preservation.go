package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PreservationRepo struct {
	db *sqlx.DB
}

func NewPreservationRepo(db *sqlx.DB) *PreservationRepo {
	return &PreservationRepo{
		db: db,
	}
}

type Preservation interface {
	Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error)
	GetLast(ctx context.Context, req *models.GetPreservationsDTO) (*models.Preservation, error)
	Create(ctx context.Context, dto *models.PreservationDTO) error
	Update(ctx context.Context, dto *models.PreservationDTO) error
	Delete(ctx context.Context, dto *models.DeletePreservationDTO) error
}

func (r *PreservationRepo) Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error) {
	query := fmt.Sprintf(`SELECT id, date_start, date_end, notes_start, notes_end, created_at FROM %s WHERE instrument_id=$1
		ORDER BY date DESC`,
		PreservationTable,
	)
	data := []*models.Preservation{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *PreservationRepo) GetLast(ctx context.Context, req *models.GetPreservationsDTO) (*models.Preservation, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date_start, date_end, notes_start, notes_end, created_at FROM %s
		WHERE instrument_id=$1 ORDER BY date_start DESC LIMIT 1`,
		PreservationTable,
	)
	data := &models.Preservation{}

	if err := r.db.GetContext(ctx, data, query, req.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *PreservationRepo) Create(ctx context.Context, dto *models.PreservationDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date_start, date_end, notes_start, notes_end)
		VALUES (:id, :instrument_id, :date_start, :date_end, :notes_start, :notes_end)`,
		PreservationTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PreservationRepo) Update(ctx context.Context, dto *models.PreservationDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date_start=:date_start, date_end=:date_end, notes_start=:notes_start, notes_end=:notes_end 
		WHERE id=:id`, PreservationTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *PreservationRepo) Delete(ctx context.Context, dto *models.DeletePreservationDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, PreservationTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
