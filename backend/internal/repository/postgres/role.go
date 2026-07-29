package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type RoleRepo struct {
	db *sqlx.DB
	Transaction
}

func NewRoleRepo(db *sqlx.DB, transaction Transaction) *RoleRepo {
	return &RoleRepo{
		db:          db,
		Transaction: transaction,
	}
}

type Role interface {
	GetAll(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	GetAllWithNames(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	Get(context.Context, string) (*models.Role, error)
	GetRules(context.Context, string) ([]string, error)
	GetByRealm(context.Context, *models.GetRoleByRealmDTO) (*models.RoleFull, error)
	GetWithRealm(context.Context, *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error)

	// New permission system methods
	GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]*models.Role, error)
	GetIdsBySlugs(ctx context.Context, slugs []string) (map[string]string, error)
	IsExists(ctx context.Context, slug string) (bool, error)
	CreateRole(ctx context.Context, tx Tx, dto *models.RoleDTO) error
	UpdateRole(ctx context.Context, tx Tx, dto *models.RoleDTO) error
	DeleteRole(ctx context.Context, tx Tx, dto *models.DeleteRoleDTO) error
	AssignPermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error
	DeletePermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error
	AssignPermissions(ctx context.Context, tx Tx, roleId string, permissionIds []string) error
	GetUserCount(ctx context.Context, roleIds []string) (map[string]int, error)

	Create(context.Context, *models.RoleDTO) error
	Update(context.Context, *models.RoleDTO) error
	Delete(context.Context, string) error
}

func (r *RoleRepo) GetAll(ctx context.Context, req *models.GetRolesDTO) ([]*models.RoleFull, error) {
	var data []*models.RoleFullDTO
	query := fmt.Sprintf(`SELECT id, name, level, description, COALESCE(extends, '{}') AS extends 
		FROM %s WHERE is_show=true ORDER BY level, name`,
		RoleTable,
	)

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	roles := []*models.RoleFull{}
	for _, rfd := range data {
		roles = append(roles, &models.RoleFull{
			ID:          rfd.ID,
			Name:        rfd.Name,
			Level:       rfd.Level,
			Extends:     rfd.Extends,
			Description: rfd.Description,
		})
	}

	return roles, nil
}

func (r *RoleRepo) GetAllWithNames(ctx context.Context, req *models.GetRolesDTO) ([]*models.RoleFull, error) {
	var data []*models.RoleFullDTO
	query := fmt.Sprintf(`SELECT r.id, name, level, CASE WHEN extends IS NOT NULL THEN
		ARRAY(SELECT name FROM roles WHERE ARRAY[id] <@ r.extends) ELSE '{}' END AS extends
		FROM %s AS r ORDER BY level, name`,
		RoleTable,
	)

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	roles := []*models.RoleFull{}
	for _, rfd := range data {
		roles = append(roles, &models.RoleFull{
			ID:          rfd.ID,
			Name:        rfd.Name,
			Level:       rfd.Level,
			Extends:     rfd.Extends,
			Description: rfd.Description,
		})
	}

	return roles, nil
}

