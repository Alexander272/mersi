package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type SectionService struct {
	repo repository.Section
}

func NewSectionService(repo repository.Section) *SectionService {
	return &SectionService{
		repo: repo,
	}
}

type Section interface {
	Get(ctx context.Context, req *models.GetSectionsDTO) ([]*models.Section, error)
	Create(ctx context.Context, dto *models.SectionDTO) error
	Update(ctx context.Context, dto *models.SectionDTO) error
	Delete(ctx context.Context, dto *models.DeleteSectionDTO) error
}

func (s *SectionService) Get(ctx context.Context, req *models.GetSectionsDTO) ([]*models.Section, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sections. error: %w", err)
	}
	return data, nil
}

func (s *SectionService) Create(ctx context.Context, dto *models.SectionDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create section. error: %w", err)
	}
	return nil
}

func (s *SectionService) Update(ctx context.Context, dto *models.SectionDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update section. error: %w", err)
	}
	return nil
}

func (s *SectionService) Delete(ctx context.Context, dto *models.DeleteSectionDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete section. error: %w", err)
	}
	return nil
}
