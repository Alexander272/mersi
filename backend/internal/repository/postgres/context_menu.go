package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ContextRepo struct {
	db *sqlx.DB
}

func NewContextRepo(db *sqlx.DB) *ContextRepo {
	return &ContextRepo{
		db: db,
	}
}

type ContextMenu interface {
	Get(ctx context.Context, req *models.GetContextMenuDTO) ([]*models.ContextMenu, error)
	Create(ctx context.Context, dto *models.ContextMenuDTO) error
	Update(ctx context.Context, dto *models.ContextMenuDTO) error
	Delete(ctx context.Context, dto *models.DeleteContextMenuDTO) error
}

func (r *ContextRepo) Get(ctx context.Context, req *models.GetContextMenuDTO) ([]*models.ContextMenu, error) {
	query := fmt.Sprintf(`SELECT c.id, position, section_id, c.name, label, (r.name || ':' || r.method) AS rule FROM %s AS c
		INNER JOIN %s AS r ON rule_item_id=r.id
		WHERE section_id=$1 ORDER BY position`,
		ContextTable, RuleItemTable,
	)

	data := []*models.ContextMenu{}
	if err := r.db.SelectContext(ctx, &data, query, req.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *ContextRepo) Create(ctx context.Context, dto *models.ContextMenuDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, position, section_id, name, label, rule_item_id)
		VALUES (:id, :position, :section_id, :name, :label, :rule_item_id)`,
		ContextTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) Update(ctx context.Context, dto *models.ContextMenuDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET position=:position, name=:name, label=:label, rule_item_id=:rule_item_id, updated_at=now()
		WHERE id=:id`,
		ColumnsTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) Delete(ctx context.Context, dto *models.DeleteContextMenuDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, ContextTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
