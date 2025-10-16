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

type ContextRepo struct {
	db *sqlx.DB
}

func NewContextRepo(db *sqlx.DB) *ContextRepo {
	return &ContextRepo{
		db: db,
	}
}

type ContextMenu interface {
	Get(ctx context.Context, req *models.GetContextMenuDTO) ([]*models.ContextMenu, error)
	Create(ctx context.Context, dto *models.ContextMenuDTO) error
	CreateSeveral(ctx context.Context, dto []*models.ContextMenuDTO) error
	Update(ctx context.Context, dto *models.ContextMenuDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.ContextMenuDTO) error
	Delete(ctx context.Context, dto *models.DeleteContextMenuDTO) error
	DeleteSeveral(ctx context.Context, dto []string) error
}

func (r *ContextRepo) Get(ctx context.Context, req *models.GetContextMenuDTO) ([]*models.ContextMenu, error) {
	// query := fmt.Sprintf(`SELECT c.id, position, section_id, c.name, label, (r.name || ':' || r.method) AS rule FROM %s AS c
	// 	INNER JOIN %s AS r ON rule_item_id=r.id
	// 	WHERE section_id=$1 ORDER BY position`,
	// 	ContextTable, RuleItemTable,
	// )

	query := fmt.Sprintf(`WITH Items AS (
		(SELECT id, position, section_id, name, label, rule_item_id FROM %s
			WHERE section_id=$1 ORDER BY position)
		UNION ALL
		(SELECT c.id, position+10, section_id, name, label, rule_item_id FROM %s AS c
			INNER JOIN tools_menu AS t ON tools_menu_id=t.id
			WHERE section_id=$1 AND user_id=$2 ORDER BY position)
	)
	
	SELECT c.id, position, section_id, c.name, label, (r.name || ':' || r.method) AS rule
		FROM Items AS c INNER JOIN %s AS r ON rule_item_id=r.id ORDER BY position`,
		ContextTable, CustomContextMenuTable, RuleItemTable,
	)

	data := []*models.ContextMenu{}
	if err := r.db.SelectContext(ctx, &data, query, req.SectionId, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *ContextRepo) Create(ctx context.Context, dto *models.ContextMenuDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, position, section_id, name, label, rule_item_id)
		VALUES (:id, :position, :section_id, :name, :label, :rule_item_id)`,
		ContextTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) CreateSeveral(ctx context.Context, dto []*models.ContextMenuDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, position, section_id, name, label, rule_item_id)
		VALUES (:id, :position, :section_id, :name, :label, :rule_item_id)`,
		ContextTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) Update(ctx context.Context, dto *models.ContextMenuDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET position=:position, name=:name, label=:label, rule_item_id=:rule_item_id, updated_at=now()
		WHERE id=:id`,
		ContextTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) UpdateSeveral(ctx context.Context, dto []*models.ContextMenuDTO) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		tmp := []interface{}{v.Id, v.Position, v.Name, v.Label, v.RuleItemId}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET position=s.position::integer, name=s.name, label=s.label, rule_item_id=s.rule_item_id, updated_at=now()
		FROM (VALUES %s) AS s(id, position, name, label, rule_item_id) 
		WHERE t.id=s.id::uuid`,
		ContextTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) Delete(ctx context.Context, dto *models.DeleteContextMenuDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, ContextTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *ContextRepo) DeleteSeveral(ctx context.Context, dto []string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=ANY(:id)`, ContextTable)

	if _, err := r.db.NamedExecContext(ctx, query, pq.Array(dto)); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
