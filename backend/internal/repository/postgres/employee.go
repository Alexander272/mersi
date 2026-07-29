package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

type EmployeeRepo struct {
	db *sqlx.DB
}

func NewEmployeeRepo(db *sqlx.DB) *EmployeeRepo {
	return &EmployeeRepo{
		db: db,
	}
}

type Employee interface {
	GetAll(context.Context, *models.GetEmployeesDTO) ([]*models.Employee, error)
	GetUnique(context.Context, *models.GetUniqueEmployeeDTO) ([]*models.Employee, error)
	GetByDepartment(context.Context, string) ([]*models.Employee, error)
	GetByName(ctx context.Context, req *models.GetEmployeeByNameDTO) (*models.Employee, error)
	GetById(context.Context, string) (*models.Employee, error)
	GetByMostId(context.Context, string) (*models.EmployeeData, error)
	GetBySSOId(context.Context, string) (*models.Employee, error)
	Create(context.Context, *models.WriteEmployeeDTO) error
	Update(context.Context, *models.WriteEmployeeDTO) error
	Delete(context.Context, string) error
}

func (r *EmployeeRepo) GetAll(ctx context.Context, req *models.GetEmployeesDTO) ([]*models.Employee, error) {
	query := fmt.Sprintf(`SELECT e.id, e.name, notes, department_id, most_id, is_lead FROM %s AS e
		INNER JOIN %s AS d ON e.department_id=d.id`,
		EmployeeTable, DepartmentTable,
	)

	params := []interface{}{}
	filters := []string{}

	i := 0
	for k, v := range req.Filters {
		i++
		filters = append(filters, fmt.Sprintf("%s=$%d", k, i))
		params = append(params, v)
	}
	filter := ""
	if len(filters) > 0 {
		filter = fmt.Sprintf(" WHERE %s", strings.Join(filters, " AND "))
	}
	sort := " ORDER BY department_id, name"

	query += filter + sort

	employees := []*models.Employee{}
	if err := r.db.SelectContext(ctx, &employees, query, params...); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return employees, nil
}

func (r *EmployeeRepo) GetUnique(ctx context.Context, dto *models.GetUniqueEmployeeDTO) ([]*models.Employee, error) {
	query := fmt.Sprintf(`SELECT  DISTINCT(e.name) FROM %s AS e INNER JOIN %s AS d ON e.department_id=d.id
		WHERE realm_id=$1 ORDER BY e.name`,
		EmployeeTable, DepartmentTable,
	)
	employees := []*models.Employee{}

	if err := r.db.SelectContext(ctx, &employees, query, dto.Realm); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return employees, nil
}

func (r *EmployeeRepo) GetByDepartment(ctx context.Context, departmentId string) ([]*models.Employee, error) {
	query := fmt.Sprintf(`SELECT id, name, most_id FROM %s WHERE department_id=$1 ORDER BY name`, EmployeeTable)

	employees := []*models.Employee{}
	if err := r.db.SelectContext(ctx, &employees, query, departmentId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return employees, nil
}

func (r *EmployeeRepo) GetByName(ctx context.Context, req *models.GetEmployeeByNameDTO) (*models.Employee, error) {
	query := fmt.Sprintf(`SELECT id, name, most_id FROM %s WHERE name=$1 AND department_id=2`, EmployeeTable)
	data := &models.Employee{}

	if err := r.db.GetContext(ctx, data, query, req.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return data, nil
}

func (r *EmployeeRepo) GetById(ctx context.Context, id string) (*models.Employee, error) {
	query := fmt.Sprintf(`SELECT id, name, department_id, notes, most_id FROM %s WHERE id=$1`, EmployeeTable)

	employee := &models.Employee{}
	if err := r.db.GetContext(ctx, employee, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return employee, nil
}

func (r *EmployeeRepo) GetByMostId(ctx context.Context, mostId string) (*models.EmployeeData, error) {
	query := fmt.Sprintf(`SELECT e.id, e.name, d.name AS department, most_id, d.leader_id=e.id AS is_lead
		FROM %s AS e INNER JOIN %s AS d ON department_id=d.id WHERE most_id=$1`,
		EmployeeTable, DepartmentTable,
	)
	employee := &models.EmployeeData{}

	if err := r.db.GetContext(ctx, employee, query, mostId); err != nil {
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return employee, nil
}

func (r *EmployeeRepo) GetBySSOId(ctx context.Context, id string) (*models.Employee, error) {
	//TODO тут может быть несколько строк и тогда запрос выдаст ошибку
	query := fmt.Sprintf(`SELECT id, name, department_id, most_id, is_lead FROM %s WHERE sso_id=$1`, EmployeeTable)

	employee := &models.Employee{}
	if err := r.db.GetContext(ctx, employee, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRows
		}
		return nil, fmt.Errorf("failed to execute query. error: %w", err)
	}
	return employee, nil
}

func (r *EmployeeRepo) Create(ctx context.Context, employee *models.WriteEmployeeDTO) error {
	query := fmt.Sprintf(`INSERT INTO %s(id, name, notes, department_id, most_id) VALUES ($1, $2, $3, $4, $5)`, EmployeeTable)
	employee.Id = uuid.NewString()

	_, err := r.db.ExecContext(ctx, query, employee.Id, employee.Name, employee.Notes, employee.DepartmentId, employee.MattermostId)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	return nil
}

func (r *EmployeeRepo) Update(ctx context.Context, employee *models.WriteEmployeeDTO) error {
	query := fmt.Sprintf(`UPDATE %s SET name=$1, department_id=$2, notes=$3 WHERE id=$4`, EmployeeTable)

	res, err := r.db.ExecContext(ctx, query, employee.Name, employee.DepartmentId, employee.Notes, employee.Id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}

func (r *EmployeeRepo) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id=$1", EmployeeTable)

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to execute query. error: %w", err)
	}
	if err := checkRowsAffected(res); err != nil {
		return err
	}
	return nil
}
