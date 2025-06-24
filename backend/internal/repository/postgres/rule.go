package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RuleRepo struct {
	db *sqlx.DB
}

func NewRuleRepo(db *sqlx.DB) *RuleRepo {
	return &RuleRepo{
		db: db,
	}
}

type Rule interface {
	GetAll(context.Context) ([]*models.Rule, error)
	Create(context.Context, *models.RuleDTO) error
	Update(context.Context, *models.RuleDTO) error
	Delete(context.Context, string) error
}

// func (r *RuleRepo) GetAll(ctx context.Context) ([]*models.Rule, error) {
// 	var data []*pq_models.RuleDTO
// 	query := fmt.Sprintf(`SELECT m.id, role_id, name, level, Rule_item_id, CASE WHEN extends IS NOT NULL THEN
// 		ARRAY(SELECT name FROM %s WHERE ARRAY[id] <@ r.extends) ELSE '{}' END AS extends
// 		FROM %s AS m INNER JOIN %s AS r ON role_id=r.id ORDER BY level, name`,
// 		RoleTable, RuleTable, RoleTable,
// 	)

// 	if err := r.db.SelectContext(ctx, &data, query); err != nil {
// 		return nil, fmt.Errorf("failed to execute query. error: %w", err)
// 	}

// 	Rule := []*models.Rule{}
// 	for _, mpd := range data {
// 		Rule = append(Rule, &models.Rule{
// 			Id:          mpd.Id,
// 			RoleId:      mpd.RoleId,
// 			RoleName:    mpd.RoleName,
// 			RoleLevel:   mpd.RoleLevel,
// 			RoleExtends: mpd.RoleExtends,
// 			RuleItemId:  mpd.RuleItemId,
// 		})
// 	}

// 	return Rule, nil
// }

func (r *RuleRepo) GetAll(ctx context.Context) ([]*models.Rule, error) {
	query := fmt.Sprintf(`SELECT m.id, r.name, role_id, rule_item_id, i.name AS item_name, i.method
		FROM %s AS m INNER JOIN %s AS r ON role_id=r.id INNER JOIN %s AS i ON i.id=rule_item_id ORDER BY level`,
		RuleTable, RoleTable, RuleItemTable,
	)

	var Rule []*models.Rule
	if err := r.db.SelectContext(ctx, &Rule, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return Rule, nil
}

func (r *RuleRepo) Create(ctx context.Context, Rule *models.RuleDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, role_id, rule_item_id) VALUES ($1, $2, $3)`, RuleTable)
	id := uuid.New()

	_, err := r.db.ExecContext(ctx, query, id, Rule.RoleId, Rule.RuleItemId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RuleRepo) Update(ctx context.Context, Rule *models.RuleDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET role_id=$1, rule_item_id=$2 WHERE id=$3`, RuleTable)

	_, err := r.db.ExecContext(ctx, query, Rule.RoleId, Rule.RuleItemId, Rule.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RuleRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RuleTable)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
