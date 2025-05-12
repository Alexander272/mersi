package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CreateFormRepo struct {
	db *sqlx.DB
}

func NewCreateFormRepo(db *sqlx.DB) *CreateFormRepo {
	return &CreateFormRepo{
		db: db,
	}
}

type CreateForm interface {
	Get(ctx context.Context, req *models.GetCreateFormDTO) ([]*models.CreateFormStep, error)
	Create(ctx context.Context, dto *models.CreateFormFieldDTO) error
	Update(ctx context.Context, dto *models.CreateFormFieldDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.CreateFormFieldDTO) error
	Delete(ctx context.Context, dto *models.DeleteCreateFormFieldDTO) error
}

func (r *CreateFormRepo) Get(ctx context.Context, req *models.GetCreateFormDTO) ([]*models.CreateFormStep, error) {
	query := fmt.Sprintf(`SELECT id, section_id, step, step_name, field, field_name, path, type, is_required, position FROM %s
		WHERE section_id=$1 ORDER BY step, position`,
		CreatingFormTable,
	)

	tmp := []*models.CreateFormField{}
	if err := r.db.SelectContext(ctx, &tmp, query, req.SectionID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []*models.CreateFormStep{}
	for i, v := range tmp {

		if i == 0 || data[len(data)-1].Step != v.Step {
			data = append(data, &models.CreateFormStep{
				Step:     v.Step,
				StepName: v.StepName,
				Fields:   []*models.CreateFormField{v},
			})
		} else {
			data[len(data)-1].Fields = append(data[len(data)-1].Fields, v)
		}
	}
	return data, nil
}

func (r *CreateFormRepo) Create(ctx context.Context, dto *models.CreateFormFieldDTO) error {
	//TODO надо еще добавить всякие ограничения и доп параметры
	// к примеру: мин и макс для чисел, многострочность для текста, группу для файлов и тому подобное
	query := fmt.Sprintf(`INSERT INTO %s (id, section_id, step, step_name, field, field_name, path, type, is_required, position)
		VALUES (:id, :section_id, :step, :step_name, :field, :field_name, :path, :type, :is_required, :position)`,
		CreatingFormTable,
	)
	dto.ID = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *CreateFormRepo) Update(ctx context.Context, dto *models.CreateFormFieldDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET step=:step, step_name=:step_name, field=:field, type=:type, position=:position,
		field_name=:field_name, path=:path, is_required=:is_required
		WHERE id=:id`, CreatingFormTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *CreateFormRepo) UpdateSeveral(ctx context.Context, dto []*models.CreateFormFieldDTO) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		tmp := []interface{}{v.ID, v.Step, v.StepName, v.Field, v.Type, v.Position, v.FieldName, v.Path, v.IsRequired}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET step=s.step::integer, step_name=s.step_name, field=s.field, type=s.type, position=s.position::integer,
		field_name=s.field_name, path=s.path, is_required=s.is_required::boolean
		FROM (VALUES %s) AS s(id, step, step_name, field, type, position, field_name, path, is_required) 
		WHERE t.id=s.id::uuid`,
		CreatingFormTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *CreateFormRepo) Delete(ctx context.Context, dto *models.DeleteCreateFormFieldDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, CreatingFormTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
