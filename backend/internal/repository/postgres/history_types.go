package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type HistoryTypeRepo struct {
	db *sqlx.DB
}

func NewHistoryTypeRepo(db *sqlx.DB) *HistoryTypeRepo {
	return &HistoryTypeRepo{
		db: db,
	}
}

type HistoryType interface {
	Get(ctx context.Context, dto *models.GetHistoryTypesDTO) ([]*models.HistoryType, error)
	Create(ctx context.Context, dto *models.HistoryTypeDTO) error
	Update(ctx context.Context, dto *models.HistoryTypeDTO) error
	Delete(ctx context.Context, dto *models.DeleteHistoryTypeDTO) error
}

func (r *HistoryTypeRepo) Get(ctx context.Context, dto *models.GetHistoryTypesDTO) ([]*models.HistoryType, error) {
	query := fmt.Sprintf(`SELECT id, "group", label, position, created_at FROM %s WHERE section_id=$1 ORDER BY position`,
		HistoryTypesTable,
	)
	data := []*models.HistoryType{}

	if err := r.db.SelectContext(ctx, &data, query, dto.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil

}

func (r *HistoryTypeRepo) Create(ctx context.Context, dto *models.HistoryTypeDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, "group", label, position) 
		VALUES (:id, :section_id, :group, :label, :position)`,
		HistoryTypesTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *HistoryTypeRepo) Update(ctx context.Context, dto *models.HistoryTypeDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET group=:group, label=:label, position=:position WHERE id=:id`,
		HistoryTypesTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *HistoryTypeRepo) Delete(ctx context.Context, dto *models.DeleteHistoryTypeDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, HistoryTypesTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
