package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WriteOffRepo struct {
	db *sqlx.DB
	Transaction
}

func NewWriteOffRepo(db *sqlx.DB, transaction Transaction) *WriteOffRepo {
	return &WriteOffRepo{
		db:          db,
		Transaction: transaction,
	}
}

type WriteOff interface {
	Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error)
	GetById(ctx context.Context, dto *models.GetWriteOffDTO) (*models.WriteOff, error)
	GetByInstrumentAndDate(ctx context.Context, tx Tx, instrumentId string, date time.Time) (*models.WriteOff, error)
	GetLast(ctx context.Context, tx Tx, req *models.GetWriteOffDTO) (*models.WriteOff, error)
	Create(ctx context.Context, tx Tx, dto *models.WriteOffDTO) error
	CreateSeveral(ctx context.Context, dto []*models.WriteOffDTO) error
	Update(ctx context.Context, tx Tx, dto *models.WriteOffDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteWriteOffDTO) error
}

func (r *WriteOffRepo) Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s 
		WHERE instrument_id=$1 ORDER BY date DESC, created_at DESC`,
		WriteOffTable,
	)
	data := []*models.WriteOff{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *WriteOffRepo) GetById(ctx context.Context, dto *models.GetWriteOffDTO) (*models.WriteOff, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s 
		WHERE id=$1`,
		WriteOffTable,
	)
	data := &models.WriteOff{}

	if err := r.db.GetContext(ctx, data, query, dto.Id); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *WriteOffRepo) GetByInstrumentAndDate(ctx context.Context, tx Tx, instrumentId string, date time.Time) (*models.WriteOff, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s 
		WHERE instrument_id=$1 AND date=$2 LIMIT 1`,
		WriteOffTable,
	)
	data := &models.WriteOff{}

	if err := r.getExec(tx).GetContext(ctx, data, query, instrumentId, date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *WriteOffRepo) GetLast(ctx context.Context, tx Tx, req *models.GetWriteOffDTO) (*models.WriteOff, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s
		WHERE instrument_id=$1 ORDER BY date DESC LIMIT 1`,
		WriteOffTable,
	)
	data := &models.WriteOff{}

	if err := r.getExec(tx).GetContext(ctx, data, query, req.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *WriteOffRepo) Create(ctx context.Context, tx Tx, dto *models.WriteOffDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, notes, doc_id, doc_name) 
		VALUES (:id, :instrument_id, :date, :notes, :doc_id, :doc_name)`,
		WriteOffTable,
	)
	dto.Id = uuid.NewString()
	if dto.DocId == "" {
		dto.DocId = uuid.Nil.String()
	}

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *WriteOffRepo) CreateSeveral(ctx context.Context, dto []*models.WriteOffDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, notes, doc_id, doc_name) 
		VALUES (:id, :instrument_id, :date, :notes, :doc_id, :doc_name)`,
		WriteOffTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
		if dto[i].DocId == "" {
			dto[i].DocId = uuid.Nil.String()
		}
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *WriteOffRepo) Update(ctx context.Context, tx Tx, dto *models.WriteOffDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date=:date, notes=:notes, doc_id=:doc_id, doc_name=:doc_name WHERE id=:id`,
		WriteOffTable,
	)
	if dto.DocId == "" {
		dto.DocId = uuid.Nil.String()
	}

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *WriteOffRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteWriteOffDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, WriteOffTable)

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
