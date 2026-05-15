package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type WriteOffService struct {
	repo        repository.WriteOff
	txManager   TransactionManager
	instrument  Instrument
	docs        Document
	activityLog ActivityLog
}

func NewWriteOffService(repo repository.WriteOff, txManager TransactionManager, instrument Instrument, docs Document, activityLog ActivityLog) *WriteOffService {
	return &WriteOffService{
		repo:        repo,
		txManager:   txManager,
		instrument:  instrument,
		docs:        docs,
		activityLog: activityLog,
	}
}

type WriteOff interface {
	Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error)
	GetById(ctx context.Context, req *models.GetWriteOffDTO) (*models.WriteOff, error)
	GetLast(ctx context.Context, req *models.GetWriteOffDTO) (*models.WriteOff, error)
	Create(ctx context.Context, dto *models.WriteOffDTO) error
	CreateSeveral(ctx context.Context, dto []*models.WriteOffDTO) error
	Update(ctx context.Context, dto *models.WriteOffDTO) error
	Delete(ctx context.Context, dto *models.DeleteWriteOffDTO) error
}

func (s *WriteOffService) Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get write off. error: %w", err)
	}
	return data, nil
}

func (s *WriteOffService) GetById(ctx context.Context, req *models.GetWriteOffDTO) (*models.WriteOff, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get write off by id. error: %w", err)
	}
	return data, nil
}

func (s *WriteOffService) GetLast(ctx context.Context, req *models.GetWriteOffDTO) (*models.WriteOff, error) {
	data, err := s.repo.GetLast(ctx, nil, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last write off. error: %w", err)
	}
	return data, nil
}

func (s *WriteOffService) Create(ctx context.Context, dto *models.WriteOffDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		candidate, err := s.repo.GetByInstrumentAndDate(ctx, tx, dto.InstrumentId, dto.Date)
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return err
		}
		if candidate != nil {
			return models.ErrAlreadyExists
		}

		if err := s.repo.Create(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create write off. error: %w", err)
		}

		if dto.DocId != "" {
			pathDTO := &models.PathParts{
				InstrumentId: dto.InstrumentId,
				Group:        "writeOff",
				UserId:       dto.UserId,
			}
			if err := s.docs.ChangePath(ctx, pathDTO); err != nil {
				return err
			}
		}

		instrumentDTO := &models.UpdateStatus{
			Id:     dto.InstrumentId,
			Status: models.InstrumentStatusWriteOff,
		}
		if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
			return err
		}

		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "write_off",
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

func (s *WriteOffService) CreateSeveral(ctx context.Context, dto []*models.WriteOffDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create several write offs. error: %w", err)
	}
	return nil
}

func (s *WriteOffService) Update(ctx context.Context, dto *models.WriteOffDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		oldData, err := s.repo.GetById(ctx, &models.GetWriteOffDTO{Id: dto.Id})
		if err != nil {
			return fmt.Errorf("failed to get old write off data. error: %w", err)
		}

		if err := s.repo.Update(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to update write off. error: %w", err)
		}

		// Обработка удаляемых документов
		if len(dto.DeletedDocs) >0 {
			for _, doc := range dto.DeletedDocs {
				deleteDto := &models.DeleteDocumentDTO{
					Id:           doc.DocId,
					Filename:     doc.Filename,
					Group:        "writeOff",
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
				TableName:  "write_off",
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

func (s *WriteOffService) Delete(ctx context.Context, dto *models.DeleteWriteOffDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получение записи для удаления
		oldData, err := s.GetById(ctx, &models.GetWriteOffDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get write off data. error: %w", err)
		}

		// Проверка: является ли удаляемая запись последней (один вызов GetLast)
		isLast := false
		if oldData != nil {
			lastData, err := s.GetLast(ctx, &models.GetWriteOffDTO{InstrumentId: oldData.InstrumentId})
			if err != nil && !errors.Is(err, models.ErrNoRows) {
				return err
			}
			if lastData != nil && lastData.Id == oldData.Id {
				isLast = true
			}
		}

		// Удаление записи
		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete write off. error: %w", err)
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
				TableName:  "write_off",
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
