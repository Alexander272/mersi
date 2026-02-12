package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type TransferToDepService struct {
	repo       repository.TransferToDepartment
	txManager  TransactionManager
	instrument Instrument
	docs       Document
}

func NewTransferToDepService(repo repository.TransferToDepartment, txManager TransactionManager, instrument Instrument, docs Document) *TransferToDepService {
	return &TransferToDepService{
		repo:       repo,
		txManager:  txManager,
		instrument: instrument,
		docs:       docs,
	}
}

type TransferToDepartment interface {
	Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error)
	Create(ctx context.Context, dto *models.TransferToDepartmentDTO) error
	CreateSeveral(ctx context.Context, dto []*models.TransferToDepartmentDTO) error
	Update(ctx context.Context, dto *models.TransferToDepartmentDTO) error
	Delete(ctx context.Context, dto *models.DeleteTransferToDepDTO) error
}

func (s *TransferToDepService) Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers to department. error: %w", err)
	}
	return data, nil
}

func (s *TransferToDepService) Create(ctx context.Context, dto *models.TransferToDepartmentDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create transfer to department. error: %w", err)
		}

		if dto.DocId != "" {
			pathDTO := &models.PathParts{
				InstrumentId: dto.InstrumentId,
				Group:        "transferToDep",
				UserId:       dto.UserId,
			}
			if err := s.docs.ChangePath(ctx, pathDTO); err != nil {
				return err
			}
		}

		status := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusTransferred,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, status); err != nil {
			return fmt.Errorf("failed to change instrument status. error: %w", err)
		}

		return nil
	})
}

func (s *TransferToDepService) CreateSeveral(ctx context.Context, dto []*models.TransferToDepartmentDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several transfers to department. error: %w", err)
	}
	return nil
}

func (s *TransferToDepService) Update(ctx context.Context, dto *models.TransferToDepartmentDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update transfer to department. error: %w", err)
	}
	return nil
}

func (s *TransferToDepService) Delete(ctx context.Context, dto *models.DeleteTransferToDepDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete transfer to department. error: %w", err)
	}
	return nil
}
