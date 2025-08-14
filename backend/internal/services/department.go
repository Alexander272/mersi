package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type DepartmentService struct {
	repo     repository.Department
	location Location
}

func NewDepartmentService(repo repository.Department, location Location) *DepartmentService {
	return &DepartmentService{
		repo:     repo,
		location: location,
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

func (s *DepartmentService) GetAll(ctx context.Context, req *models.GetDepartmentsDTO) ([]*models.Department, error) {
	departments, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all departments. error: %w", err)
	}
	if departments == nil {
		departments = []*models.Department{}
	}
	return departments, nil
}

func (s *DepartmentService) GetById(ctx context.Context, req *models.GetDepartmentByIdDTO) (*models.Department, error) {
	department, err := s.repo.GetById(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get department by id. error: %w", err)
	}
	return department, nil
}

func (s *DepartmentService) GetBySSOId(ctx context.Context, id string) ([]*models.Department, error) {
	departments, err := s.repo.GetBySSOId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments by sso id. error: %w", err)
	}
	return departments, nil
}

func (s *DepartmentService) Create(ctx context.Context, department *models.DepartmentDTO) (string, error) {
	id, err := s.repo.Create(ctx, department)
	if err != nil {
		return id, fmt.Errorf("failed to create department. error: %w", err)
	}
	return id, nil
}

func (s *DepartmentService) Update(ctx context.Context, department *models.DepartmentDTO) error {
	if err := s.repo.Update(ctx, department); err != nil {
		return fmt.Errorf("failed to update department. error: %w", err)
	}
	return nil
}

func (s *DepartmentService) Delete(ctx context.Context, id string) error {
	// надо запретить удалять подразделение если у него есть инструменты
	locs, err := s.location.GetUsedByDepartment(ctx, &models.GetLocationByDepartmentDTO{DepartmentId: id})
	if err != nil {
		return err
	}
	if len(locs) > 0 {
		return models.ErrHasInstrument
	}

	// при удалении подразделения надо записать его название в таблицу с перемещениями (чтобы оно сохранилось в истории)
	// и надо записать фио всех работников в этом подразделении тк они тоже удаляться
	if err := s.location.SetDepartment(ctx, id); err != nil {
		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete department. error: %w", err)
	}
	return nil
}
