package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

type User interface {
	GetAll(ctx context.Context) ([]*models.UserData, error)
	GetByAccess(ctx context.Context, req *models.GetByAccessDTO) ([]*models.UserData, error)
	GetByRealm(ctx context.Context, req *models.GetByRealmDTO) ([]*models.UserData, error)
	GetById(ctx context.Context, id string) (*models.UserData, error)
	GetBySSOId(ctx context.Context, id string) (*models.UserData, error)
	Create(ctx context.Context, dto *models.UserData) error
	CreateSeveral(ctx context.Context, dto []*models.UserData) error
	Update(ctx context.Context, dto *models.UserData) error
	UpdateSeveral(ctx context.Context, dto []*models.UserData) error
	Delete(ctx context.Context, id string) error
	DeleteSeveral(ctx context.Context, ids []string) error
}

func (r *UserRepo) GetAll(ctx context.Context) ([]*models.UserData, error) {
	query := fmt.Sprintf(`SELECT id, username, email, sso_id, first_name, last_name FROM %s ORDER BY last_name`, UserTable)
	data := []*models.UserData{}

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *UserRepo) GetByAccess(ctx context.Context, req *models.GetByAccessDTO) ([]*models.UserData, error) {
	query := fmt.Sprintf(`SELECT u.id, sso_id, username, first_name, last_name, email
		FROM %s AS u
		INNER JOIN %s AS a ON user_id=u.id
		INNER JOIN %s AS r ON role_id=r.id
		WHERE realm_id=$1 AND name=$2 ORDER BY last_name, first_name`,
		UserTable, AccessTable, RoleTable,
	)
	data := []*models.UserData{}

	if err := r.db.SelectContext(ctx, &data, query, req.RealmID, req.Role); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *UserRepo) GetByRealm(ctx context.Context, req *models.GetByRealmDTO) ([]*models.UserData, error) {
	include := ""
	if req.Include {
		include = " NOT"
	}

	query := fmt.Sprintf(`SELECT id, sso_id, username, first_name, last_name, email
		FROM %s AS u
		LEFT JOIN LATERAL (SELECT user_id FROM %s WHERE realm_id=$1 AND user_id=u.id) AS a ON true
		WHERE user_id IS%s NULL ORDER BY last_name, first_name`,
		UserTable, AccessTable, include,
	)
	data := []*models.UserData{}

	if err := r.db.SelectContext(ctx, &data, query, req.RealmID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *UserRepo) GetById(ctx context.Context, id string) (*models.UserData, error) {
	query := fmt.Sprintf(`SELECT id, username, email, sso_id, first_name, last_name FROM %s WHERE id=$1`, UserTable)
	data := &models.UserData{}

	if err := r.db.GetContext(ctx, data, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *UserRepo) GetBySSOId(ctx context.Context, id string) (*models.UserData, error) {
	query := fmt.Sprintf(`SELECT id, username, email, sso_id, first_name, last_name FROM %s WHERE sso_id=$1`, UserTable)
	data := &models.UserData{}

	if err := r.db.GetContext(ctx, data, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *UserRepo) Create(ctx context.Context, dto *models.UserData) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, username, email, sso_id, first_name, last_name) VALUES ($1, $2, $3, $4, $5, $6)`, UserTable)

	_, err := r.db.ExecContext(ctx, query, dto.ID, dto.Username, dto.Email, dto.SSO_ID, dto.FirstName, dto.LastName)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *UserRepo) CreateSeveral(ctx context.Context, dto []*models.UserData) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, username, email, sso_id, first_name, last_name) VALUES `, UserTable)

	args := make([]interface{}, 0)
	values := make([]string, 0, len(dto))
	c := 6
	for i, d := range dto {
		id := uuid.New()
		args = append(args, id, d.Username, d.Email, d.SSO_ID, d.FirstName, d.LastName)
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5, i*c+6))
	}
	query += strings.Join(values, ",")

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *UserRepo) Update(ctx context.Context, dto *models.UserData) error {
	query := fmt.Sprintf(`UPDATE %s SET username=$1, email=$2, sso_id=$3, first_name=$4, last_name=$5 WHERE id=$6`, UserTable)

	_, err := r.db.ExecContext(ctx, query, dto.Username, dto.Email, dto.SSO_ID, dto.FirstName, dto.LastName, dto.ID)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *UserRepo) UpdateSeveral(ctx context.Context, dto []*models.UserData) error {
	values := []string{}
	args := []interface{}{}
	c := 5
	for i, v := range dto {
		args = append(args, v.Username, v.Email, v.SSO_ID, v.FirstName, v.LastName)
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*c+1, i*c+2, i*c+3, i*c+4, i*c+5))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET username=s.username, email=s.email, first_name=s.first_name, last_name=s.last_name 
		FROM (VALUES %s) AS s(username, email, sso_id, first_name, last_name) WHERE t.sso_id=s.sso_id`,
		UserTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *UserRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, UserTable)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *UserRepo) DeleteSeveral(ctx context.Context, ids []string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=ANY($1)`, UserTable)

	if _, err := r.db.ExecContext(ctx, query, pq.Array(ids)); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
