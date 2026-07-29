package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleHierarchyRepo struct {
	db *sqlx.DB
	Transaction
}

func NewRoleHierarchyRepo(db *sqlx.DB, tr Transaction) *RoleHierarchyRepo {
	return &RoleHierarchyRepo{
		db:          db,
		Transaction: tr,
	}
}

type RoleHierarchy interface {
	LoadPolicy(ctx context.Context) ([]*models.SyncRoleInheritance, error)
	GetRoleDescendants(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error)
	GetDirectChildren(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error)
	AddInheritance(ctx context.Context, tx Tx, dto *models.RoleHierarchyDTO) error
	AddInheritances(ctx context.Context, tx Tx, roleId string, parentRoleIds []string) error
	RemoveInheritance(ctx context.Context, tx Tx, dto *models.RoleHierarchyDTO) error
	RemoveInheritances(ctx context.Context, tx Tx, roleId string, parentRoleIds []string) error
}

func (r *RoleHierarchyRepo) GetRoleDescendants(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error) {
	query := fmt.Sprintf(`WITH RECURSIVE descendants_tree AS (
			SELECT
				r1.id as root_id,
				r1.name as root_name,
				r2.id as child_id,
				r2.name as child_name
			FROM %s ri
			JOIN %s r1 ON ri.parent_role_id = r1.id
			JOIN %s r2 ON ri.role_id = r2.id
			WHERE r1.name = ANY($1)
			AND r2.is_active = true

			UNION ALL

			SELECT
				dt.root_id,
				dt.root_name,
				r3.id,
				r3.name
			FROM descendants_tree dt
			JOIN %s ri ON ri.parent_role_id = dt.child_id
			JOIN %s r3 ON ri.role_id = r3.id
			WHERE r3.is_active = true
		)
		SELECT DISTINCT root_name, child_name
		FROM descendants_tree`,
		RoleHierarchyTable, RoleTable, RoleTable,
		RoleHierarchyTable, RoleTable,
	)

	rows, err := r.db.QueryContext(ctx, query, pq.Array(req.Roles))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]string)
	for _, name := range req.Roles {
		result[name] = []string{}
	}

	for rows.Next() {
		var root, child string
		if err := rows.Scan(&root, &child); err != nil {
			return nil, err
		}
		result[root] = append(result[root], child)
	}

	return result, nil
}

func (r *RoleHierarchyRepo) GetDirectChildren(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error) {
	query := fmt.Sprintf(`SELECT r1.name, r2.name
		FROM %s ri
		JOIN %s r1 ON ri.parent_role_id = r1.id
		JOIN %s r2 ON ri.role_id = r2.id
		WHERE r1.name = ANY($1) AND r2.is_active = true`,
		RoleHierarchyTable, RoleTable, RoleTable,
	)

	rows, err := r.db.QueryContext(ctx, query, pq.Array(req.Roles))
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	result := make(map[string][]string)
	for _, name := range req.Roles {
		result[name] = []string{}
	}

	for rows.Next() {
		var parent, child string
		if err := rows.Scan(&parent, &child); err != nil {
			return nil, err
		}
		result[parent] = append(result[parent], child)
	}

	return result, nil
}

func (r *RoleHierarchyRepo) LoadPolicy(ctx context.Context) ([]*models.SyncRoleInheritance, error) {
	query := fmt.Sprintf(`SELECT r1.name, r2.name, d.id
        FROM %s rh
        JOIN %s r1 ON rh.role_id = r1.id
        JOIN %s r2 ON rh.parent_role_id = r2.id
		JOIN %s d ON true
        WHERE r1.is_active = true AND r2.is_active = true`,
		RoleHierarchyTable, RoleTable, RoleTable, RealmTable,
	)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	data := make([]*models.SyncRoleInheritance, 0, 5)

	for rows.Next() {
		item := &models.SyncRoleInheritance{}
		if err := rows.Scan(&item.Role, &item.ParentRole, &item.Realm); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		data = append(data, item)
	}

	return data, nil
}

func (r *RoleHierarchyRepo) AddInheritance(ctx context.Context, tx Tx, dto *models.RoleHierarchyDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (role_id, parent_role_id) 
		VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		RoleHierarchyTable,
	)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.RoleId, dto.ParentRoleId)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleHierarchyRepo) AddInheritances(ctx context.Context, tx Tx, roleId string, parentRoleIds []string) error {
	if len(parentRoleIds) == 0 {
		return nil
	}

	values := make([]string, 0, len(parentRoleIds))
	args := make([]interface{}, 0, len(parentRoleIds)*2)
	for i, parentId := range parentRoleIds {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, roleId, parentId)
	}

	query := fmt.Sprintf(`INSERT INTO %s (role_id, parent_role_id) VALUES %s ON CONFLICT DO NOTHING`,
		RoleHierarchyTable, strings.Join(values, ", "))

	_, err := r.getExec(tx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleHierarchyRepo) RemoveInheritance(ctx context.Context, tx Tx, dto *models.RoleHierarchyDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id = $1 AND parent_role_id = $2`,
		RoleHierarchyTable,
	)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.RoleId, dto.ParentRoleId)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *RoleHierarchyRepo) RemoveInheritances(ctx context.Context, tx Tx, roleId string, parentRoleIds []string) error {
	if len(parentRoleIds) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(parentRoleIds))
	args := []interface{}{roleId}
	for _, parentID := range parentRoleIds {
		placeholders = append(placeholders, fmt.Sprintf("$%d", len(args)+1))
		args = append(args, parentID)
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id = $1 AND parent_role_id IN (%s)`,
		RoleHierarchyTable, strings.Join(placeholders, ", "))

	_, err := r.getExec(tx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
