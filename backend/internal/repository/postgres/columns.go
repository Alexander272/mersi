package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ColumnRepo struct {
	db *sqlx.DB
}

func NewColumnRepo(db *sqlx.DB) *ColumnRepo {
	return &ColumnRepo{
		db: db,
	}
}

type Columns interface {
	Get(ctx context.Context, req *models.GetColumnsDTO) ([]*models.Column, error)
	Create(ctx context.Context, dto *models.ColumnsDTO) error
	CreateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error
	Update(ctx context.Context, dto *models.ColumnsDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error
	Delete(ctx context.Context, dto *models.DeleteColumnDTO) error
	DeleteAll(ctx context.Context, dto *models.DeleteColumnsDTO) error
}

func (r *ColumnRepo) Get(ctx context.Context, req *models.GetColumnsDTO) ([]*models.Column, error) {
	query := fmt.Sprintf(`SELECT id, section_id, name, field, position, type, width, COALESCE(parent_id::text, '') AS parent_id, 
		allow_sort, allow_filter, created_at
		FROM %s WHERE section_id=$1 ORDER BY position, created_at`,
		ColumnsTable,
	)
	tmp := []*models.Column{}

	if err := r.db.SelectContext(ctx, &tmp, query, req.SectionID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := make([]*models.Column, 0, len(tmp))
	for i, v := range tmp {
		if i != 0 && v.ParentID == data[len(data)-1].ID {
			data[len(data)-1].Children = append(data[len(data)-1].Children, v)
		} else {
			// v.Children = make([]*models.Column, 0)
			data = append(data, v)
		}
	}
	return data, nil
}

func (r *ColumnRepo) Create(ctx context.Context, dto *models.ColumnsDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, name, field, position, type, width, parent_id, allow_sort, allow_filter)
		VALUES (:id, :section_id, :name, :field, :position, :type, :width, :parent_id, :allow_sort, :allow_filter)`,
		ColumnsTable,
	)

	dto.ID = uuid.NewString()
	if dto.ParentID == "" {
		dto.ParentID = uuid.Nil.String()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ColumnRepo) CreateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, name, field, position, type, width, parent_id, allow_sort, allow_filter)
		VALUES (:id, :section_id, :name, :field, :position, :type, :width, :parent_id, :allow_sort, :allow_filter)`,
		ColumnsTable,
	)
	for i, d := range dto {
		dto[i].ID = uuid.NewString()
		if d.ParentID == "" {
			dto[i].ParentID = uuid.Nil.String()
		}
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ColumnRepo) Update(ctx context.Context, dto *models.ColumnsDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, field=:field, position=:position, type=:type, width=:width,
		parent_id=:parent_id, allow_sort=:allow_sort, allow_filter=:allow_filter WHERE id=:id`,
		ColumnsTable,
	)

	if dto.ParentID == "" {
		dto.ParentID = uuid.Nil.String()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ColumnRepo) UpdateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		if v.ParentID == "" {
			dto[i].ParentID = uuid.Nil.String()
		}

		tmp := []interface{}{v.ID, v.Name, v.Field, v.Position, v.Type, v.Width, v.ParentID, v.AllowSort, v.AllowFilter}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	//TODO надо проверить нормально ли отрабатывает
	query := fmt.Sprintf(`UPDATE %s AS t SET name=s.name, field=s.field, position=s.position, type=s.type, width=s.width,
		parent_id=s.parent_id, allow_sort=s.allow_sort, allow_filter=s.allow_filter 
		FROM (VALUES %s) AS s(id, name, field, position::integer, type, width::integer, parent_id::uuid, allow_sort::boolean, allow_filter::boolean) 
		WHERE t.id=s.id::integer`,
		ColumnsTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ColumnRepo) Delete(ctx context.Context, dto *models.DeleteColumnDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, ColumnsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ColumnRepo) DeleteAll(ctx context.Context, dto *models.DeleteColumnsDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE section_id=:section_id`, ColumnsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
