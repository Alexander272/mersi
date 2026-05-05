package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type ActivityLogRepo struct {
	db *sqlx.DB
}

func NewActivityLogRepo(db *sqlx.DB) *ActivityLogRepo {
	return &ActivityLogRepo{db: db}
}

type ActivityLog interface {
	Create(ctx context.Context, dto *models.CreateActivityLogDTO) error
	CreateSeveral(ctx context.Context, dto []*models.CreateActivityLogDTO) error
	GetByRecord(ctx context.Context, dto *models.GetActivityLogDTO) ([]*models.ActivityLog, error)
	GetAll(ctx context.Context, dto *models.ActivityLogFilter) ([]*models.ActivityLog, error)
}

func (r *ActivityLogRepo) Create(ctx context.Context, dto *models.CreateActivityLogDTO) error {
	query := `INSERT INTO activity_log (table_name, record_id, record_name, action, field_name, old_value, new_value, user_id, user_name) 
		VALUES (:table_name, :record_id, :record_name, :action, :field_name, :old_value, :new_value, :user_id, :user_name)`

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to create activity log record. error: %w", err)
	}
	return nil
}

func (r *ActivityLogRepo) CreateSeveral(ctx context.Context, dto []*models.CreateActivityLogDTO) error {
	if len(dto) == 0 {
		return nil
	}

	query := `INSERT INTO activity_log (table_name, record_id, record_name, action, field_name, old_value, new_value, user_id, user_name) 
		VALUES (:table_name, :record_id, :record_name, :action, :field_name, :old_value, :new_value, :user_id, :user_name)`

	if _, err := r.db.NamedExecContext(ctx, query, dto); err != nil {
		return fmt.Errorf("failed to create several activity log records. error: %w", err)
	}
	return nil
}

func (r *ActivityLogRepo) GetByRecord(ctx context.Context, dto *models.GetActivityLogDTO) ([]*models.ActivityLog, error) {
	query := `SELECT id, table_name, record_id, record_name, action, field_name, old_value, new_value, user_id, user_name, created_at 
		FROM activity_log WHERE table_name=$1 AND record_id=$2 ORDER BY created_at DESC`

	if dto.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", dto.Limit)
	}
	if dto.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", dto.Offset)
	}

	data := []*models.ActivityLog{}
	if err := r.db.SelectContext(ctx, &data, query, dto.TableName, dto.RecordId); err != nil {
		return nil, fmt.Errorf("failed to get activity log by record. error: %w", err)
	}
	return data, nil
}

func (r *ActivityLogRepo) GetAll(ctx context.Context, dto *models.ActivityLogFilter) ([]*models.ActivityLog, error) {
	conditions := []string{}
	args := []interface{}{}
	argIdx := 1

	if dto.TableName != "" {
		conditions = append(conditions, fmt.Sprintf("table_name=$%d", argIdx))
		args = append(args, dto.TableName)
		argIdx++
	}
	if dto.RecordId != "" {
		conditions = append(conditions, fmt.Sprintf("record_id=$%d", argIdx))
		args = append(args, dto.RecordId)
		argIdx++
	}
	if dto.UserId != "" {
		conditions = append(conditions, fmt.Sprintf("user_id=$%d", argIdx))
		args = append(args, dto.UserId)
		argIdx++
	}
	if dto.Action != "" {
		conditions = append(conditions, fmt.Sprintf("action=$%d", argIdx))
		args = append(args, dto.Action)
		argIdx++
	}
	if !dto.DateFrom.IsZero() {
		conditions = append(conditions, fmt.Sprintf("created_at>=$%d", argIdx))
		args = append(args, dto.DateFrom)
		argIdx++
	}
	if !dto.DateTo.IsZero() {
		conditions = append(conditions, fmt.Sprintf("created_at<=$%d", argIdx))
		args = append(args, dto.DateTo)
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	query := fmt.Sprintf(`SELECT id, table_name, record_id, record_name, action, field_name, old_value, new_value, user_id, user_name, created_at 
		FROM activity_log %s ORDER BY created_at DESC`, whereClause)

	if dto.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", dto.Limit)
	}
	if dto.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", dto.Offset)
	}

	data := []*models.ActivityLog{}
	if err := r.db.SelectContext(ctx, &data, query, args...); err != nil {
		return nil, fmt.Errorf("failed to get activity log. error: %w", err)
	}
	return data, nil
}