func (r *RoleRepo) Get(ctx context.Context, roleName string) (*models.Role, error) {
	cond := "name=$1"
	params := []interface{}{roleName}

	query := fmt.Sprintf(`SELECT id, slug, name, description, is_active, is_system, is_editable, created_at, updated_at
		FROM %s WHERE %s LIMIT 1`, RoleTable, cond,
	)

	role := &models.Role{}
	if err := r.db.GetContext(ctx, role, query, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

func (r *RoleRepo) GetRules(ctx context.Context, roleName string) ([]string, error) {
	var data []models.RoleWithRuleDTO
	query := fmt.Sprintf(`SELECT r.id, name, COALESCE(extends, '{}') AS extends,
		ARRAY(SELECT DISTINCT(i.name || ':' || i.method) FROM %s AS m INNER JOIN rule_item AS i ON m.rule_item_id=i.id WHERE role_id=r.id) AS rules
		FROM %s AS r
		ORDER BY level, name`,
		RuleTable, RoleTable,
	)

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	rules := make(map[string][]string)
	extends := make(map[string]struct{})

	for _, r := range data {
		m, exist := rules[r.ID]
		if !exist {
			rules[r.ID] = r.Rules
			if r.Name == roleName {
				extends[r.ID] = struct{}{}
				for _, v := range r.Extends {
					extends[v] = struct{}{}
				}
			}
		} else {
			m = append(m, r.Rules...)
			rules[r.ID] = m
		}
	}

	for i := 1; i < len(extends); i++ {
		for _, r := range data {
			if _, exist := extends[r.ID]; exist {
				for _, v := range r.Extends {
					extends[v] = struct{}{}
				}
				break
			}
		}
	}

	roleMenu := make(map[string]struct{})
	for k := range extends {
		for _, v := range rules[k] {
			roleMenu[v] = struct{}{}
		}
	}

	result := make([]string, 0, len(roleMenu))
	for k := range roleMenu {
		result = append(result, k)
	}

	return result, nil
}

func (r *RoleRepo) GetByRealm(ctx context.Context, req *models.GetRoleByRealmDTO) (*models.RoleFull, error) {
	cond := ""
	params := []interface{}{req.UserID}
	if req.RealmID != "" {
		cond = "AND realm_id=$2"
		params = append(params, req.RealmID)
	}

	query := fmt.Sprintf(`SELECT r.id, name, description, level, COALESCE(extends, '{}') AS extends, realm_id
		FROM %s AS r
		INNER JOIN %s AS a ON a.role_id=r.id
		LEFT JOIN LATERAL (SELECT sso_id from %s WHERE id=a.user_id) AS u ON true
		WHERE sso_id=$1 %s ORDER BY level DESC, realm_id LIMIT 1`,
		RoleTable, AccessTable, UsersTable, cond,
	)
	tmp := &pq_models.RoleFull{}

	if err := r.db.GetContext(ctx, tmp, query, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	req.RealmID = tmp.RealmId
	data := &models.RoleFull{
		ID:          tmp.Id,
		Name:        tmp.Name,
		Level:       tmp.Level,
		Extends:     tmp.Extends,
		Description: tmp.Description,
	}

	return data, nil
}

func (r *RoleRepo) GetWithRealm(ctx context.Context, req *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error) {
	query := fmt.Sprintf(`SELECT r.id, name, description, level, realm_id
		FROM %s AS r
		INNER JOIN %s AS a ON a.role_id=r.id
		-- LEFT JOIN LATERAL (SELECT u.sso_id from %s AS u WHERE id=a.user_id) AS u ON true
		WHERE sso_id=$1 ORDER BY level DESC, realm_id`,
		RoleTable, AccessTable, UsersTable,
	)
	data := []*models.RoleWithRealm{}

	if err := r.db.SelectContext(ctx, &data, query, req.UserID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RoleRepo) Create(ctx context.Context, role *models.RoleDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, slug, name, level, description, is_active, is_system, is_editable)
		VALUES ($1, $2, $3, $4, $5, true, $6, true)`, RoleTable,
	)
	id := uuid.NewString()

	_, err := r.db.ExecContext(ctx, query, id, role.Slug, role.Name, role.Level, role.Description, role.IsSystem)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RoleRepo) Update(ctx context.Context, role *models.RoleDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, slug=$2, description=$3 WHERE id=$4`, RoleTable)

	_, err := r.db.ExecContext(ctx, query, role.Name, role.Slug, role.Description, role.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RoleRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RoleTable)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RoleRepo) GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error) {
	cond := "id=$1"
	params := []interface{}{req.ID}
	if req.Slug != "" {
		cond = "slug=$1"
		params = []interface{}{req.Slug}
	}

	query := fmt.Sprintf(`SELECT id, slug, name, description, is_active, is_system, is_editable, created_at, updated_at
		FROM %s WHERE %s`, RoleTable, cond,
	)

	data := &models.Role{}
	if err := r.db.GetContext(ctx, data, query, params...); err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return data, nil
}

func (r *RoleRepo) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	query := fmt.Sprintf(`SELECT id, slug, name, description, is_active, is_system, is_editable, created_at, updated_at
		FROM %s ORDER BY name`, RoleTable,
	)

	data := make([]*models.Role, 0, 10)
	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}
	return data, nil
}

