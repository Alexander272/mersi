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
}

func NewRoleRepo(db *sqlx.DB) *RoleRepo {
	return &RoleRepo{
		db: db,
	}
}

type Role interface {
	GetAll(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	GetAllWithNames(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	Get(context.Context, string) (*models.Role, error)
	GetByRealm(context.Context, *models.GetRoleByRealmDTO) (*models.RoleFull, error)
	GetWithRealm(context.Context, *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error)
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
	var data []models.RoleWithRuleDTO
	// query := fmt.Sprintf(`SELECT r.id, r.name, COALESCE(extends, '{}') AS extends, i.name AS menu
	// 	FROM %s AS r
	// 	LEFT JOIN %s AS m ON r.id=role_id
	// 	LEFT JOIN %s AS i ON rule_item_id=i.id
	// 	WHERE i.is_show=true ORDER BY level`,
	// 	RoleTable, MenuTable, MenuItemTable,
	// )
	query := fmt.Sprintf(`SELECT r.id, name, COALESCE(extends, '{}') AS extends,
		ARRAY(SELECT DISTINCT(i.name || ':' || i.method) FROM %s AS m INNER JOIN rule_item AS i ON m.rule_item_id=i.id WHERE role_id=r.id) AS rules
		FROM %s AS r
		ORDER BY level, name`,
		RuleTable, RoleTable,
	)

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	role := &models.Role{}
	rules := make(map[string][]string, 0)
	extends := make(map[string]struct{})

	//EDIT Возможно можно это как-то покрасивее написать
	for _, r := range data {
		m, exist := rules[r.ID]
		if !exist {
			rules[r.ID] = r.Rules

			if r.Name == roleName {
				role.ID = r.ID
				role.Name = r.Name
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

	roleMenu := make(map[string]struct{}, 0)
	for k := range extends {
		for _, v := range rules[k] {
			roleMenu[v] = struct{}{}
		}
		// role.Menu = append(role.Menu, menu[k]...)
	}

	role.Rules = make([]string, 0, len(roleMenu))
	for k := range roleMenu {
		role.Rules = append(role.Rules, k)
	}

	return role, nil
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
		LEFT JOIN LATERAL (SELECT sso_id from %s WHERE id=a.user_id) AS u ON true
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
	query := fmt.Sprintf(`INSERT INTO %s(id, name, level, extends, description) VALUES ($1, $2, $3, $4, $5)`, RoleTable)
	id := uuid.New()

	_, err := r.db.ExecContext(ctx, query, id, role.Name, role.Level, pq.Array(role.Extends))
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RoleRepo) Update(ctx context.Context, role *models.RoleDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, level=$2, extends=$3, description=$4 WHERE id=$5`, RoleTable)

	_, err := r.db.ExecContext(ctx, query, role.Name, role.Level, pq.Array(role.Extends), role.Description, role.ID)
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
