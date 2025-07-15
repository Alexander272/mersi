package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type TransferToSaveService struct {
	repo repository.TransferToSave
}

func NewTransferToSaveService(repo repository.TransferToSave) *TransferToSaveService {
	return &TransferToSaveService{
		repo: repo,
	}
}

type TransferToSave interface {
	Get(ctx context.Context, req *models.GetTransferToSaveDTO) ([]*models.TransferToSave, error)
	GetLast(ctx context.Context, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error)
	Create(ctx context.Context, dto *models.TransferToSaveDTO) error
	Update(ctx context.Context, dto *models.TransferToSaveDTO) error
	Delete(ctx context.Context, dto *models.DeleteTransferToSaveDTO) error
}

func (s *TransferToSaveService) Get(ctx context.Context, req *models.GetTransferToSaveDTO) ([]*models.TransferToSave, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers to save. error: %w", err)
	}
	return data, nil
}

func (s *TransferToSaveService) GetLast(ctx context.Context, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error) {
	data, err := s.repo.GetLast(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last transfer to save. error: %w", err)
	}
	return data, nil
}

func (s *TransferToSaveService) Create(ctx context.Context, dto *models.TransferToSaveDTO) error {
	candidate, err := s.GetLast(ctx, &models.GetTransferToSaveDTO{InstrumentId: dto.InstrumentId})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return err
	}
	if candidate != nil && candidate.DateEnd > dto.DateStart {
		return models.ErrNotValid
	}

	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create transfer to save. error: %w", err)
	}
	return nil
}

func (s *TransferToSaveService) Update(ctx context.Context, dto *models.TransferToSaveDTO) error {
	if dto.DateEnd < dto.DateStart {
		return models.ErrNotValid
	}

	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update transfer to save. error: %w", err)
	}
	return nil
}

func (s *TransferToSaveService) Delete(ctx context.Context, dto *models.DeleteTransferToSaveDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete transfer to save. error: %w", err)
	}
	return nil
}
