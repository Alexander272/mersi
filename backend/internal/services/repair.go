package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type RepairService struct {
	repo       repository.Repair
	txManager  TransactionManager
	instrument Instrument
}

func NewRepairService(repo repository.Repair, txManager TransactionManager, instrument Instrument) *RepairService {
	return &RepairService{
		repo:       repo,
		txManager:  txManager,
		instrument: instrument,
	}
}

type Repair interface {
	Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error)
	GetLast(ctx context.Context, req *models.GetRepairDTO) (*models.Repair, error)
	Create(ctx context.Context, dto *models.RepairDTO) error
	CreateSeveral(ctx context.Context, dto []*models.RepairDTO) error
	Update(ctx context.Context, dto *models.RepairDTO) error
	Delete(ctx context.Context, dto *models.DeleteRepairDTO) error
}

func (s *RepairService) Get(ctx context.Context, req *models.GetRepairDTO) ([]*models.Repair, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get repair info. error: %w", err)
	}
	return data, nil
}

func (s *RepairService) GetLast(ctx context.Context, req *models.GetRepairDTO) (*models.Repair, error) {
	data, err := s.repo.GetLast(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last repair info. error: %w", err)
	}
	return data, nil
}

func (s *RepairService) Create(ctx context.Context, dto *models.RepairDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create repair info. error: %w", err)
		}

		status := models.InstrumentStatusWork
		if dto.PeriodEnd.IsZero() {
			status = models.InstrumentStatusRepair
		}
		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: status,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}

		return nil
	})
}

func (s *RepairService) CreateSeveral(ctx context.Context, dto []*models.RepairDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create repair info. error: %w", err)
	}
	return nil
}

func (s *RepairService) Update(ctx context.Context, dto *models.RepairDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update repair info. error: %w", err)
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

func (s *RepairService) Delete(ctx context.Context, dto *models.DeleteRepairDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete repair info. error: %w", err)
	}
	return nil
}
