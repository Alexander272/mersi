package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type VerificationDocService struct {
	repo repository.VerificationDoc
}

func NewVerificationDocService(repo repository.VerificationDoc) *VerificationDocService {
	return &VerificationDocService{
		repo: repo,
	}
}

type VerificationDoc interface {
	Get(ctx context.Context, req *models.GetVerificationDocsDTO) ([]*models.VerificationDoc, error)
	GetGrouped(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error)
	CreateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error
	Update(ctx context.Context, dto *models.VerificationDocDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerificationDocDTO) error
}

func (s *VerificationDocService) Get(ctx context.Context, req *models.GetVerificationDocsDTO) ([]*models.VerificationDoc, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification documents. error: %w", err)
	}
	return data, nil
}

func (s *VerificationDocService) GetGrouped(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error) {
	data, err := s.repo.GetGrouped(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get grouped verification documents. error: %w", err)
	}
	return data, nil
}

func (s *VerificationDocService) CreateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several verification documents. error: %w", err)
	}
	return nil
}

func (s *VerificationDocService) Update(ctx context.Context, dto *models.VerificationDocDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update verification document. error: %w", err)
	}
	return nil
}

func (s *VerificationDocService) UpdateSeveral(ctx context.Context, dto []*models.VerificationDocDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.UpdateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to update several verification documents. error: %w", err)
	}
	return nil
}

func (s *VerificationDocService) Delete(ctx context.Context, dto *models.DeleteVerificationDocDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete verification document. error: %w", err)
	}
	return nil
}
