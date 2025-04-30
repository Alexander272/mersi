package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type ColumnsService struct {
	repo repository.Columns
}

func NewColumnsService(repo repository.Columns) *ColumnsService {
	return &ColumnsService{
		repo: repo,
	}
}

type Columns interface {
	Get(ctx context.Context, req *models.GetColumnsDTO) ([]*models.Column, error)
	Create(ctx context.Context, dto *models.ColumnsDTO) error
	CreateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error
	Update(ctx context.Context, dto *models.ColumnsDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error
	UpdatePositions(ctx context.Context, dto []*models.UpdateColumnPosition) error
	Delete(ctx context.Context, dto *models.DeleteColumnDTO) error
	DeleteAll(ctx context.Context, dto *models.DeleteColumnsDTO) error
}

func (s *ColumnsService) Get(ctx context.Context, req *models.GetColumnsDTO) ([]*models.Column, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get columns. error: %w", err)
	}
	return data, nil
}

func (s *ColumnsService) Create(ctx context.Context, dto *models.ColumnsDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create column. error: %w", err)
	}
	return nil
}

func (s *ColumnsService) CreateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error {
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create columns. error: %w", err)
	}
	return nil
}

func (s *ColumnsService) Update(ctx context.Context, dto *models.ColumnsDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update column. error: %w", err)
	}
	return nil
}

func (s *ColumnsService) UpdateSeveral(ctx context.Context, dto []*models.ColumnsDTO) error {
	if err := s.repo.UpdateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to update few columns. error: %w", err)
	}
	return nil
}

func (s *ColumnsService) UpdatePositions(ctx context.Context, dto []*models.UpdateColumnPosition) error {
	if err := s.repo.UpdatePositions(ctx, dto); err != nil {
		return fmt.Errorf("failed to update column positions. error: %w", err)
	}
	return nil
}

func (s *ColumnsService) Delete(ctx context.Context, dto *models.DeleteColumnDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete column. error: %w", err)
	}
	return nil
}

func (s *ColumnsService) DeleteAll(ctx context.Context, dto *models.DeleteColumnsDTO) error {
	if err := s.repo.DeleteAll(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete columns by section. error: %w", err)
	}
	return nil
}
