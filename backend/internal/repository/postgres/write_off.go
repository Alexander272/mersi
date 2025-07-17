package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type WriteOffRepo struct {
	db *sqlx.DB
}

func NewWriteOffRepo(db *sqlx.DB) *WriteOffRepo {
	return &WriteOffRepo{
		db: db,
	}
}

type WriteOff interface {
	Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error)
	Create(ctx context.Context, dto *models.WriteOffDTO) error
	Update(ctx context.Context, dto *models.WriteOffDTO) error
	Delete(ctx context.Context, dto *models.DeleteWriteOffDTO) error
}

func (r *WriteOffRepo) Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error) {
	query := fmt.Sprintf(`SELECT id, instrument_id, date, notes, doc_id, doc_name, created_at FROM %s 
		WHERE instrument_id=$1 ORDER BY date DESC`,
		WriteOffTable,
	)
	data := []*models.WriteOff{}

	if err := r.db.SelectContext(ctx, &data, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *WriteOffRepo) Create(ctx context.Context, dto *models.WriteOffDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, instrument_id, date, notes, doc_id, doc_name) 
		VALUES (:id, :instrument_id, :date, :notes, :doc_id, :doc_name)`,
		WriteOffTable,
	)
	dto.Id = uuid.NewString()
	if dto.DocId == "" {
		dto.DocId = uuid.Nil.String()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *WriteOffRepo) Update(ctx context.Context, dto *models.WriteOffDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET date=:date, notes=:notes, doc_id=:doc_id, doc_name=:doc_name WHERE id=:id`,
		WriteOffTable,
	)
	if dto.DocId == "" {
		dto.DocId = uuid.Nil.String()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *WriteOffRepo) Delete(ctx context.Context, dto *models.DeleteWriteOffDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, WriteOffTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
