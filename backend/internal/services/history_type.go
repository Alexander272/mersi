package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type HistoryTypeService struct {
	repo repository.HistoryType
}

func NewHistoryTypeService(repo repository.HistoryType) *HistoryTypeService {
	return &HistoryTypeService{repo: repo}
}

type HistoryType interface {
	Get(ctx context.Context, dto *models.GetHistoryTypesDTO) ([]*models.HistoryType, error)
	Create(ctx context.Context, dto *models.HistoryTypeDTO) error
	Update(ctx context.Context, dto *models.HistoryTypeDTO) error
	Delete(ctx context.Context, dto *models.DeleteHistoryTypeDTO) error
}

func (s *HistoryTypeService) Get(ctx context.Context, dto *models.GetHistoryTypesDTO) ([]*models.HistoryType, error) {
	data, err := s.repo.Get(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get history types. error: %w", err)
	}
	return data, nil
}

func (s *HistoryTypeService) Create(ctx context.Context, dto *models.HistoryTypeDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create history type. error: %w", err)
	}
	return nil
}

func (s *HistoryTypeService) Update(ctx context.Context, dto *models.HistoryTypeDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update history type. error: %w", err)
	}
	return nil
}

func (s *HistoryTypeService) Delete(ctx context.Context, dto *models.DeleteHistoryTypeDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete history type. error: %w", err)
	}
	return nil
}
