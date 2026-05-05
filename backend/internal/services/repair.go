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
	repo        repository.Repair
	txManager   TransactionManager
	instrument  Instrument
	activityLog ActivityLog
}

func NewRepairService(repo repository.Repair, txManager TransactionManager, instrument Instrument, activityLog ActivityLog) *RepairService {
	return &RepairService{
		repo:        repo,
		txManager:   txManager,
		instrument:  instrument,
		activityLog: activityLog,
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
		if !dto.PeriodEnd.IsZero() && dto.PeriodEnd.Before(dto.PeriodStart) {
			return models.ErrNotValid
		}

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

		recordName := fmt.Sprintf("%s %s-%s", dto.Defect, dto.PeriodStart.Format("02.01.2006"), dto.PeriodEnd.Format("02.01.2006"))
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "repair",
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
		if !dto.PeriodEnd.IsZero() && dto.PeriodEnd.Before(dto.PeriodStart) {
			return models.ErrNotValid
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update repair info. error: %w", err)
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

		recordName := fmt.Sprintf("%s %s-%s", dto.Defect, dto.PeriodStart.Format("02.01.2006"), dto.PeriodEnd.Format("02.01.2006"))
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "repair",
			RecordId:   dto.Id,
			RecordName: recordName,
			Action:     "UPDATE",
			UserId:     dto.Actor.ID,
			UserName:   dto.Actor.Name,
			NewValue:   dto,
		})
		return nil
	})
}

func (s *RepairService) Delete(ctx context.Context, dto *models.DeleteRepairDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получение записи для удаления
		oldData, err := s.repo.GetById(ctx, dto.Id)
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get repair data. error: %w", err)
		}

		// Проверка: является ли удаляемая запись последней (один вызов GetLast)
		isLast := false
		if oldData != nil {
			lastData, err := s.GetLast(ctx, &models.GetRepairDTO{InstrumentId: oldData.InstrumentId})
			if err != nil && !errors.Is(err, models.ErrNoRows) {
				return err
			}
			if lastData != nil && lastData.Id == oldData.Id {
				isLast = true
			}
		}

		// Удаление записи
		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete repair. error: %w", err)
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
			recordName := fmt.Sprintf("%s %s-%s", oldData.Defect, oldData.PeriodStart.Format("02.01.2006"), oldData.PeriodEnd.Format("02.01.2006"))
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "repair",
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
