package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type FilterService struct {
	repo repository.Filters
}

func NewFilterService(repo repository.Filters) *FilterService {
	return &FilterService{repo: repo}
}

type Filters interface {
	Get(ctx context.Context, dto *models.GetSavedFiltersDTO) ([]*models.SavedFilter, error)
	CreateOne(ctx context.Context, dto *models.SavedFilterDTO) error
	Create(ctx context.Context, dto []*models.SavedFilterDTO) error
	Change(ctx context.Context, dto []*models.SavedFilterDTO) error
	Delete(ctx context.Context, dto *models.DeleteSavedFiltersDTO) error
}

func (s *FilterService) Get(ctx context.Context, req *models.GetSavedFiltersDTO) ([]*models.SavedFilter, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get saved filters. error: %w", err)
	}
	return data, nil
}

func (s *FilterService) CreateOne(ctx context.Context, dto *models.SavedFilterDTO) error {
	if err := s.repo.CreateOne(ctx, dto); err != nil {
		return fmt.Errorf("failed to create saved filter. error: %w", err)
	}
	return nil
}

func (s *FilterService) Create(ctx context.Context, dto []*models.SavedFilterDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create saved filters. error: %w", err)
	}
	return nil
}

func (s *FilterService) Change(ctx context.Context, dto []*models.SavedFilterDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if err := s.Delete(ctx, &models.DeleteSavedFiltersDTO{
		UserId:    dto[0].UserId,
		SectionId: dto[0].SectionId,
	}); err != nil {
		return err
	}

	if err := s.Create(ctx, dto); err != nil {
		return err
	}
	return nil
}

func (s *FilterService) Delete(ctx context.Context, dto *models.DeleteSavedFiltersDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete saved filters. error: %w", err)
	}
	return nil
}
