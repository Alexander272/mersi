package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type SortingService struct {
	repo repository.Sorting
}

func NewSortingService(repo repository.Sorting) *SortingService {
	return &SortingService{
		repo: repo,
	}
}

type Sorting interface {
	Get(ctx context.Context, req *models.GetSortingDTO) ([]*models.Sorting, error)
	Create(ctx context.Context, dto *models.SortingDTO) error
	CreateSeveral(ctx context.Context, dto []*models.SortingDTO) error
	Update(ctx context.Context, dto *models.SortingDTO) error
	Change(ctx context.Context, dto []*models.SortingDTO) error
	Delete(ctx context.Context, dto *models.DeleteSortingDTO) error
	DeleteAll(ctx context.Context, dto *models.DeleteSortingDTO) error
}

func (s *SortingService) Get(ctx context.Context, req *models.GetSortingDTO) ([]*models.Sorting, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sorting. error: %w", err)
	}
	return data, nil
}

func (s *SortingService) Create(ctx context.Context, dto *models.SortingDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create sorting. error: %w", err)
	}
	return nil
}

func (s *SortingService) CreateSeveral(ctx context.Context, dto []*models.SortingDTO) error {
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several sorting. error: %w", err)
	}
	return nil
}

func (s *SortingService) Update(ctx context.Context, dto *models.SortingDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update sorting. error: %w", err)
	}
	return nil
}

func (s *SortingService) Change(ctx context.Context, dto []*models.SortingDTO) error {
	if err := s.DeleteAll(ctx, &models.DeleteSortingDTO{
		UserId:    dto[0].UserId,
		SectionId: dto[0].SectionId,
	}); err != nil {
		return err
	}

	if err := s.CreateSeveral(ctx, dto); err != nil {
		return err
	}
	return nil
}

func (s *SortingService) Delete(ctx context.Context, dto *models.DeleteSortingDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete sorting. error: %w", err)
	}
	return nil
}

func (s *SortingService) DeleteAll(ctx context.Context, dto *models.DeleteSortingDTO) error {
	if err := s.repo.DeleteAll(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete all sorting. error: %w", err)
	}
	return nil
}
