package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type DepartmentAccessService struct {
	repo repository.DepartmentAccess
}

func NewDepartmentAccessService(repo repository.DepartmentAccess) *DepartmentAccessService {
	return &DepartmentAccessService{repo: repo}
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

func (s *DepartmentAccessService) Get(ctx context.Context, req *models.GetDepartmentAccessDTO) ([]*models.DepartmentAccess, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get department accesses. error: %w", err)
	}
	return data, nil
}

func (s *DepartmentAccessService) GetByUserId(ctx context.Context, req *models.GetDepartmentAccessDTO) ([]*models.DepartmentAccess, error) {
	data, err := s.repo.GetByUserId(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get department accesses. error: %w", err)
	}
	return data, nil
}

func (s *DepartmentAccessService) Replace(ctx context.Context, dto *models.ReplaceDepartmentAccessDTO) error {
	err := s.repo.Replace(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to replace department accesses. error: %w", err)
	}
	return nil
}

func (s *DepartmentAccessService) Create(ctx context.Context, dto *models.DepartmentAccessDTO) error {
	err := s.repo.Create(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to create department access. error: %w", err)
	}
	return nil
}

func (s *DepartmentAccessService) CreateSeveral(ctx context.Context, dto []*models.DepartmentAccessDTO) error {
	err := s.repo.CreateSeveral(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to create several department accesses. error: %w", err)
	}
	return nil
}

func (s *DepartmentAccessService) Update(ctx context.Context, dto *models.DepartmentAccessDTO) error {
	err := s.repo.Update(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to update department access. error: %w", err)
	}
	return nil
}

func (s *DepartmentAccessService) UpdateSeveral(ctx context.Context, dto []*models.DepartmentAccessDTO) error {
	err := s.repo.UpdateSeveral(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to update several department accesses. error: %w", err)
	}
	return nil
}

func (s *DepartmentAccessService) Delete(ctx context.Context, dto *models.DeleteDepartmentAccessDTO) error {
	err := s.repo.Delete(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to delete department access. error: %w", err)
	}
	return nil
}
