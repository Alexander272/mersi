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

type VerificationRepo struct {
	db *sqlx.DB
	Transaction
}

func NewVerificationRepo(db *sqlx.DB, transaction Transaction) *VerificationRepo {
	return &VerificationRepo{
		db:          db,
		Transaction: transaction,
	}
}

type Verification interface {
	Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error)
	GetById(ctx context.Context, id string) (*models.Verification, error)
	GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error)
	CreateInTx(ctx context.Context, tx Tx, dto *models.VerificationDTO) error
	CreateSeveral(ctx context.Context, tx Tx, dto []*models.VerificationDTO) error
	Update(ctx context.Context, tx Tx, dto *models.VerificationDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteVerificationDTO) error
}

func (r *VerificationRepo) Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, next_date, register_link, not_verified, notes, status
		FROM %s WHERE instrument_id=$1 ORDER BY date DESC, created_at DESC`,
		VerificationTable,
	)
	data := []*models.Verification{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *VerificationRepo) GetById(ctx context.Context, id string) (*models.Verification, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, next_date, register_link, not_verified, notes, status 
		FROM %s WHERE id=$1`,
		VerificationTable,
	)
	verification := &models.Verification{}

	if err := r.db.GetContext(ctx, verification, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return verification, nil
}

func (r *VerificationRepo) GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, next_date, register_link, not_verified, notes, status
		FROM %s WHERE instrument_id=$1 ORDER BY date DESC, created_at DESC LIMIT 1`,
		VerificationTable,
	)
	verification := &models.Verification{}

	if err := r.db.GetContext(ctx, verification, query, req.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return verification, nil
}

func (r *VerificationRepo) CreateInTx(ctx context.Context, tx Tx, dto *models.VerificationDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, next_date, register_link, not_verified, notes, status)
		VALUES (:id, :instrument_id, :date, :next_date, :register_link, :not_verified, :notes, :status)`,
		VerificationTable,
	)
	dto.Id = uuid.NewString()

	if _, err := tx.TX().NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationRepo) Create(ctx context.Context, dto *models.VerificationDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, next_date, register_link, not_verified, notes, status)
		VALUES (:id, :instrument_id, :date, :next_date, :register_link, :not_verified, :notes, :status)`,
		VerificationTable,
	)
	dto.Id = uuid.NewString()

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationRepo) CreateSeveral(ctx context.Context, tx Tx, dto []*models.VerificationDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, next_date, register_link, not_verified, notes, status)
		VALUES (:id, :instrument_id, :date, :next_date, :register_link, :not_verified, :notes, :status)`,
		VerificationTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationRepo) Update(ctx context.Context, tx Tx, dto *models.VerificationDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date=:date, next_date=:next_date, register_link=:register_link, not_verified=:not_verified,
		notes=:notes, status=:status, updated_at=now() WHERE id=:id`,
		VerificationTable,
	)

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteVerificationDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, VerificationTable)

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
