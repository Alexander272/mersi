package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccessesRepo struct {
	db *sqlx.DB
}

func NewAccessesRepo(db *sqlx.DB) *AccessesRepo {
	return &AccessesRepo{
		db: db,
	}
}

type Accesses interface {
	Get(ctx context.Context, req *models.GetAccessesDTO) ([]*models.Accesses, error)
	Create(ctx context.Context, dto *models.AccessesDTO) error
	Update(ctx context.Context, dto *models.AccessesDTO) error
	Delete(ctx context.Context, dto *models.DeleteAccessesDTO) error
}

func (r *AccessesRepo) Get(ctx context.Context, req *models.GetAccessesDTO) ([]*models.Accesses, error) {
	query := fmt.Sprintf(`SELECT a.id, realm_id, user_id, sso_id, username, first_name, last_name, email, role_id, name, a.created_at 
		FROM %s AS a 
		INNER JOIN %s AS u ON a.user_id=u.id 
		INNER JOIN %s AS r ON r.id=a.role_id 
		WHERE realm_id=$1 ORDER BY name, last_name, first_name`,
		AccessTable, UsersTable, RoleTable,
	)
	tmp := []*pq_models.Accesses{}

	if err := r.db.SelectContext(ctx, &tmp, query, req.RealmID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := []*models.Accesses{}
	for _, v := range tmp {
		data = append(data, &models.Accesses{
			ID:      v.Id,
			RealmID: v.RealmId,
			Created: v.Created,
			User: &models.UserData{
				ID:        v.UserId,
				SSO_ID:    v.SSOId,
				Username:  v.Username,
				FirstName: v.FirstName,
				LastName:  v.LastName,
				Email:     v.Email,
			},
			Role: &models.Role{
				ID:   v.RoleId,
				Name: v.RoleName,
			},
		})
	}

	return data, nil
}

func (r *AccessesRepo) Create(ctx context.Context, dto *models.AccessesDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, realm_id, user_id, role_id) 
		VALUES (:id, :realm_id, :user_id, :role_id)`,
		AccessTable,
	)
	dto.ID = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AccessesRepo) Update(ctx context.Context, dto *models.AccessesDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET realm_id=:realm_id, user_id=:user_id, role_id=:role_id WHERE id=:id`,
		AccessTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *AccessesRepo) Delete(ctx context.Context, dto *models.DeleteAccessesDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, AccessTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
