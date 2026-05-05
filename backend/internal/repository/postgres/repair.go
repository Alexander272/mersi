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

type RepairRepo struct {
	db *sqlx.DB
	Transaction
}

func NewRepairRepo(db *sqlx.DB, transaction Transaction) *RepairRepo {
	return &RepairRepo{
		db:          db,
		Transaction: transaction,
	}
}

type Repair interface {
	Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error)
	GetById(ctx context.Context, id string) (*models.Repair, error)
	GetLast(ctx context.Context, req *models.GetRepairDTO) (*models.Repair, error)
	Create(ctx context.Context, tx Tx, dto *models.RepairDTO) error
	CreateSeveral(ctx context.Context, dto []*models.RepairDTO) error
	Update(ctx context.Context, tx Tx, dto *models.RepairDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteRepairDTO) error
}

func (r *RepairRepo) Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, defect, work, period_start, period_end, description, created_at FROM %s
		WHERE instrument_id=$1 ORDER BY period_start DESC`,
		RepairTable,
	)
	data := []*models.Repair{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RepairRepo) GetById(ctx context.Context, id string) (*models.Repair, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, defect, work, period_start, period_end, description, created_at FROM %s
		WHERE id=$1`,
		RepairTable,
	)
	repair := &models.Repair{}

	if err := r.db.GetContext(ctx, repair, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return repair, nil
}

func (r *RepairRepo) GetLast(ctx context.Context, req *models.GetRepairDTO) (*models.Repair, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, defect, work, period_start, period_end, description, created_at FROM %s
		WHERE instrument_id=$1 ORDER BY period_start DESC LIMIT 1`,
		RepairTable,
	)
	data := &models.Repair{}

	if err := r.db.GetContext(ctx, data, query, req.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RepairRepo) Create(ctx context.Context, tx Tx, dto *models.RepairDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, defect, work, period_start, period_end, description)
		VALUES (:id, :instrument_id, :defect, :work, :period_start, :period_end, :description)`,
		RepairTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RepairRepo) CreateSeveral(ctx context.Context, dto []*models.RepairDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, defect, work, period_start, period_end, description)
		VALUES (:id, :instrument_id, :defect, :work, :period_start, :period_end, :description)`,
		RepairTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RepairRepo) Update(ctx context.Context, tx Tx, dto *models.RepairDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET defect=:defect, work=:work, period_start=:period_start, period_end=:period_end,
		description=:description WHERE id=:id`,
		RepairTable,
	)

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RepairRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteRepairDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, RepairTable)

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
