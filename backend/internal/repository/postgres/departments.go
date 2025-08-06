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

type DepartmentRepo struct {
	db *sqlx.DB
}

func NewDepartmentRepo(db *sqlx.DB) *DepartmentRepo {
	return &DepartmentRepo{
		db: db,
	}
}

type Department interface {
	GetAll(ctx context.Context, req *models.GetDepartmentsDTO) ([]*models.Department, error)
	GetById(ctx context.Context, req *models.GetDepartmentByIdDTO) (*models.Department, error)
	GetBySSOId(context.Context, string) ([]*models.Department, error)
	Create(context.Context, *models.DepartmentDTO) (string, error)
	Update(context.Context, *models.DepartmentDTO) error
	Delete(context.Context, string) error
}

func (r *DepartmentRepo) GetAll(ctx context.Context, req *models.GetDepartmentsDTO) ([]*models.Department, error) {
	query := fmt.Sprintf(`SELECT id, name FROM %s WHERE realm_id=$1 ORDER BY name`, DepartmentTable)
	departments := []*models.Department{}

	if err := r.db.SelectContext(ctx, &departments, query, req.RealmId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return departments, nil
}

func (r *DepartmentRepo) GetById(ctx context.Context, req *models.GetDepartmentByIdDTO) (*models.Department, error) {
	query := fmt.Sprintf(`SELECT d.id, d.name, COALESCE(d.leader_id::text, '') AS leader_id, COALESCE(d.channel_id::text, '') AS channel_id, 
		COALESCE(c.name, '') AS channel_name FROM %s AS d LEFT JOIN %s AS c ON d.channel_id=c.id WHERE d.id=$1`,
		DepartmentTable, ChannelTable,
	)
	data := &models.Department{}

	if err := r.db.GetContext(ctx, data, query, req.Id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *DepartmentRepo) GetBySSOId(ctx context.Context, id string) ([]*models.Department, error) {
	query := fmt.Sprintf(`SELECT d.id, d.name FROM %s AS d INNER JOIN %s AS e ON e.department_id=d.id WHERE e.is_lead=true AND sso_id=$1`,
		DepartmentTable, EmployeeTable,
	)
	data := []*models.Department{}

	if err := r.db.SelectContext(ctx, &data, query, id); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *DepartmentRepo) Create(ctx context.Context, dto *models.DepartmentDTO) (string, error) {
	query := fmt.Sprintf(`INSERT INTO %s(id, name, leader_id, channel_id, realm_id) VALUES ($1, $2, $3, $4, $5)`, DepartmentTable)
	dto.Id = uuid.NewString()

	var leaderId *string = &dto.LeaderId
	if dto.LeaderId == "" {
		leaderId = nil
	}
	var channelId *string = &dto.ChannelId
	if dto.ChannelId == "" {
		channelId = nil
	}

	_, err := r.db.ExecContext(ctx, query, dto.Id, dto.Name, leaderId, channelId, dto.RealmId)
	if err != nil {
		return "", fmt.Errorf("failed to execute query. error: %w", err)
	}
	return dto.Id, nil
}

func (r *DepartmentRepo) Update(ctx context.Context, dto *models.DepartmentDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, leader_id=$2, channel_id=$3 WHERE id=$4`, DepartmentTable)

	var leaderId *string = &dto.LeaderId
	if dto.LeaderId == "" {
		leaderId = nil
	}
	var channelId *string = &dto.ChannelId
	if dto.ChannelId == "" {
		channelId = nil
	}

	_, err := r.db.ExecContext(ctx, query, dto.Name, leaderId, channelId, dto.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *DepartmentRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, DepartmentTable)

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}
