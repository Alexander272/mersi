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
	CreateSeveral(ctx context.Context, dto []*models.HistoryTypeDTO) error
	Update(ctx context.Context, dto *models.HistoryTypeDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.HistoryTypeDTO) error
	Delete(ctx context.Context, dto *models.DeleteHistoryTypeDTO) error
	DeleteSeveral(ctx context.Context, dto []string) error
}

func (r *HistoryTypeRepo) Get(ctx context.Context, dto *models.GetHistoryTypesDTO) ([]*models.HistoryType, error) {
	query := fmt.Sprintf(`SELECT id, section_id, "group", label, position, created_at FROM %s WHERE section_id=$1 ORDER BY position`,
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

func (r *HistoryTypeRepo) CreateSeveral(ctx context.Context, dto []*models.HistoryTypeDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, "group", label, position) 
		VALUES (:id, :section_id, :group, :label, :position)`,
		HistoryTypesTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

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

func (r *HistoryTypeRepo) UpdateSeveral(ctx context.Context, dto []*models.HistoryTypeDTO) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		tmp := []interface{}{v.Id, v.Group, v.Label, v.Position}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET "group"=s."group", label=s.label, position=s.position::integer
		FROM (VALUES %s) AS s(id, "group", label, position) 
		WHERE t.id=s.id::uuid`,
		HistoryTypesTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
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

func (r *HistoryTypeRepo) DeleteSeveral(ctx context.Context, dto []string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=ANY($1)`, HistoryTypesTable)

	if _, err := r.db.ExecContext(ctx, query, pq.Array(dto)); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
