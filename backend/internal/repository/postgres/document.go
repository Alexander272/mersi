package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/jmoiron/sqlx"
)

type DocumentRepo struct {
	db *sqlx.DB
}

func NewDocumentRepo(db *sqlx.DB) *DocumentRepo {
	return &DocumentRepo{
		db: db,
	}
}

type Document interface {
	GetTemp(ctx context.Context, req *models.GetTempDocumentDTO) ([]*models.Document, error)
	CreateSeveral(ctx context.Context, dto []*models.Document) error
	UpdatePath(ctx context.Context, dto *models.PathParts) (int64, error)
	Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error
}

func (r *DocumentRepo) GetTemp(ctx context.Context, req *models.GetTempDocumentDTO) ([]*models.Document, error) {
	query := fmt.Sprintf(`SELECT id, label, size, path, type, belongs, user_id FROM %s 
		WHERE user_id=$1 AND belongs=$2 AND path LIKE '%%temp/%%'`,
		DocumentsTable,
	)
	data := []*models.Document{}

	if err := r.db.SelectContext(ctx, &data, query, req.UserId, req.Group); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *DocumentRepo) CreateSeveral(ctx context.Context, dto []*models.Document) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, label, size, path, type, belongs, instrument_id, user_id) 
		VALUES (:id, :label, :size, :path, :type, :belongs, :instrument_id, :user_id)`,
		DocumentsTable,
	)
	tmp := []*pq_models.Document{}
	for i := range dto {
		// dto[i].Id = uuid.NewString()

		tmp = append(tmp, &pq_models.Document{
			Id:           dto[i].Id,
			Label:        dto[i].Label,
			Size:         dto[i].Size,
			Path:         dto[i].Path,
			DocumentType: dto[i].DocumentType,
			Group:        dto[i].Group,
			UserId:       dto[i].UserId,
			InstrumentId: nil,
		})
		if dto[i].InstrumentId != "" {
			tmp[len(tmp)-1].InstrumentId = dto[i].InstrumentId
		}
	}

	if _, err := r.db.NamedExecContext(ctx, query, tmp); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DocumentRepo) UpdatePath(ctx context.Context, dto *models.PathParts) (int64, error) {
	src := fmt.Sprintf("temp/%s/%s", dto.UserId, dto.Group)
	query := fmt.Sprintf(`UPDATE %s SET instrument_id=$1::uuid, path=REPLACE(path, $2, $3||'/'||$1::text)
		WHERE (instrument_id IS NULL OR instrument_id=$1) AND belongs=$3 AND user_id=$4`,
		DocumentsTable,
	)

	res, err := r.db.ExecContext(ctx, query, dto.InstrumentId, src, dto.Group, dto.UserId)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query. error: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get number of rows affected. error: %w", err)
	}
	return count, nil
}

func (r *DocumentRepo) Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, DocumentsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
