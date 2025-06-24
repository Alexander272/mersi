package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RuleItemRepo struct {
	db *sqlx.DB
}

func NewRuleItemRepo(db *sqlx.DB) *RuleItemRepo {
	return &RuleItemRepo{
		db: db,
	}
}

type RuleItem interface {
	GetAll(context.Context) ([]*models.RuleItem, error)
	Create(context.Context, *models.RuleItemDTO) error
	Update(context.Context, *models.RuleItemDTO) error
	Delete(context.Context, string) error
}

func (r *RuleItemRepo) GetAll(ctx context.Context) ([]*models.RuleItem, error) {
	query := fmt.Sprintf(`SELECT id, name, method, description, is_show FROM %s`, RuleItemTable)
	items := []*models.RuleItem{}

	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return items, nil
}

func (r *RuleItemRepo) Create(ctx context.Context, menu *models.RuleItemDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, name, method, description, is_show) VALUES ($1, $2, $3, $4, $5)`, RuleItemTable)
	id := uuid.New()

	_, err := r.db.ExecContext(ctx, query, id, menu.Name, menu.Method, menu.Description, menu.IsShow)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RuleItemRepo) Update(ctx context.Context, menu *models.RuleItemDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, method=$2 description=$3, is_show=$4 WHERE id=$5`, RuleItemTable)

	_, err := r.db.ExecContext(ctx, query, menu.Name, menu.Method, menu.Description, menu.IsShow, menu.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RuleItemRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RuleItemTable)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
