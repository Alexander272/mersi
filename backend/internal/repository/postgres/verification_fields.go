package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VerificationFieldRepo struct {
	db *sqlx.DB
}

func NewVerificationFieldRepo(db *sqlx.DB) *VerificationFieldRepo {
	return &VerificationFieldRepo{
		db: db,
	}
}

type VerificationFields interface {
	Get(ctx context.Context, req *models.GetVerFieldsDTO) ([]*models.VerificationField, error)
	Create(ctx context.Context, dto *models.VerificationFieldDTO) error
	Update(ctx context.Context, dto *models.VerificationFieldDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerFieldDTO) error
}

func (r *VerificationFieldRepo) Get(ctx context.Context, req *models.GetVerFieldsDTO) ([]*models.VerificationField, error) {
	query := fmt.Sprintf(`SELECT id, section_id, field, label, type, position FROM %s WHERE section_id=$1 ORDER BY position`,
		VerificationFieldsTable,
	)
	data := []*models.VerificationField{}

	if err := r.db.SelectContext(ctx, &data, query, req.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *VerificationFieldRepo) Create(ctx context.Context, dto *models.VerificationFieldDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, field, label, type, position)
		VALUES (:id, :section_id, :field, :label, :type, :position)`,
		VerificationFieldsTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationFieldRepo) Update(ctx context.Context, dto *models.VerificationFieldDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET field=:field, label=:label, type=:type, position=:position WHERE id=:id`, VerificationFieldsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationFieldRepo) Delete(ctx context.Context, dto *models.DeleteVerFieldDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, VerificationFieldsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
