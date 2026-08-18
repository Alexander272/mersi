package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CustomContextRepo struct {
	db *sqlx.DB
}

func NewCustomContextRepo(db *sqlx.DB) *CustomContextRepo {
	return &CustomContextRepo{
		db: db,
	}
}

type CustomContextMenu interface {
	Create(ctx context.Context, dto *models.CustomContextMenuDTO) error
	Delete(ctx context.Context, dto *models.CustomContextMenuDTO) error
}

func (r *CustomContextRepo) Create(ctx context.Context, dto *models.CustomContextMenuDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, tools_menu_id) VALUES (:id, :user_id, :tools_menu_id)`, CustomContextMenuTable)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *CustomContextRepo) Delete(ctx context.Context, dto *models.CustomContextMenuDTO) error {
	// query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, CustomContextMenuTable)
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id=:user_id AND tools_menu_id=:tools_menu_id`, CustomContextMenuTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
