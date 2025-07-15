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

type TransferToSaveRepo struct {
	db *sqlx.DB
}

func NewTransferToSaveRepo(db *sqlx.DB) *TransferToSaveRepo {
	return &TransferToSaveRepo{
		db: db,
	}
}

type TransferToSave interface {
	Get(ctx context.Context, req *models.GetTransferToSaveDTO) ([]*models.TransferToSave, error)
	GetLast(ctx context.Context, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error)
	Create(ctx context.Context, dto *models.TransferToSaveDTO) error
	Update(ctx context.Context, dto *models.TransferToSaveDTO) error
	Delete(ctx context.Context, dto *models.DeleteTransferToSaveDTO) error
}

func (r *TransferToSaveRepo) Get(ctx context.Context, req *models.GetTransferToSaveDTO) ([]*models.TransferToSave, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date_start, notes_start, date_end, notes_end, created_at
		FROM %s WHERE instrument_id=$1`,
		TransferToSaveTable,
	)
	data := []*models.TransferToSave{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *TransferToSaveRepo) GetLast(ctx context.Context, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date_start, notes_start, date_end, notes_end, created_at
		FROM %s WHERE instrument_id=$1 ORDER BY date_start DESC LIMIT 1`,
		TransferToSaveTable,
	)
	data := &models.TransferToSave{}

	if err := r.db.GetContext(ctx, data, query, req.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *TransferToSaveRepo) Create(ctx context.Context, dto *models.TransferToSaveDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date_start, notes_start, date_end, notes_end)
		VALUES (:id, :instrument_id, :date_start, :notes_start, :date_end, :notes_end)`,
		TransferToSaveTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *TransferToSaveRepo) Update(ctx context.Context, dto *models.TransferToSaveDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date_start=:date_start, notes_start=:notes_start, date_end=:date_end, 
		notes_end=:notes_end WHERE id=:id`,
		TransferToSaveTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *TransferToSaveRepo) Delete(ctx context.Context, dto *models.DeleteTransferToSaveDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, TransferToSaveTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
