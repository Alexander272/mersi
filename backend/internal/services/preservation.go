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
	repo        repository.Preservation
	txManager   TransactionManager
	instrument  Instrument
	activityLog ActivityLog
}

func NewPreservationService(repo repository.Preservation, txManager TransactionManager, instrument Instrument, activityLog ActivityLog) *PreservationService {
	return &PreservationService{
		repo:        repo,
		txManager:   txManager,
		instrument:  instrument,
		activityLog: activityLog,
	}
}

type Preservation interface {
	Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error)
	GetById(ctx context.Context, req *models.GetPreservationsDTO) (*models.Preservation, error)
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

func (s *PreservationService) GetById(ctx context.Context, req *models.GetPreservationsDTO) (*models.Preservation, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get preservation by id. error: %w", err)
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
		if candidate != nil && candidate.DateEnd.After(dto.DateStart) {
			return models.ErrNotValid
		}

		duplicate, err := s.repo.GetByInstrumentAndDateStart(ctx, tx, dto.InstrumentId, dto.DateStart)
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}
		if duplicate != nil {
			return models.ErrPreservationAlreadyExists
		}

		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create preservation. error: %w", err)
		}

		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusWork,
		}
		if dto.DateEnd.IsZero() {
			instrumentDTO.Status = models.InstrumentStatusArchived
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}

		recordName := fmt.Sprintf("%s - %s", dto.DateStart.Format("02.01.2006"), dto.DateEnd.Format("02.01.2006"))
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "preservations",
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

func (s *PreservationService) CreateSeveral(ctx context.Context, dto []*models.PreservationDTO) error {
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several preservations. error: %w", err)
	}
	return nil
}

func (s *PreservationService) Update(ctx context.Context, dto *models.PreservationDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if !dto.DateEnd.IsZero() && dto.DateEnd.Before(dto.DateStart) {
			return models.ErrNotValid
		}

		oldData, err := s.repo.GetById(ctx, &models.GetPreservationsDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update preservation. error: %w", err)
		}

		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusWork,
		}
		if dto.DateEnd.IsZero() {
			instrumentDTO.Status = models.InstrumentStatusArchived
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}

		if oldData != nil {
			recordName := fmt.Sprintf("%s - %s", dto.DateStart.Format("02.01.2006"), dto.DateEnd.Format("02.01.2006"))
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "preservations",
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

func (s *PreservationService) Delete(ctx context.Context, dto *models.DeletePreservationDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получение записи для удаления
		oldData, err := s.GetById(ctx, &models.GetPreservationsDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get preservation data. error: %w", err)
		}

		// Проверка: является ли удаляемая запись последней (один вызов GetLast)
		isLast := false
		if oldData != nil {
			lastData, err := s.GetLast(ctx, tx, &models.GetPreservationsDTO{InstrumentId: oldData.InstrumentId})
			if err != nil && !errors.Is(err, models.ErrNoRows) {
				return err
			}
			if lastData != nil && lastData.Id == oldData.Id {
				isLast = true
			}
		}

		// Удаление записи
		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete preservation. error: %w", err)
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
				TableName:  "preservations",
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
