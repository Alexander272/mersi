package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type TransferToSaveService struct {
	repo       repository.TransferToSave
	txManager  TransactionManager
	instrument Instrument
}

func NewTransferToSaveService(repo repository.TransferToSave, txManager TransactionManager, instrument Instrument) *TransferToSaveService {
	return &TransferToSaveService{
		repo:       repo,
		txManager:  txManager,
		instrument: instrument,
	}
}

type TransferToSave interface {
	Get(ctx context.Context, req *models.GetTransferToSaveDTO) ([]*models.TransferToSave, error)
	GetLast(ctx context.Context, tx postgres.Tx, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error)
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

func (s *TransferToSaveService) GetLast(ctx context.Context, tx postgres.Tx, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error) {
	data, err := s.repo.GetLast(ctx, tx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last transfer to save. error: %w", err)
	}
	return data, nil
}

func (s *TransferToSaveService) Create(ctx context.Context, dto *models.TransferToSaveDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		candidate, err := s.GetLast(ctx, tx, &models.GetTransferToSaveDTO{InstrumentId: dto.InstrumentId})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}
		// if candidate != nil && candidate.DateEnd > dto.DateStart {
		if candidate != nil && candidate.DateEnd.After(dto.DateStart) {
			return models.ErrNotValid
		}

		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create transfer to save. error: %w", err)
		}

		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusSaved,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}
		return nil
	})
}

func (s *TransferToSaveService) CreateSeveral(ctx context.Context, dto []*models.TransferToSaveDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create transfer to save. error: %w", err)
	}
	return nil
}

func (s *TransferToSaveService) Update(ctx context.Context, dto *models.TransferToSaveDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// if dto.DateEnd < dto.DateStart {
		if dto.DateEnd.Before(dto.DateStart) {
			return models.ErrNotValid
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update transfer to save. error: %w", err)
		}

		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusWork,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}
		return nil
	})
}

func (s *TransferToSaveService) Delete(ctx context.Context, dto *models.DeleteTransferToSaveDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete transfer to save. error: %w", err)
	}
	return nil
}
