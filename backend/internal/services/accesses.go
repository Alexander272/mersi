package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type AccessesService struct {
	repo repository.Accesses
}

func NewAccessesService(repo repository.Accesses) *AccessesService {
	return &AccessesService{
		repo: repo,
	}
}

type Accesses interface {
	Get(ctx context.Context, req *models.GetAccessesDTO) ([]*models.Accesses, error)
	Create(ctx context.Context, dto *models.AccessesDTO) error
	Update(ctx context.Context, dto *models.AccessesDTO) error
	Delete(ctx context.Context, dto *models.DeleteAccessesDTO) error
}

func (s *AccessesService) Get(ctx context.Context, req *models.GetAccessesDTO) ([]*models.Accesses, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get accesses. error: %w", err)
	}
	return data, nil
}

func (s *AccessesService) Create(ctx context.Context, dto *models.AccessesDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create accesses. error: %w", err)
	}
	return nil
}

func (s *AccessesService) Update(ctx context.Context, dto *models.AccessesDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update accesses. error: %w", err)
	}
	return nil
}

func (s *AccessesService) Delete(ctx context.Context, dto *models.DeleteAccessesDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete accesses. error: %w", err)
	}
	return nil
}
