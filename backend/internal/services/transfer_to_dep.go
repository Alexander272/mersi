package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type TransferToDepService struct {
	repo        repository.TransferToDepartment
	txManager   TransactionManager
	instrument  Instrument
	docs        Document
	activityLog ActivityLog
}

func NewTransferToDepService(repo repository.TransferToDepartment, txManager TransactionManager, instrument Instrument, docs Document, activityLog ActivityLog) *TransferToDepService {
	return &TransferToDepService{
		repo:        repo,
		txManager:   txManager,
		instrument:  instrument,
		docs:        docs,
		activityLog: activityLog,
	}
}

type TransferToDepartment interface {
	Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error)
	GetById(ctx context.Context, req *models.GetTransferToDepDTO) (*models.TransferToDepartment, error)
	GetLast(ctx context.Context, req *models.GetTransferToDepDTO) (*models.TransferToDepartment, error)
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

func (s *TransferToDepService) GetById(ctx context.Context, req *models.GetTransferToDepDTO) (*models.TransferToDepartment, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfer to department by id. error: %w", err)
	}
	return data, nil
}

func (s *TransferToDepService) GetLast(ctx context.Context, req *models.GetTransferToDepDTO) (*models.TransferToDepartment, error) {
	data, err := s.repo.GetLast(ctx, nil, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last transfer to department. error: %w", err)
	}
	return data, nil
}

func (s *TransferToDepService) Create(ctx context.Context, dto *models.TransferToDepartmentDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		candidate, err := s.repo.GetByInstrumentAndDate(ctx, tx, dto.InstrumentId, dto.Date)
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}
		if candidate != nil {
			return models.ErrAlreadyExists
		}

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

		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "transfer_to_department",
			RecordId:   dto.Id,
			RecordName: dto.DocName,
			Action:     "CREATE",
			UserId:     dto.Actor.ID,
			UserName:   dto.Actor.Name,
			NewValue:   dto,
		})
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
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		oldData, err := s.repo.GetById(ctx, &models.GetTransferToDepDTO{Id: dto.Id})
		if err != nil {
			return fmt.Errorf("failed to get old transfer to department data. error: %w", err)
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update transfer to department. error: %w", err)
		}

		// Обработка удаляемых документов
		if len(dto.DeletedDocs) >0 {
			for _, doc := range dto.DeletedDocs {
				deleteDto := &models.DeleteDocumentDTO{
					Id:           doc.DocId,
					Filename:     doc.Filename,
					Group:        "transferToDep",
					InstrumentId: dto.InstrumentId,
					UserId:       dto.UserId,
					IsTemp:       false,
				}
				if err := s.docs.Delete(ctx, deleteDto); err != nil {
					return fmt.Errorf("failed to delete document: %w", err)
				}
			}
		}

		if oldData != nil {
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "transfer_to_department",
				RecordId:   dto.Id,
				RecordName: dto.DocName,
				Action:     "UPDATE",
				UserId:     dto.Actor.ID,
				UserName:   dto.Actor.Name,
				NewValue:   dto,
				OldValue:   oldData,
			})
		}
		return nil
	})
}

func (s *TransferToDepService) Delete(ctx context.Context, dto *models.DeleteTransferToDepDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получение записи для удаления
		oldData, err := s.GetById(ctx, &models.GetTransferToDepDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get transfer to department data. error: %w", err)
		}

		// Проверка: является ли удаляемая запись последней (один вызов GetLast)
		isLast := false
		if oldData != nil {
			lastData, err := s.GetLast(ctx, &models.GetTransferToDepDTO{InstrumentId: oldData.InstrumentId})
			if err != nil && !errors.Is(err, models.ErrNoRows) {
				return err
			}
			if lastData != nil && lastData.Id == oldData.Id {
				isLast = true
			}
		}

		// Удаление записи
		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete transfer to department. error: %w", err)
		}

		// Обновление статуса только если удалена последняя запись
		if isLast {
			instrumentDTO := &models.UpdateStatus{
				Id:     oldData.InstrumentId,
				Status: models.InstrumentStatusWork,
			}
			if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
				return err
			}
		}

		// Логирование действия
		if oldData != nil {
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "transfer_to_department",
				RecordId:   dto.Id,
				RecordName: oldData.DocName,
				Action:     "DELETE",
				UserId:     dto.Actor.ID,
				UserName:   dto.Actor.Name,
				OldValue:   oldData,
			})
		}

		return nil
	})
}
