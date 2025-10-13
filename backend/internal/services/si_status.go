package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type StatusService struct {
	repo repository.SiStatus
}

func NewStatusService(repo repository.SiStatus) *StatusService {
	return &StatusService{
		repo: repo,
	}
}

type SiStatus interface {
	Get(ctx context.Context, req *models.GetSiStatusDTO) ([]*models.SiStatus, error)
	Create(ctx context.Context, dto *models.SiStatusDTO) error
	Update(ctx context.Context, dto *models.SiStatusDTO) error
	Delete(ctx context.Context, dto *models.DeleteSiStatusDTO) error
}

func (s *StatusService) Get(ctx context.Context, req *models.GetSiStatusDTO) ([]*models.SiStatus, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get si status. error: %w", err)
	}
	return data, nil
}

func (s *StatusService) Create(ctx context.Context, dto *models.SiStatusDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create si status. error: %w", err)
	}
	return nil
}

func (s *StatusService) Update(ctx context.Context, dto *models.SiStatusDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update si status. error: %w", err)
	}
	return nil
}

func (s *StatusService) Delete(ctx context.Context, dto *models.DeleteSiStatusDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete si status. error: %w", err)
	}
	return nil
}
