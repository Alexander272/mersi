package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type VerificationService struct {
	repo        repository.Verification
	txManager   TransactionManager
	verDocs     VerificationDoc
	instrument  Instrument
	docs        Document
	activityLog ActivityLog
}

type VerificationDeps struct {
	Repo        repository.Verification
	TxManager   TransactionManager
	VerDocs     VerificationDoc
	Instrument  Instrument
	Docs        Document
	ActivityLog ActivityLog
}

func NewVerificationService(deps *VerificationDeps) *VerificationService {
	return &VerificationService{
		repo:        deps.Repo,
		txManager:   deps.TxManager,
		verDocs:     deps.VerDocs,
		instrument:  deps.Instrument,
		docs:        deps.Docs,
		activityLog: deps.ActivityLog,
	}
}

type Verification interface {
	Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error)
	GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error)
	Create(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	CreateSeveral(ctx context.Context, dto []*models.VerificationDTO) error
	Update(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerificationDTO) error
}

func (s *VerificationService) Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification by instrument. error: %w", err)
	}

	docs, err := s.verDocs.GetGrouped(ctx, &models.GetGroupedVerificationDocsDTO{InstrumentId: req.InstrumentId})
	if err != nil {
		return nil, err
	}

	for i := range data {
		if group, ok := docs.Groups[data[i].Id]; ok {
			data[i].Docs = group.Docs
		}
	}

	return data, nil
}

func (s *VerificationService) GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
	data, err := s.repo.GetLast(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last verification. error: %w", err)
	}

	docs, err := s.verDocs.GetGrouped(ctx, &models.GetGroupedVerificationDocsDTO{InstrumentId: req.InstrumentId})
	if err != nil {
		return nil, err
	}

	if group, ok := docs.Groups[data.Id]; ok {
		data.Docs = group.Docs
	}

	return data, nil
}

func (s *VerificationService) Create(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.ExecuteInTx(ctx, func(newTx postgres.Tx) error {
			return s.executeCreate(ctx, newTx, dto)
		})
	}

	// Если транзакция передана, используем её
	return s.executeCreate(ctx, tx, dto)
}
func (s *VerificationService) executeCreate(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if err := s.repo.CreateInTx(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create verification. error: %w", err)
	}

	if len(dto.Docs) > 0 {
		for i := range dto.Docs {
			dto.Docs[i].VerificationId = dto.Id
		}
		if err := s.verDocs.CreateSeveral(ctx, tx, dto.Docs); err != nil {
			return err
		}
	}

	instrumentDTO := &models.UpdateStatus{
		Id:     dto.InstrumentId,
		Status: models.InstrumentStatus(dto.Status),
	}
	if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
		return err
	}

	recordName := fmt.Sprintf("%s - %s",
		dto.Date.Format("02.01.2006"),
		dto.NextDate.Format("02.01.2006"),
	)
	// Логирование создания
	s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
		TableName:  "verifications",
		RecordId:   dto.Id,
		RecordName: recordName,
		Action:     "CREATE",
		UserId:     dto.Actor.ID,
		UserName:   dto.Actor.Name,
		NewValue:   dto,
	})

	return nil
}

func (s *VerificationService) CreateSeveral(ctx context.Context, dto []*models.VerificationDTO) error {
	if len(dto) == 0 {
		return nil
	}
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.repo.CreateSeveral(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create verifications. error: %w", err)
		}

		docs := []*models.VerificationDocDTO{}
		for i := range dto {
			for j := range dto[i].Docs {
				dto[i].Docs[j].VerificationId = dto[i].Id
			}
			docs = append(docs, dto[i].Docs...)
		}
		if err := s.verDocs.CreateSeveral(ctx, tx, docs); err != nil {
			return err
		}

		return nil
	})
}

func (s *VerificationService) Update(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.ExecuteInTx(ctx, func(newTx postgres.Tx) error {
			return s.executeUpdate(ctx, newTx, dto)
		})
	}

	// Если транзакция передана, используем её
	return s.executeUpdate(ctx, tx, dto)
}
func (s *VerificationService) executeUpdate(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	// Получаем старые данные для логирования
	oldData, err := s.repo.GetById(ctx, dto.Id)
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return fmt.Errorf("failed to get old verification data. error: %w", err)
	}

	if err := s.repo.Update(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to update verification. error: %w", err)
	}

	newDocs := []*models.VerificationDocDTO{}
	updatedDocs := []*models.VerificationDocDTO{}
	for i := range dto.Docs {
		if dto.Docs[i].Id == "" {
			dto.Docs[i].VerificationId = dto.Id
			newDocs = append(newDocs, dto.Docs[i])
		} else {
			updatedDocs = append(updatedDocs, dto.Docs[i])
		}
	}

	if len(newDocs) > 0 {
		if err := s.verDocs.CreateSeveral(ctx, tx, newDocs); err != nil {
			return err
		}
	}
	if len(updatedDocs) > 0 {
		if err := s.verDocs.UpdateSeveral(ctx, tx, updatedDocs); err != nil {
			return err
		}
	}

	// Обработка удаляемых документов
	if len(dto.DeletedDocs) > 0 {
		for _, doc := range dto.DeletedDocs {
			deleteDto := &models.DeleteDocumentDTO{
				Id:           doc.DocId,
				Filename:     doc.Filename,
				Group:        "verification",
				InstrumentId: dto.InstrumentId,
				UserId:       dto.UserId,
				IsTemp:       false,
			}
			if err := s.docs.Delete(ctx, deleteDto); err != nil {
				return fmt.Errorf("failed to delete document: %w", err)
			}
			// Удаляем связь из verification_docs
			if err := s.verDocs.DeleteByDocId(ctx, tx, doc.DocId); err != nil {
				return err
			}
		}
	}

	instrumentDTO := &models.UpdateStatus{
		Id:     dto.InstrumentId,
		Status: models.InstrumentStatus(dto.Status),
	}
	if err := s.instrument.ChangeStatus(ctx, tx, instrumentDTO); err != nil {
		return err
	}

	// Логирование обновления
	if oldData != nil {
		recordName := fmt.Sprintf("%s - %s",
			dto.Date.Format("02.01.2006"),
			dto.NextDate.Format("02.01.2006"),
		)

		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "verifications",
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
}

func (s *VerificationService) Delete(ctx context.Context, dto *models.DeleteVerificationDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получение записи для удаления
		oldData, err := s.repo.GetById(ctx, dto.Id)
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get verification data. error: %w", err)
		}

		// Проверка: является ли удаляемая запись последней (один вызов GetLast)
		isLast := false
		if oldData != nil {
			lastData, err := s.GetLast(ctx, &models.GetVerificationDTO{InstrumentId: oldData.InstrumentId})
			if err != nil && !errors.Is(err, models.ErrNoRows) {
				return err
			}
			if lastData != nil && lastData.Id == oldData.Id {
				isLast = true
			}
		}

		// Удаление записи
		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to delete verification. error: %w", err)
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
			recordName := fmt.Sprintf("%s - %s",
				oldData.Date.Format("02.01.2006"),
				oldData.NextDate.Format("02.01.2006"),
			)
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "verifications",
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
