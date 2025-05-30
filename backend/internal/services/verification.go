package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type VerificationService struct {
	repo repository.Verification
	docs VerificationDoc
}

func NewVerificationService(repo repository.Verification, docs VerificationDoc) *VerificationService {
	return &VerificationService{
		repo: repo,
		docs: docs,
	}
}

type Verification interface {
	Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error)
	GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error)
	Create(ctx context.Context, dto *models.VerificationDTO) error
	Update(ctx context.Context, dto *models.VerificationDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerificationDTO) error
}

func (s *VerificationService) Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification by instrument. error: %w", err)
	}

	docs, err := s.docs.GetGrouped(ctx, &models.GetGroupedVerificationDocsDTO{InstrumentId: req.InstrumentId})
	if err != nil {
		return nil, err
	}

	for i := range data {
		group, exists := docs.Groups[data[i].Id]
		if exists {
			data[i].Docs = group.Docs
		}
	}
	return data, nil
}

func (s *VerificationService) GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
	data, err := s.repo.GetLast(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last verification. error: %w", err)
	}

	docs, err := s.docs.Get(ctx, &models.GetVerificationDocsDTO{VerificationId: data.Id})
	if err != nil {
		return nil, err
	}
	data.Docs = docs

	return data, nil
}

func (s *VerificationService) Create(ctx context.Context, dto *models.VerificationDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create verification. error: %w", err)
	}

	for i := range dto.Docs {
		dto.Docs[i].VerificationId = dto.Id
	}
	if err := s.docs.CreateSeveral(ctx, dto.Docs); err != nil {
		s.Delete(ctx, &models.DeleteVerificationDTO{Id: dto.Id})
		return err
	}
	return nil
}

func (s *VerificationService) Update(ctx context.Context, dto *models.VerificationDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update verification. error: %w", err)
	}

	if err := s.docs.UpdateSeveral(ctx, dto.Docs); err != nil {
		return err
	}
	return nil
}

func (s *VerificationService) Delete(ctx context.Context, dto *models.DeleteVerificationDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete verification. error: %w", err)
	}
	return nil
}
