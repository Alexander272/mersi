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

type PermissionRepo struct {
	db *sqlx.DB
	Transaction
}

func NewPermissionRepo(db *sqlx.DB, tr Transaction) *PermissionRepo {
	return &PermissionRepo{
		db:          db,
		Transaction: tr,
	}
}

type Permissions interface {
	LoadPolicy(ctx context.Context) ([]*models.Permission, error)
	Sync(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error
	GetById(ctx context.Context, id string) (*models.Permission, error)
	GetAll(ctx context.Context) ([]*models.Permission, error)
	GetByRole(ctx context.Context, req *models.GetPermsByRoleDTO) ([]*models.Permission, error)
	GetInheritedByRole(ctx context.Context, roleId string) (map[string]struct{}, error)
	GetRolePermissionsMap(ctx context.Context, roleId string) (map[string]bool, error)
	ReplacePermissions(ctx context.Context, tx Tx, roleId string, permissionIds []string) error
	Create(ctx context.Context, tx Tx, dto *models.PermissionDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeletePermissionDTO) error
	DeleteByKeys(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error
}

func (r *PermissionRepo) LoadPolicy(ctx context.Context) ([]*models.Permission, error) {
	query := fmt.Sprintf(`SELECT r.slug, d.id, p.object, p.action
		FROM %s rp
		JOIN %s r ON r.id = rp.role_id
		INNER JOIN %s d ON true
		JOIN %s p ON p.id = rp.permission_id`,
		RolePermissionsTable, RoleTable, RealmTable, PermissionsTable,
	)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	permissions := make([]*models.Permission, 0, 50)
	for rows.Next() {
		item := &models.Permission{}
		if err := rows.Scan(&item.Role, &item.Realm, &item.Object, &item.Action); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		permissions = append(permissions, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}
	return permissions, nil
}

func (r *PermissionRepo) Sync(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error {
	if len(dto) == 0 {
		return nil
	}
	values := []string{}
	args := []interface{}{}

	for _, v := range dto {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			len(args)+1, len(args)+2, len(args)+3, len(args)+4, len(args)+5,
		))
		args = append(args, uuid.New(), v.Object, v.Action, v.Name, v.Description)
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, object, action, name, description)
			VALUES %s
			ON CONFLICT (object, action) 
			DO UPDATE SET description = EXCLUDED.description, name = EXCLUDED.name`,
		PermissionsTable, strings.Join(values, ", "),
	)

	_, err := r.getExec(tx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) GetById(ctx context.Context, id string) (*models.Permission, error) {
	query := fmt.Sprintf(`SELECT id, object, action FROM %s WHERE id=$1`, PermissionsTable)
	data := &models.Permission{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&data.Id, &data.Object, &data.Action)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	return data, nil
}

func (r *PermissionRepo) GetAll(ctx context.Context) ([]*models.Permission, error) {
	query := fmt.Sprintf(`SELECT id, object, action, name, description FROM %s ORDER BY object, action`, PermissionsTable)

	data := make([]*models.Permission, 0, 50)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.Permission{}
		if err := rows.Scan(&item.Id, &item.Object, &item.Action, &item.Name, &item.Description); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		data = append(data, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return data, nil
}

func (r *PermissionRepo) GetByRole(ctx context.Context, req *models.GetPermsByRoleDTO) ([]*models.Permission, error) {
	query := fmt.Sprintf(`SELECT p.id, p.object, p.action, p.name, p.description
		FROM %s rp
		JOIN %s r ON r.id = rp.role_id
		JOIN %s p ON p.id = rp.permission_id
		WHERE r.slug = $1`,
		RolePermissionsTable, RoleTable, PermissionsTable,
	)

	data := make([]*models.Permission, 0, 50)
	rows, err := r.db.QueryContext(ctx, query, req.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &models.Permission{}
		if err := rows.Scan(&item.Id, &item.Object, &item.Action, &item.Name, &item.Description); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		data = append(data, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return data, nil
}

func (r *PermissionRepo) GetRolePermissionsMap(ctx context.Context, roleId string) (map[string]bool, error) {
	query := fmt.Sprintf(`SELECT permission_id FROM %s WHERE role_id = $1`, RolePermissionsTable)

	rows, err := r.db.QueryContext(ctx, query, roleId)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	result := make(map[string]bool)
	for rows.Next() {
		var permID string
		if err := rows.Scan(&permID); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		result[permID] = true
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func (r *PermissionRepo) GetInheritedByRole(ctx context.Context, roleID string) (map[string]struct{}, error) {
	query := fmt.Sprintf(`WITH RECURSIVE sub_roles AS (
			SELECT role_id 
			FROM %s 
			WHERE parent_role_id = $1
			
			UNION ALL
			
			SELECT rh.role_id 
			FROM %s rh
			JOIN sub_roles sr ON rh.parent_role_id = sr.role_id
		)
		SELECT DISTINCT rp.permission_id 
		FROM %s rp 
		WHERE rp.role_id IN (SELECT role_id FROM sub_roles)`,
		RoleHierarchyTable, RoleHierarchyTable, RolePermissionsTable,
	)

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	result := make(map[string]struct{})
	for rows.Next() {
		var permID string
		if err := rows.Scan(&permID); err != nil {
			return nil, fmt.Errorf("scan row error: %w", err)
		}
		result[permID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return result, nil
}

func (r *PermissionRepo) ReplacePermissions(ctx context.Context, tx Tx, roleId string, permissionIDs []string) error {
	exec := r.getExec(tx)

	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id = $1`, RolePermissionsTable)
	_, err := exec.ExecContext(ctx, query, roleId)
	if err != nil {
		return fmt.Errorf("failed to delete old permissions: %w", err)
	}

	if len(permissionIDs) == 0 {
		return nil
	}

	values := make([]string, 0, len(permissionIDs))
	args := make([]interface{}, 0, len(permissionIDs)*2)
	for i, id := range permissionIDs {
		values = append(values, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, roleId, id)
	}

	query = fmt.Sprintf(`INSERT INTO %s (role_id, permission_id) VALUES %s`, RolePermissionsTable, strings.Join(values, ", "))

	_, err = exec.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert permissions: %w", err)
	}

	return nil
}

func (r *PermissionRepo) Create(ctx context.Context, tx Tx, dto *models.PermissionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, object, action, name, description) VALUES ($1, $2, $3, $4, $5)`,
		PermissionsTable,
	)
	dto.Id = uuid.NewString()

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.Id, dto.Object, dto.Action, dto.Name, dto.Description)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) Delete(ctx context.Context, tx Tx, dto *models.DeletePermissionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, PermissionsTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) DeleteByKeys(ctx context.Context, tx Tx, dto []*models.PermissionDTO) error {
	if len(dto) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(dto)*2)
	args := make([]interface{}, 0, len(dto)*2)
	for _, v := range dto {
		placeholders = append(placeholders, fmt.Sprintf("($%d::text, $%d::text)", len(args)+1, len(args)+2))
		args = append(args, v.Object, v.Action)
	}

	query := fmt.Sprintf(`DELETE FROM %s 
        WHERE (object, action) NOT IN (
            SELECT * FROM (VALUES %s) AS t(obj, act)
        )`,
		PermissionsTable,
		strings.Join(placeholders, ","),
	)

	_, err := r.getExec(tx).ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) AssignPermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (role_id, permission_id) VALUES ($1, $2)`, RolePermissionsTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.RoleId, dto.PermissionId)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *PermissionRepo) DeletePermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id=$1 AND permission_id=$2`, RolePermissionsTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.RoleId, dto.PermissionId)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

// Ensure pq is used (for potential array operations)
var _ = pq.Array
