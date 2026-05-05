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

type TransferToDepRepo struct {
	db *sqlx.DB
	Transaction
}

func NewTransferToDepRepo(db *sqlx.DB, transaction Transaction) *TransferToDepRepo {
	return &TransferToDepRepo{
		db:          db,
		Transaction: transaction,
	}
}

type TransferToDepartment interface {
	Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error)
	GetById(ctx context.Context, dto *models.GetTransferToDepDTO) (*models.TransferToDepartment, error)
	GetLast(ctx context.Context, tx Tx, req *models.GetTransferToDepDTO) (*models.TransferToDepartment, error)
	Create(ctx context.Context, tx Tx, dto *models.TransferToDepartmentDTO) error
	CreateSeveral(ctx context.Context, dto []*models.TransferToDepartmentDTO) error
	Update(ctx context.Context, tx Tx, dto *models.TransferToDepartmentDTO) error
	Delete(ctx context.Context, tx Tx, dto *models.DeleteTransferToDepDTO) error
}

func (r *TransferToDepRepo) Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s 
		WHERE instrument_id=$1 ORDER BY date DESC`,
		TransferToDepTable,
	)
	data := []*models.TransferToDepartment{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *TransferToDepRepo) GetById(ctx context.Context, dto *models.GetTransferToDepDTO) (*models.TransferToDepartment, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s 
		WHERE id=$1`,
		TransferToDepTable,
	)
	data := &models.TransferToDepartment{}

	if err := r.db.GetContext(ctx, data, query, dto.Id); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *TransferToDepRepo) GetLast(ctx context.Context, tx Tx, req *models.GetTransferToDepDTO) (*models.TransferToDepartment, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s
		WHERE instrument_id=$1 ORDER BY date DESC LIMIT 1`,
		TransferToDepTable,
	)
	data := &models.TransferToDepartment{}

	if err := r.getExec(tx).GetContext(ctx, data, query, req.InstrumentId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *TransferToDepRepo) Create(ctx context.Context, tx Tx, dto *models.TransferToDepartmentDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, notes, doc_id, doc_name) 
		VALUES (:id, :instrument_id, :date, :notes, :doc_id, :doc_name)`,
		TransferToDepTable,
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

func (r TransferToDepRepo) CreateSeveral(ctx context.Context, dto []*models.TransferToDepartmentDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, notes, doc_id, doc_name) 
		VALUES (:id, :instrument_id, :date, :notes, :doc_id, :doc_name)`,
		TransferToDepTable,
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

func (r *TransferToDepRepo) Update(ctx context.Context, tx Tx, dto *models.TransferToDepartmentDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date=:date, notes=:notes, doc_id=:doc_id, doc_name=:doc_name WHERE id=:id`,
		TransferToDepTable,
	)
	if dto.DocId == "" {
		dto.DocId = uuid.Nil.String()
	}

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *TransferToDepRepo) Delete(ctx context.Context, tx Tx, dto *models.DeleteTransferToDepDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, TransferToDepTable)

	if _, err := r.getExec(tx).NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