func (r *RoleRepo) GetIdsBySlugs(ctx context.Context, slugs []string) (map[string]string, error) {
	query := fmt.Sprintf(`SELECT id, slug FROM %s WHERE slug = ANY($1)`, RoleTable)

	rows, err := r.db.QueryContext(ctx, query, pq.Array(slugs))
	if err != nil {
		return nil, fmt.Errorf("failed to get ids by slugs: %w", err)
	}
	defer rows.Close()

	result := make(map[string]string, len(slugs))
	for rows.Next() {
		var id, slug string
		if err := rows.Scan(&id, &slug); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[slug] = id
	}
	return result, rows.Err()
}

func (r *RoleRepo) IsExists(ctx context.Context, slug string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE slug=$1)`, RoleTable)

	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, slug); err != nil {
		return false, fmt.Errorf("failed to check role existence: %w", err)
	}
	return exists, nil
}

func (r *RoleRepo) CreateRole(ctx context.Context, tx Tx, dto *models.RoleDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, slug, name, description, is_active, is_system, is_editable)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`, RoleTable,
	)

	id := uuid.NewString()
	_, err := r.getExec(tx).ExecContext(ctx, query, id, dto.Slug, dto.Name, dto.Description, true, false, true)
	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}
	dto.ID = id
	return nil
}

func (r *RoleRepo) UpdateRole(ctx context.Context, tx Tx, dto *models.RoleDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, description=$2, slug=$3 WHERE id=$4`, RoleTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.Name, dto.Description, dto.Slug, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}
	return nil
}

func (r *RoleRepo) DeleteRole(ctx context.Context, tx Tx, dto *models.DeleteRoleDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, RoleTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

func (r *RoleRepo) AssignPermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, RolePermissionsTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.RoleId, dto.PermissionId)
	if err != nil {
		return fmt.Errorf("failed to assign permission: %w", err)
	}
	return nil
}

func (r *RoleRepo) DeletePermission(ctx context.Context, tx Tx, dto *models.RolePermissionDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id=$1 AND permission_id=$2`, RolePermissionsTable)

	_, err := r.getExec(tx).ExecContext(ctx, query, dto.RoleId, dto.PermissionId)
	if err != nil {
		return fmt.Errorf("failed to delete permission: %w", err)
	}
	return nil
}

func (r *RoleRepo) AssignPermissions(ctx context.Context, tx Tx, roleId string, permissionIds []string) error {
	if len(permissionIds) == 0 {
		return nil
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE role_id=$1`, RolePermissionsTable)
	if _, err := r.getExec(tx).ExecContext(ctx, query, roleId); err != nil {
		return fmt.Errorf("failed to clear permissions: %w", err)
	}

	for _, permId := range permissionIds {
		insertQuery := fmt.Sprintf(`INSERT INTO %s(role_id, permission_id) VALUES ($1, $2)`, RolePermissionsTable)
		if _, err := r.getExec(tx).ExecContext(ctx, insertQuery, roleId, permId); err != nil {
			return fmt.Errorf("failed to insert permission: %w", err)
		}
	}
	return nil
}

func (r *RoleRepo) GetUserCount(ctx context.Context, roleIds []string) (map[string]int, error) {
	if len(roleIds) == 0 {
		return map[string]int{}, nil
	}

	query := fmt.Sprintf(`SELECT role_id, COUNT(*) as count FROM %s WHERE role_id = ANY($1) GROUP BY role_id`, AccessTable)

	rows, err := r.db.QueryContext(ctx, query, pq.Array(roleIds))
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}
	defer rows.Close()

	result := make(map[string]int, len(roleIds))
	for rows.Next() {
		var roleId string
		var count int
		if err := rows.Scan(&roleId, &count); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		result[roleId] = count
	}
	return result, rows.Err()
}
