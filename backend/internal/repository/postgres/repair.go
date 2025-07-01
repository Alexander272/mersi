package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RepairRepo struct {
	db *sqlx.DB
}

func NewRepairRepo(db *sqlx.DB) *RepairRepo {
	return &RepairRepo{
		db: db,
	}
}

type Repair interface {
	Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error)
	Create(ctx context.Context, dto *models.RepairDTO) error
	Update(ctx context.Context, dto *models.RepairDTO) error
	Delete(ctx context.Context, dto *models.DeleteRepairDTO) error
}

func (r *RepairRepo) Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error) {
	query := fmt.Sprintf(`SELECT id, defect, work, period_start, period_end, description, created_at FROM %s
		WHERE instrument_id=$1`,
		RepairTable,
	)
	data := []*models.Repair{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RepairRepo) Create(ctx context.Context, dto *models.RepairDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, defect, work, period_start, period_end, description)
		VALUES (:id, :instrument_id, :defect, :work, :period_start, :period_end, :description)`,
		RepairTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RepairRepo) Update(ctx context.Context, dto *models.RepairDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET defect=:defect, work=:work, period_start=:period_start, period_end=:period_end
		description=:description WHERE id=:id`,
		RepairTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RepairRepo) Delete(ctx context.Context, dto *models.DeleteRepairDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, RepairTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
