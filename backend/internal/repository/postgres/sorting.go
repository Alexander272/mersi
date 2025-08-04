package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres/pq_models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SortingRepo struct {
	db *sqlx.DB
}

func NewSortingRepo(db *sqlx.DB) *SortingRepo {
	return &SortingRepo{
		db: db,
	}
}

type Sorting interface {
	Get(ctx context.Context, req *models.GetSortingDTO) (models.SortingMap, error)
	Create(ctx context.Context, dto *models.SortingDTO) error
	CreateSeveral(ctx context.Context, dto []*models.SortingDTO) error
	Update(ctx context.Context, dto *models.SortingDTO) error
	Delete(ctx context.Context, dto *models.DeleteSortingDTO) error
	DeleteAll(ctx context.Context, dto *models.DeleteSortingDTO) error
}

func (r *SortingRepo) Get(ctx context.Context, req *models.GetSortingDTO) (models.SortingMap, error) {
	query := fmt.Sprintf(`SELECT id, name, order_type FROM %s WHERE user_id=$1 AND section_id=$2 ORDER BY created_at`, SortingTable)

	tmp := []*pq_models.Sorting{}
	if err := r.db.SelectContext(ctx, &tmp, query, req.UserId, req.SectionId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}

	data := make(map[string]string, 0)
	for _, v := range tmp {
		data[v.Name] = v.OrderType
	}
	return data, nil
}

func (r *SortingRepo) Create(ctx context.Context, dto *models.SortingDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, section_id, name, order_type) 
		VALUES (:id, :user_id, :section_id, :name, :order_type)`,
		SortingTable,
	)
	dto.Id = uuid.NewString()

	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SortingRepo) CreateSeveral(ctx context.Context, dto []*models.SortingDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, user_id, section_id, name, order_type) 
		VALUES (:id, :user_id, :section_id, :name, :order_type)`,
		SortingTable,
	)
	for i := range dto {
		dto[i].Id = uuid.NewString()
	}

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SortingRepo) Update(ctx context.Context, dto *models.SortingDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET order_type=:order_type WHERE name=:name AND user_id=:user_id AND section_id=:section_id`, SortingTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SortingRepo) Delete(ctx context.Context, dto *models.DeleteSortingDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE name=:name AND user_id=:user_id AND section_id=:section_id`, SortingTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *SortingRepo) DeleteAll(ctx context.Context, dto *models.DeleteSortingDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id=:user_id AND section_id=:section_id`, SortingTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
