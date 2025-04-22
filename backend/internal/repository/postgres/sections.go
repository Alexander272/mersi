package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type SectionRepo struct {
	db *sqlx.DB
}

func NewSectionRepo(db *sqlx.DB) *SectionRepo {
	return &SectionRepo{
		db: db,
	}
}

type Section interface {
	Get(ctx context.Context, req *models.GetSectionsDTO) ([]*models.Section, error)
	Create(ctx context.Context, dto *models.SectionDTO) error
	Update(ctx context.Context, dto *models.SectionDTO) error
	Delete(ctx context.Context, dto *models.DeleteSectionDTO) error
}

func (r *SectionRepo) Get(ctx context.Context, req *models.GetSectionsDTO) ([]*models.Section, error) {
	query := fmt.Sprintf(`SELECT id, name, realm_id, position FROM %s WHERE realm_id=$1 ORDER BY position`, SectionTable)
	data := []*models.Section{}

	if err := r.db.SelectContext(ctx, &data, query, req.RealmID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *SectionRepo) Create(ctx context.Context, dto *models.SectionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, name, realm_id, position) VALUES (:id, :name, :realm_id, :position)`, SectionTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SectionRepo) Update(ctx context.Context, dto *models.SectionDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, realm_id=:realm_id, position=:position WHERE id=:id`, SectionTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SectionRepo) Delete(ctx context.Context, dto *models.DeleteSectionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, SectionTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
