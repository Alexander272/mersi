package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type CustomContextService struct {
	repo repository.CustomContextMenu
}

func NewCustomContextService(repo repository.CustomContextMenu) *CustomContextService {
	return &CustomContextService{
		repo: repo,
	}
}

type CustomContextMenu interface {
	Create(ctx context.Context, dto *models.CustomContextMenuDTO) error
	Delete(ctx context.Context, dto *models.CustomContextMenuDTO) error
}

func (s *CustomContextService) Create(ctx context.Context, dto *models.CustomContextMenuDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create custom context menu. error: %w", err)
	}
	return nil
}

func (s *CustomContextService) Delete(ctx context.Context, dto *models.CustomContextMenuDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete custom context menu. error: %w", err)
	}
	return nil
}
