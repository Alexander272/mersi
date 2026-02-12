package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type PreservationService struct {
	repo       repository.Preservation
	txManager  TransactionManager
	instrument Instrument
}

func NewPreservationService(repo repository.Preservation, txManager TransactionManager, instrument Instrument) *PreservationService {
	return &PreservationService{
		repo:       repo,
		txManager:  txManager,
		instrument: instrument,
	}
}

type Preservation interface {
	Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error)
	GetLast(ctx context.Context, tx postgres.Tx, req *models.GetPreservationsDTO) (*models.Preservation, error)
	Create(ctx context.Context, dto *models.PreservationDTO) error
	CreateSeveral(ctx context.Context, dto []*models.PreservationDTO) error
	Update(ctx context.Context, dto *models.PreservationDTO) error
	Delete(ctx context.Context, dto *models.DeletePreservationDTO) error
}

func (s *PreservationService) Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get preservations by instrument id. error: %w", err)
	}
	return data, nil
}

func (s *PreservationService) GetLast(ctx context.Context, tx postgres.Tx, req *models.GetPreservationsDTO) (*models.Preservation, error) {
	data, err := s.repo.GetLast(ctx, tx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last preservation. error: %w", err)
	}
	return data, nil
}

func (s *PreservationService) Create(ctx context.Context, dto *models.PreservationDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		candidate, err := s.GetLast(ctx, tx, &models.GetPreservationsDTO{InstrumentId: dto.InstrumentId})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}
		// if candidate != nil && candidate.DateEnd > dto.DateStart {
		if candidate != nil && candidate.DateEnd.After(dto.DateStart) {
			return models.ErrNotValid
		}

		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create preservation. error: %w", err)
		}

		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusArchived,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}
		return nil
	})
}

func (s *PreservationService) CreateSeveral(ctx context.Context, dto []*models.PreservationDTO) error {
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several preservations. error: %w", err)
	}
	return nil
}

func (s *PreservationService) Update(ctx context.Context, dto *models.PreservationDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// if dto.DateEnd < dto.DateStart {
		if dto.DateEnd.Before(dto.DateStart) {
			return models.ErrNotValid
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update preservation. error: %w", err)
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

func (s *PreservationService) Delete(ctx context.Context, dto *models.DeletePreservationDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete preservation. error: %w", err)
	}
	return nil
}
