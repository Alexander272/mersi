package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type InstrumentService struct {
	repo        repository.Instrument
	docs        Document
	activityLog ActivityLog
	txManager   TransactionManager
}

func NewInstrumentService(repo repository.Instrument, docs Document, activityLog ActivityLog, txManager TransactionManager) *InstrumentService {
	return &InstrumentService{
		repo:        repo,
		docs:        docs,
		activityLog: activityLog,
		txManager:   txManager,
	}
}

type Instrument interface {
	GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	Create(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error
	// Create(ctx context.Context, dto *models.InstrumentDTO) error
	CreateSeveral(ctx context.Context, dto []*models.InstrumentDTO) error
	Update(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error
	ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error
	ChangeStatus(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error
	ChangeSeveralStatuses(ctx context.Context, dto []*models.UpdateStatus) error
	Delete(ctx context.Context, dto *models.DeleteSiDTO) error
}

func (s *InstrumentService) GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get instrument by id. error: %w", err)
	}
	return data, nil
}

func (s *InstrumentService) GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error) {
	data, err := s.repo.GetUniqueData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique data for field. error: %w", err)
	}
	return data, nil
}

func (s *InstrumentService) Create(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	if tx == nil {
		return s.txManager.ExecuteInTx(ctx, func(newTx postgres.Tx) error {
			return s.executeCreate(ctx, newTx, dto)
		})
	}
	return s.executeCreate(ctx, tx, dto)
}

func (s *InstrumentService) executeCreate(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	if err := s.repo.CreateInTx(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create instrument. error: %w", err)
	}

	if dto.ActOfEnteringId != "" {
		pathDTO := &models.PathParts{
			InstrumentId: dto.Id,
			Group:        "act",
			UserId:       dto.UserId,
			IdWasEmpty:   true,
		}
		if err := s.docs.ChangePath(ctx, pathDTO); err != nil {
			return err
		}
	}

	// Логирование создания
	s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
		TableName:  "instruments",
		RecordId:   dto.Id,
		RecordName: dto.Name,
		Action:     "CREATE",
		UserId:     dto.Actor.ID,
		UserName:   dto.Actor.Name,
		NewValue:   dto,
	})
	return nil
}

// func (s *InstrumentService) Create(ctx context.Context, dto *models.InstrumentDTO) error {
// 	if err := s.repo.Create(ctx, dto); err != nil {
// 		return fmt.Errorf("failed to create instrument. error: %w", err)
// 	}

// 	if dto.ActOfEnteringId != "" {
// 		pathDTO := &models.PathParts{
// 			InstrumentId: dto.Id,
// 			Group:        "act",
// 			UserId:       dto.UserId,
// 			IdWasEmpty:   true,
// 		}
// 		if err := s.docs.ChangePath(ctx, pathDTO); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (s *InstrumentService) CreateSeveral(ctx context.Context, dto []*models.InstrumentDTO) error {
	if len(dto) == 0 {
		return nil
	}

	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several instruments. error: %w", err)
	}
	return nil
}

func (s *InstrumentService) Update(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
	// Получаем старые данные для логирования
	oldData, err := s.repo.GetById(ctx, &models.GetInstrumentByIdDTO{Id: dto.Id})
	if err != nil {
		return fmt.Errorf("failed to get old instrument data. error: %w", err)
	}

	if err := s.repo.Update(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to update instrument. error: %w", err)
	}

	// Логирование обновления
	s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
		TableName:  "instruments",
		RecordId:   dto.Id,
		RecordName: dto.Name,
		Action:     "UPDATE",
		UserId:     dto.Actor.ID,
		UserName:   dto.Actor.Name,
		OldValue:   oldData,
		NewValue:   dto,
	})
	return nil
}

func (s *InstrumentService) ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error {
	if err := s.repo.ChangePosition(ctx, dto); err != nil {
		return fmt.Errorf("failed to change position. error: %w", err)
	}
	return nil
}

func (s *InstrumentService) ChangeStatus(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error {
	if err := s.repo.ChangeStatus(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to change instrument status. error: %w", err)
	}
	return nil
}
func (s *InstrumentService) ChangeSeveralStatuses(ctx context.Context, dto []*models.UpdateStatus) error {
	if err := s.repo.ChangeSeveralStatuses(ctx, dto); err != nil {
		return fmt.Errorf("failed to change several instrument statuses. error: %w", err)
	}
	return nil
}

func (s *InstrumentService) Delete(ctx context.Context, dto *models.DeleteSiDTO) error {
	// Получаем данные перед удалением для логирования
	oldData, err := s.repo.GetById(ctx, &models.GetInstrumentByIdDTO{Id: dto.Id})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return fmt.Errorf("failed to get instrument before delete. error: %w", err)
	}

	if err := s.repo.Delete(ctx, dto.Id); err != nil {
		return fmt.Errorf("failed to delete instrument. error: %w", err)
	}

	// Логирование удаления
	if oldData != nil {
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "instruments",
			RecordId:   dto.Id,
			RecordName: oldData.Name,
			Action:     "DELETE",
			UserId:     dto.Actor.ID,
			UserName:   dto.Actor.Name,
			OldValue:   oldData,
		})
	}
	return nil
}
