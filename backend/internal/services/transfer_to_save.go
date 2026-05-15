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
	repo        repository.TransferToSave
	txManager   TransactionManager
	instrument  Instrument
	activityLog ActivityLog
}

func NewTransferToSaveService(repo repository.TransferToSave, txManager TransactionManager, instrument Instrument, activityLog ActivityLog) *TransferToSaveService {
	return &TransferToSaveService{
		repo:        repo,
		txManager:   txManager,
		instrument:  instrument,
		activityLog: activityLog,
	}
}

type TransferToSave interface {
	Get(ctx context.Context, req *models.GetTransferToSaveDTO) ([]*models.TransferToSave, error)
	GetById(ctx context.Context, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error)
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

func (s *TransferToSaveService) GetById(ctx context.Context, req *models.GetTransferToSaveDTO) (*models.TransferToSave, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get transfer to save by id. error: %w", err)
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
		if candidate != nil && candidate.DateEnd.After(dto.DateStart) {
			return models.ErrNotValid
		}

		duplicate, err := s.repo.GetByInstrumentAndDateStart(ctx, tx, dto.InstrumentId, dto.DateStart)
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}
		if duplicate != nil {
			return models.ErrAlreadyExists
		}

		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create transfer to save. error: %w", err)
		}

		status := models.InstrumentStatusWork
		if dto.DateEnd.IsZero() {
			status = models.InstrumentStatusSaved
		}
		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: status,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}

		recordName := fmt.Sprintf("%s - %s", dto.DateStart.Format("02.01.2006"), dto.DateEnd.Format("02.01.2006"))
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "transfer_to_save",
			RecordId:   dto.Id,
			RecordName: recordName,
			Action:     "CREATE",
			UserId:     dto.Actor.ID,
			UserName:   dto.Actor.Name,
			NewValue:   dto,
		})
		return nil
	})
}

func (s *TransferToSaveService) CreateSeveral(ctx context.Context, dto []*models.TransferToSaveDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several transfers to save. error: %w", err)
	}
	return nil
}

func (s *TransferToSaveService) Update(ctx context.Context, dto *models.TransferToSaveDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if !dto.DateEnd.IsZero() && dto.DateEnd.Before(dto.DateStart) {
			return models.ErrNotValid
		}

		oldData, err := s.repo.GetById(ctx, &models.GetTransferToSaveDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update transfer to save. error: %w", err)
		}

		status := models.InstrumentStatusWork
		if dto.DateEnd.IsZero() {
			status = models.InstrumentStatusSaved
		}
		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: status,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}

		if oldData != nil {
			recordName := fmt.Sprintf("%s - %s", dto.DateStart.Format("02.01.2006"), dto.DateEnd.Format("02.01.2006"))
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "transfer_to_save",
				RecordId:   dto.Id,
				RecordName: recordName,
				Action:     "UPDATE",
				UserId:     dto.Actor.ID,
				UserName:   dto.Actor.Name,
				OldValue:   oldData,
				NewValue:   dto,
			})
		}
		return nil
	})
}

func (s *TransferToSaveService) Delete(ctx context.Context, dto *models.DeleteTransferToSaveDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получение записи для удаления
		oldData, err := s.GetById(ctx, &models.GetTransferToSaveDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get transfer to save data. error: %w", err)
		}

		// Проверка: является ли удаляемая запись последней (один вызов GetLast)
		isLast := false
		if oldData != nil {
			lastData, err := s.GetLast(ctx, tx, &models.GetTransferToSaveDTO{InstrumentId: oldData.InstrumentId})
			if err != nil && !errors.Is(err, models.ErrNoRows) {
				return err
			}
			if lastData != nil && lastData.Id == oldData.Id {
				isLast = true
			}
		}

		// Удаление записи
		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete transfer to save. error: %w", err)
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
			recordName := fmt.Sprintf("%s - %s", oldData.DateStart.Format("02.01.2006"), oldData.DateEnd.Format("02.01.2006"))
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "transfer_to_save",
				RecordId:   dto.Id,
				RecordName: recordName,
				Action:     "DELETE",
				UserId:     dto.Actor.ID,
				UserName:   dto.Actor.Name,
				OldValue:   oldData,
			})
		}

		return nil
	})
}
