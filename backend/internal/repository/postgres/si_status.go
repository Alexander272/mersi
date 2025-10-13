package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type StatusRepo struct {
	db *sqlx.DB
}

func NewStatusRepo(db *sqlx.DB) *StatusRepo {
	return &StatusRepo{
		db: db,
	}
}

type Status interface {
	Get(ctx context.Context, req *models.GetSiStatusDTO) ([]*models.SiStatus, error)
	Create(ctx context.Context, dto *models.SiStatusDTO) error
	Update(ctx context.Context, dto *models.SiStatusDTO) error
	Delete(ctx context.Context, dto *models.DeleteSiStatusDTO) error
}

func (r *StatusRepo) Get(ctx context.Context, req *models.GetSiStatusDTO) ([]*models.SiStatus, error) {
	query := fmt.Sprintf(`SELECT id, section_id, position, value, label FROM %s WHERE section_id=$1 ORDER BY position`,
		SiStatusTable,
	)

	data := []*models.SiStatus{}
	if err := r.db.SelectContext(ctx, &data, query, req.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *StatusRepo) Create(ctx context.Context, dto *models.SiStatusDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, position, value, label) 
		VALUES (:id, :section_id, :position, :value, :label)`,
		SiStatusTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *StatusRepo) Update(ctx context.Context, dto *models.SiStatusDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET value=:value, label=:label, position=:position WHERE id=:id`,
		SiStatusTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *StatusRepo) Delete(ctx context.Context, dto *models.DeleteSiStatusDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, SiStatusTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
