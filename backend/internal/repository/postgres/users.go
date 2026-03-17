package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepo struct {
	db *sqlx.DB
	Transaction
}

func NewUserRepo(db *sqlx.DB, transaction Transaction) *UserRepo {
	return &UserRepo{
		db:          db,
		Transaction: transaction,
	}
}

type User interface {
	GetAll(ctx context.Context) ([]*models.UserData, error)
	GetByAccess(ctx context.Context, req *models.GetByAccessDTO) ([]*models.UserData, error)
	GetByRealm(ctx context.Context, req *models.GetByRealmDTO) ([]*models.UserData, error)
	GetById(ctx context.Context, id string) (*models.UserData, error)
	GetBySSOId(ctx context.Context, id string) (*models.UserData, error)
	Create(ctx context.Context, dto *models.UserData) error
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.UserData) error
	Update(ctx context.Context, dto *models.UserData) error
	UpdateSeveral(ctx context.Context, tx Tx, dto []*models.UserData) error
	Delete(ctx context.Context, id string) error
	DeleteSeveral(ctx context.Context, tx Tx, ids []string) error
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
	query := fmt.Sprintf(`SELECT u.id, u.sso_id, username, first_name, last_name, email
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

func (r *UserRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.UserData) error {
	if len(dto) == 0 {
		return nil
	}

	n := len(dto)
	ids := make([]string, n)
	usernames := make([]string, n)
	emails := make([]string, n)
	ssoIds := make([]string, n)
	firstNames := make([]string, n)
	lastNames := make([]string, n)

	for i, d := range dto {
		d.ID = uuid.NewString()

		ids[i] = d.ID
		usernames[i] = d.Username
		emails[i] = d.Email
		ssoIds[i] = d.SSO_ID
		firstNames[i] = d.FirstName
		lastNames[i] = d.LastName
	}

	query := fmt.Sprintf(`
        INSERT INTO %s (id, username, email, sso_id, first_name, last_name)
        SELECT * FROM UNNEST(
            $1::uuid[], 
            $2::text[], 
            $3::text[], 
            $4::text[], 
            $5::text[], 
            $6::text[]
        )`, UserTable)

	_, err := r.getExec(tx).ExecContext(ctx, query,
		pq.Array(ids),
		pq.Array(usernames),
		pq.Array(emails),
		pq.Array(ssoIds),
		pq.Array(firstNames),
		pq.Array(lastNames),
	)
	if err != nil {
		return fmt.Errorf("failed bulk insert into %s: %w", UserTable, err)
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

func (r *UserRepo) UpdateSeveral(ctx context.Context, tx Tx, dto []*models.UserData) error {
	if len(dto) == 0 {
		return nil
	}

	n := len(dto)
	usernames := make([]string, n)
	emails := make([]string, n)
	ssoIds := make([]string, n)
	firstNames := make([]string, n)
	lastNames := make([]string, n)

	for i, v := range dto {
		usernames[i] = v.Username
		emails[i] = v.Email
		ssoIds[i] = v.SSO_ID
		firstNames[i] = v.FirstName
		lastNames[i] = v.LastName
	}

	query := fmt.Sprintf(`
        UPDATE %s AS t
        SET 
            username = s.username,
            email = s.email,
            first_name = s.first_name,
            last_name = s.last_name
        FROM (
            SELECT * FROM UNNEST(
                $1::text[], 
                $2::text[], 
                $3::text[], 
                $4::text[], 
                $5::text[]
            ) AS s(username, email, sso_id, first_name, last_name)
        ) AS s
        WHERE t.sso_id = s.sso_id`,
		UserTable,
	)

	_, err := r.getExec(tx).ExecContext(ctx, query,
		pq.Array(usernames),
		pq.Array(emails),
		pq.Array(ssoIds),
		pq.Array(firstNames),
		pq.Array(lastNames),
	)
	if err != nil {
		return fmt.Errorf("failed to execute bulk update: %w", err)
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

func (r *UserRepo) DeleteSeveral(ctx context.Context, tx Tx, ids []string) error {
	if len(ids) == 0 {
		return nil
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE id=ANY($1::uuid[])`, UserTable)

	if _, err := r.getExec(tx).ExecContext(ctx, query, pq.Array(ids)); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
