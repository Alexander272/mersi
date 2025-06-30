package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ToolsMenuRepo struct {
	db *sqlx.DB
}

func NewToolsMenuRepo(db *sqlx.DB) *ToolsMenuRepo {
	return &ToolsMenuRepo{
		db: db,
	}
}

type ToolsMenu interface {
	Get(ctx context.Context, req *models.GetToolsMenuDTO) ([]*models.ToolsMenu, error)
	Create(ctx context.Context, dto *models.ToolsMenuDTO) error
	Update(ctx context.Context, dto *models.ToolsMenuDTO) error
	// ToggleFavorite(ctx context.Context, dto *models.ChangeFavoriteDTO) error
	Delete(ctx context.Context, dto *models.DeleteToolsMenuDTO) error
}

func (r *ToolsMenuRepo) Get(ctx context.Context, req *models.GetToolsMenuDTO) ([]*models.ToolsMenu, error) {
	// query := fmt.Sprintf(`SELECT c.id, position, section_id, c.name, label, can_be_favorite, favorite,
	// 	(r.name || ':' || r.method) AS rule FROM %s AS c
	// 	INNER JOIN %s AS r ON rule_item_id=r.id
	// 	WHERE section_id=$1 ORDER BY position`,
	// 	ToolsMenuTable, RuleItemTable,
	// )

	query := fmt.Sprintf(`SELECT t.id, position, section_id, t.name, label, can_be_favorite, c.id IS NOT NULL AS favorite,
		(r.name || ':' || r.method) AS rule FROM %s AS t
		INNER JOIN %s AS r ON rule_item_id=r.id
		LEFT JOIN LATERAL (SELECT id, tools_menu_id FROM %s WHERE user_id=$2) AS c ON tools_menu_id=t.id
		WHERE section_id=$1 ORDER BY position`,
		ToolsMenuTable, RuleItemTable, CustomContextMenuTable,
	)

	data := []*models.ToolsMenu{}
	if err := r.db.SelectContext(ctx, &data, query, req.SectionId, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *ToolsMenuRepo) Create(ctx context.Context, dto *models.ToolsMenuDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, position, section_id, name, label, can_be_favorite, rule_item_id)
		VALUES (:id, :position, :section_id, :name, :label, :can_be_favorite, :rule_item_id)`,
		ToolsMenuTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ToolsMenuRepo) Update(ctx context.Context, dto *models.ToolsMenuDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET position=:position, name=:name, label=:label, can_be_favorite=:can_be_favorite, 
		rule_item_id=:rule_item_id WHERE id=:id`,
		ToolsMenuTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

// func (r *ToolsMenuRepo) ToggleFavorite(ctx context.Context, dto *models.ChangeFavoriteDTO) error {
// 	query := fmt.Sprintf(`UPDATE %s SET favorite=:favorite WHERE id=:id`, ToolsMenuTable)

// 	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
// 		return fmt.Errorf("failed to execute query. error: %w", err)
// 	}
// 	return nil
// }

func (r *ToolsMenuRepo) Delete(ctx context.Context, dto *models.DeleteToolsMenuDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, ToolsMenuTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
