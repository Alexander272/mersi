package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type VerificationDocRepo struct {
	db *sqlx.DB
}

func NewVerificationDocRepo(db *sqlx.DB) *VerificationDocRepo {
	return &VerificationDocRepo{
		db: db,
	}
}

type VerificationDoc interface {
	Get(ctx context.Context, req *models.GetVerificationDocsDTO) ([]*models.VerificationDoc, error)
	GetGrouped(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error)
	CreateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error
	Update(ctx context.Context, dto *models.VerificationDocDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerificationDocDTO) error
}

func (r *VerificationDocRepo) Get(ctx context.Context, req *models.GetVerificationDocsDTO) ([]*models.VerificationDoc, error) {
	query := fmt.Sprintf(`SELECT id, name, doc_id FROM %s WHERE verification_id=$1 ORDER BY created_at DESC`, VerificationDocsTable)
	data := []*models.VerificationDoc{}

	if err := r.db.SelectContext(ctx, &data, query, req.VerificationId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *VerificationDocRepo) GetGrouped(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error) {
	query := fmt.Sprintf(`SELECT d.id, verification_id, name, doc_id FROM %s AS d
		INNER JOIN %s AS v ON v.id=verification_id WHERE instrument_id=$1 ORDER BY verification_id, created_at`,
		VerificationDocsTable, VerificationTable,
	)
	tmp := []*pq_models.VerificationDoc{}

	if err := r.db.SelectContext(ctx, &tmp, query, req.InstrumentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := make(map[string]*models.Groups, 0)
	for _, v := range tmp {
		t := &models.VerificationDoc{
			Id:    v.Id,
			Name:  v.Name,
			DocId: v.DocId,
		}

		_, exists := data[v.VerificationId]
		if exists {
			data[v.VerificationId].Docs = append(data[v.VerificationId].Docs, t)
		} else {
			data[v.VerificationId].Docs = []*models.VerificationDoc{t}
		}
	}
	return &models.GroupedVerificationDocs{Groups: data}, nil
}

func (r *VerificationDocRepo) CreateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, verification_id, name, doc_id) 
		VALUES (:id, :verification_id, :name, CAST(NULLIF(:doc_id, '') AS uuid))`,
		VerificationDocsTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationDocRepo) Update(ctx context.Context, dto *models.VerificationDocDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=:name, updated_at=now() WHERE id=:id`, VerificationDocsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationDocRepo) UpdateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error {
	values := []string{}
	args := []interface{}{}
	for i, v := range dto {
		tmp := []interface{}{v.Id, v.Name}
		args = append(args, tmp...)
		numbers := []string{}
		for j := range tmp {
			numbers = append(numbers, fmt.Sprintf("$%d", i*len(tmp)+j+1))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(numbers, ",")))
	}

	query := fmt.Sprintf(`UPDATE %s AS t SET name=s.name, updated_at=now()
		FROM (VALUES %s) AS s(id, name) WHERE t.id=s.id::uuid`,
		VerificationTable, strings.Join(values, ","),
	)

	if _, err := r.db.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *VerificationDocRepo) Delete(ctx context.Context, dto *models.DeleteVerificationDocDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, VerificationDocsTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
