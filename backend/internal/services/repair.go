package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type RepairService struct {
	repo repository.Repair
}

func NewRepairService(repo repository.Repair) *RepairService {
	return &RepairService{
		repo: repo,
	}
}

type Repair interface {
	Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error)
	Create(ctx context.Context, dto *models.RepairDTO) error
	Update(ctx context.Context, dto *models.RepairDTO) error
	Delete(ctx context.Context, dto *models.DeleteRepairDTO) error
}

func (s *RepairService) Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair info. error: %w", err)
	}
	return data, nil
}

func (s *RepairService) Create(ctx context.Context, dto *models.RepairDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create repair info. error: %w", err)
	}
	return nil
}

func (s *RepairService) Update(ctx context.Context, dto *models.RepairDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update repair info. error: %w", err)
	}
	return nil
}

func (s *RepairService) Delete(ctx context.Context, dto *models.DeleteRepairDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete repair info. error: %w", err)
	}
	return nil
}
