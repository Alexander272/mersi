package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type VerificationFieldService struct {
	repo repository.VerificationFields
}

func NewVerificationFieldService(repo repository.VerificationFields) *VerificationFieldService {
	return &VerificationFieldService{
		repo: repo,
	}
}

type VerificationFields interface {
	Get(ctx context.Context, req *models.GetVerFieldsDTO) ([]*models.VerificationField, error)
	Create(ctx context.Context, dto *models.VerificationFieldDTO) error
	Update(ctx context.Context, dto *models.VerificationFieldDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerFieldDTO) error
}

func (s *VerificationFieldService) Get(ctx context.Context, req *models.GetVerFieldsDTO) ([]*models.VerificationField, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification fields. error: %w", err)
	}
	return data, nil
}

func (s *VerificationFieldService) Create(ctx context.Context, dto *models.VerificationFieldDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create verification field. error: %w", err)
	}
	return nil
}

func (s *VerificationFieldService) Update(ctx context.Context, dto *models.VerificationFieldDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update verification field. error: %w", err)
	}
	return nil
}

func (s *VerificationFieldService) Delete(ctx context.Context, dto *models.DeleteVerFieldDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete verification field. error: %w", err)
	}
	return nil
}
