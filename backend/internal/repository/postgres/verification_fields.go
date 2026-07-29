package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	CreateSeveral(ctx context.Context, dto []*models.VerificationFieldDTO) error
	Update(ctx context.Context, dto *models.VerificationFieldDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.VerificationFieldDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerFieldDTO) error
	DeleteSeveral(ctx context.Context, dto []string) error
}

func (r *VerificationFieldRepo) Get(ctx context.Context, req *models.GetVerFieldsDTO) ([]*models.VerificationField, error) {
	query := fmt.Sprintf(`SELECT id, section_id, field, label, type, position, width FROM %s WHERE section_id=$1 and "group"=$2 
		ORDER BY position`,
		VerificationFieldsTable,
	)
	data := []*models.VerificationField{}

	if err := r.db.SelectContext(ctx, &data, query, req.SectionId, req.Group); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *VerificationFieldRepo) Create(ctx context.Context, dto *models.VerificationFieldDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, field, label, type, position, "group", width)
		VALUES (:id, :section_id, :field, :label, :type, :position, :group, :width)`,
		VerificationFieldsTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationFieldRepo) CreateSeveral(ctx context.Context, dto []*models.VerificationFieldDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, field, label, type, position, "group", width)
		VALUES (:id, :section_id, :field, :label, :type, :position, :group, :width)`,
		VerificationFieldsTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationFieldRepo) Update(ctx context.Context, dto *models.VerificationFieldDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET field=:field, label=:label, type=:type, position=:position, "group"=:group, width=:width 
		WHERE id=:id`,
		VerificationFieldsTable,
	)

	res, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}

func (r *VerificationFieldRepo) UpdateSeveral(ctx context.Context, dto []*models.VerificationFieldDTO) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		tmp := []interface{}{v.Id, v.Field, v.Label, v.Type, v.Position, v.Group, v.Width}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET field=s.field, label=s.label, type=s.type, position=s.position::integer, "group"=s."group", width=s.width::integer
		FROM (VALUES %s) AS s(id, field, label, type, position, "group", width) 
		WHERE t.id=s.id::uuid`,
		VerificationFieldsTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationFieldRepo) Delete(ctx context.Context, dto *models.DeleteVerFieldDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, VerificationFieldsTable)

	res, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}

func (r *VerificationFieldRepo) DeleteSeveral(ctx context.Context, dto []string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=ANY($1)`, VerificationFieldsTable)

	if _, err := r.db.ExecContext(ctx, query, pq.Array(dto)); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
