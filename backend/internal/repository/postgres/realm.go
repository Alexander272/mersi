package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type RealmRepo struct {
	db *sqlx.DB
}

func NewRealmRepo(db *sqlx.DB) *RealmRepo {
	return &RealmRepo{
		db: db,
	}
}

type Realm interface {
	Get(ctx context.Context, req *models.GetRealmsDTO) ([]*models.Realm, error)
	GetById(ctx context.Context, req *models.GetRealmByIdDTO) (*models.Realm, error)
	GetByUser(ctx context.Context, req *models.GetRealmByUserDTO) ([]*models.Realm, error)
	Create(ctx context.Context, dto *models.RealmDTO) error
	Update(ctx context.Context, dto *models.RealmDTO) error
	Delete(ctx context.Context, dto *models.DeleteRealmDTO) error
}

func (r *RealmRepo) Get(ctx context.Context, req *models.GetRealmsDTO) ([]*models.Realm, error) {
	condition := "WHERE is_active=true"
	if req.All {
		condition = ""
	}

	query := fmt.Sprintf(`SELECT id, name, realm, is_active, notification_channel, expiration_notice, location_type, has_responsible,
		has_employees, need_confirmed, has_commissioning_cert, has_preservations, has_transfer, created_at
		FROM %s %s ORDER BY created_at`,
		RealmTable, condition,
	)
	data := []*models.Realm{}

	if err := r.db.SelectContext(ctx, &data, query); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RealmRepo) GetByUser(ctx context.Context, req *models.GetRealmByUserDTO) ([]*models.Realm, error) {
	query := fmt.Sprintf(`SELECT r.id, name, realm, notification_channel, expiration_notice, location_type, is_active, has_responsible,
		has_employees, need_confirmed, has_commissioning_cert, has_preservations, has_transfer, created_at
		FROM %s AS r
		LEFT JOIN LATERAL (SELECT a.id FROM %s AS a INNER JOIN %s AS u ON a.user_id=u.id WHERE sso_id=$1 AND realm_id=r.id) AS a ON true
		WHERE a.id IS NOT NULL ORDER BY created_at`,
		RealmTable, AccessTable, UsersTable,
	)
	data := []*models.Realm{}

	if err := r.db.SelectContext(ctx, &data, query, req.UserID); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RealmRepo) GetById(ctx context.Context, req *models.GetRealmByIdDTO) (*models.Realm, error) {
	query := fmt.Sprintf(`SELECT id, name, realm, is_active, notification_channel, expiration_notice, location_type,
		has_responsible, has_employees, need_confirmed, has_commissioning_cert, has_preservations, has_transfer, created_at
		FROM %s WHERE id=$1`,
		RealmTable,
	)
	data := &models.Realm{}

	if err := r.db.GetContext(ctx, data, query, req.ID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *RealmRepo) Create(ctx context.Context, dto *models.RealmDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, name, realm, is_active, notification_channel, expiration_notice, location_type,
		has_responsible, has_employees, need_confirmed, has_commissioning_cert, has_preservations, has_transfer)
		VALUES (:id, :name, :realm, :is_active, :notification_channel, :expiration_notice, :location_type, :has_responsible,
		:has_employees, :need_confirmed, :has_commissioning_cert, :has_preservations, :has_transfer)`,
		RealmTable,
	)
	dto.ID = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RealmRepo) Update(ctx context.Context, dto *models.RealmDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, is_active=:is_active, notification_channel=:notification_channel,
		expiration_notice=:expiration_notice, location_type=:location_type, has_responsible=:has_responsible,
		has_employees=:has_employees, need_confirmed=:need_confirmed, has_commissioning_cert=:has_commissioning_cert,
		has_preservations=:has_preservations, has_transfer=:has_transfer WHERE id=:id`,
		RealmTable,
	)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *RealmRepo) Delete(ctx context.Context, dto *models.DeleteRealmDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, RealmTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
