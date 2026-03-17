package postgres

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type DepartmentAccessRepo struct {
	db *sqlx.DB
}

func NewDepartmentAccessRepo(db *sqlx.DB) *DepartmentAccessRepo {
	return &DepartmentAccessRepo{
		db: db,
	}
}

type DepartmentAccess interface {
	Get(ctx context.Context, req *models.GetDepartmentAccessDTO) ([]*models.DepartmentAccess, error)
	GetByUserId(ctx context.Context, req *models.GetDepartmentAccessDTO) ([]*models.DepartmentAccess, error)
	Replace(ctx context.Context, dto *models.ReplaceDepartmentAccessDTO) error
	Create(ctx context.Context, dto *models.DepartmentAccessDTO) error
	CreateSeveral(ctx context.Context, dto []*models.DepartmentAccessDTO) error
	Update(ctx context.Context, dto *models.DepartmentAccessDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.DepartmentAccessDTO) error
	Delete(ctx context.Context, dto *models.DeleteDepartmentAccessDTO) error
}

func (r *DepartmentAccessRepo) Get(ctx context.Context, req *models.GetDepartmentAccessDTO) ([]*models.DepartmentAccess, error) {
	query := fmt.Sprintf(`SELECT id, department_id, sso_id FROM %s WHERE department_id=$1`, DepartmentAccessTable)

	data := []*models.DepartmentAccess{}
	if err := r.db.SelectContext(ctx, &data, query, req.DepartmentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *DepartmentAccessRepo) GetByUserId(ctx context.Context, req *models.GetDepartmentAccessDTO) ([]*models.DepartmentAccess, error) {
	query := fmt.Sprintf(`SELECT id, department_id, sso_id FROM %s WHERE sso_id=$1`, DepartmentAccessTable)

	data := []*models.DepartmentAccess{}
	if err := r.db.SelectContext(ctx, &data, query, req.UserId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *DepartmentAccessRepo) Replace(ctx context.Context, dto *models.ReplaceDepartmentAccessDTO) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Удаляем старые записи
	deleteQuery := fmt.Sprintf(`DELETE FROM %s WHERE department_id=:department_id`, DepartmentAccessTable)
	if _, err = tx.NamedExecContext(ctx, deleteQuery, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	if len(dto.UserIDs) == 0 {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction. error: %w", err)
		}
		return nil
	}

	// 3. Массовая вставка через unnest
	query := fmt.Sprintf(`INSERT INTO %s (id, department_id, sso_id)
        SELECT gen_random_uuid(), $1, unnest($2::text[])`,
		DepartmentAccessTable,
	)
	if _, err = tx.ExecContext(ctx, query, dto.DepartmentId, pq.Array(dto.UserIDs)); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. error: %w", err)
	}
	return nil
}

func (r *DepartmentAccessRepo) Create(ctx context.Context, dto *models.DepartmentAccessDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s (id, department_id, sso_id) VALUES (:id, :department_id, :sso_id)`, DepartmentAccessTable)
	dto.Id = uuid.NewString()

	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DepartmentAccessRepo) CreateSeveral(ctx context.Context, dto []*models.DepartmentAccessDTO) error {
	if len(dto) == 0 {
		return nil
	}

	ids := make([]string, len(dto))
	deptIds := make([]string, len(dto))
	userIds := make([]string, len(dto))

	for i := range dto {
		dto[i].Id = uuid.NewString()

		ids[i] = dto[i].Id
		deptIds[i] = dto[i].DepartmentId
		userIds[i] = dto[i].UserId
	}

	query := fmt.Sprintf(`INSERT INTO %s (id, department_id, sso_id)
		SELECT * FROM UNNEST($1::uuid[], $2::uuid[], $3::text[])
	`, DepartmentAccessTable)

	if _, err := r.db.ExecContext(ctx, query, pq.Array(ids), pq.Array(deptIds), pq.Array(userIds)); err != nil {
		return fmt.Errorf("bulk insert via unnest failed: %w", err)
	}

	return nil
}

func (r *DepartmentAccessRepo) Update(ctx context.Context, dto *models.DepartmentAccessDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET department_id=:department_id, sso_id=:sso_id WHERE id=:id`, DepartmentAccessTable)

	_, err := r.db.NamedExecContext(ctx, query, dto)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DepartmentAccessRepo) UpdateSeveral(ctx context.Context, dto []*models.DepartmentAccessDTO) error {
	if len(dto) == 0 {
		return nil
	}

	// Разделяем данные на отдельные слайсы для передачи в UNNEST
	ids := make([]string, len(dto))
	deps := make([]string, len(dto))
	users := make([]string, len(dto))

	for i, v := range dto {
		ids[i] = v.Id
		deps[i] = v.DepartmentId
		users[i] = v.UserId
	}

	query := fmt.Sprintf(`UPDATE %s AS t 
		SET department_id = s.department_id, sso_id = s.sso_id 
		FROM (
			SELECT UNNEST($1::uuid[]) as id, 
				   UNNEST($2::uuid[]) as department_id, 
				   UNNEST($3::text[]) as sso_id, 
		) AS s 
		WHERE t.id = s.id`,
		DepartmentAccessTable,
	)

	// Передаем слайсы как массивы PostgreSQL
	_, err := r.db.ExecContext(ctx, query,
		pq.Array(ids),
		pq.Array(deps),
		pq.Array(users),
	)

	if err != nil {
		return fmt.Errorf("failed to execute unnest update: %w", err)
	}
	return nil
}

func (r *DepartmentAccessRepo) Delete(ctx context.Context, dto *models.DeleteDepartmentAccessDTO) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=:id`, DepartmentAccessTable)

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
