package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type FilterRepo struct {
	db *sqlx.DB
}

func NewFilterRepo(db *sqlx.DB) *FilterRepo {
	return &FilterRepo{
		db: db,
	}
}

type Filters interface {
	Get(ctx context.Context, dto *models.GetSavedFiltersDTO) ([]*models.SavedFilter, error)
	CreateOne(ctx context.Context, dto *models.SavedFilterDTO) error
	Create(ctx context.Context, dto []*models.SavedFilterDTO) error
	Delete(ctx context.Context, dto *models.DeleteSavedFiltersDTO) error
}

func (r *FilterRepo) Get(ctx context.Context, dto *models.GetSavedFiltersDTO) ([]*models.SavedFilter, error) {
	query := fmt.Sprintf(`SELECT id, name, compare_type, value FROM %s WHERE sso_id=$1 AND section_id=$2`, FiltersTable)
	data := []*models.SavedFilter{}

	if err := r.db.SelectContext(ctx, &data, query, dto.UserId, dto.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *FilterRepo) CreateOne(ctx context.Context, dto *models.SavedFilterDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, sso_id, section_id, name, compare_type, value)
		VALUES (:id, :sso_id, :section_id, :name, :compare_type, :value)`,
		FiltersTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *FilterRepo) Create(ctx context.Context, dto []*models.SavedFilterDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, sso_id, section_id, name, compare_type, value)
		VALUES (:id, :sso_id, :section_id, :name, :compare_type, :value)`,
		FiltersTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

// func (r *FilterRepo) Update()

func (r *FilterRepo) Delete(ctx context.Context, dto *models.DeleteSavedFiltersDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE sso_id=:sso_id AND section_id=:section_id`, FiltersTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
