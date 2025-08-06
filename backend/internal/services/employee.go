package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type EmployeeService struct {
	repo repository.Employee
	// location Location
}

func NewEmployeeService(repo repository.Employee) *EmployeeService {
	return &EmployeeService{
		repo: repo,
		// location: location,
	}
}

type Employee interface {
	GetAll(context.Context, *models.GetEmployeesDTO) ([]*models.Employee, error)
	GetUnique(context.Context) ([]*models.Employee, error)
	GetByDepartment(context.Context, string) ([]*models.Employee, error)
	GetById(context.Context, string) (*models.Employee, error)
	GetByMostId(context.Context, string) (*models.EmployeeData, error)
	GetBySSOId(context.Context, string) (*models.Employee, error)
	Create(context.Context, *models.WriteEmployeeDTO) error
	Update(context.Context, *models.WriteEmployeeDTO) error
	Delete(context.Context, string) error
}

func (s *EmployeeService) GetAll(ctx context.Context, req *models.GetEmployeesDTO) ([]*models.Employee, error) {
	employees, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees. error: %w", err)
	}
	if employees == nil {
		employees = []*models.Employee{}
	}
	return employees, nil
}

func (s *EmployeeService) GetUnique(ctx context.Context) ([]*models.Employee, error) {
	employees, err := s.repo.GetUnique(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique employees. error: %w", err)
	}
	if employees == nil {
		employees = []*models.Employee{}
	}
	return employees, nil
}

func (s *EmployeeService) GetBySSOId(ctx context.Context, id string) (*models.Employee, error) {
	employee, err := s.repo.GetBySSOId(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get employee by sso id. error: %w", err)
	}
	return employee, nil
}

func (s *EmployeeService) GetByDepartment(ctx context.Context, departmentId string) ([]*models.Employee, error) {
	users, err := s.repo.GetByDepartment(ctx, departmentId)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by department. error: %w", err)
	}
	return users, nil
}

// func (s *EmployeeService) GetDepartments(ctx context.Context, req *models.GetDepartmentsDTO) ([]string, error) {
// 	// employee, err := s.repo.GetBySSOId(ctx, req.UserId)
// 	// if err != nil {
// 	// 	if errors.Is(err, models.ErrNoRows) {
// 	// 		return nil, err
// 	// 	}
// 	// 	return nil, fmt.Errorf("failed to get employee by sso id. error: %w", err)
// 	// }
// 	// departments := []string{}
// 	// // employee
// 	// return employee, nil
// 	return nil, fmt.Errorf("not implemented")
// }

func (s *EmployeeService) GetById(ctx context.Context, id string) (*models.Employee, error) {
	employee, err := s.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get employee by id. error: %w", err)
	}
	return employee, nil
}

func (s *EmployeeService) GetByMostId(ctx context.Context, mostId string) (*models.EmployeeData, error) {
	employee, err := s.repo.GetByMostId(ctx, mostId)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee by most id. error: %w", err)
	}
	return employee, nil
}

func (s *EmployeeService) Create(ctx context.Context, employee *models.WriteEmployeeDTO) error {
	if err := s.repo.Create(ctx, employee); err != nil {
		return fmt.Errorf("failed to create employee. error: %w", err)
	}
	return nil
}

func (s *EmployeeService) Update(ctx context.Context, employee *models.WriteEmployeeDTO) error {
	if err := s.repo.Update(ctx, employee); err != nil {
		return fmt.Errorf("failed to update employee. error: %w", err)
	}

	//TODO
	// if err := s.location.UpdatePerson(ctx, &models.UpdatePlaceDTO{PersonId: employee.Id}); err != nil {
	// 	return err
	// }
	return nil
}

func (s *EmployeeService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete employee. error: %w", err)
	}
	return nil
}
