package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/services/most"
	"github.com/Alexander272/mersi/backend/pkg/logger"
)

type ReceivingService struct {
	location     Location
	repo         repository.Location
	responsible  Responsible
	notification Notification
	most         most.Dialog
	activityLog  ActivityLog
	txManager    TransactionManager
}

type ReceivingDeps struct {
	Location     Location
	Repo         repository.Location
	Responsible  Responsible
	Notification Notification
	Most         most.Dialog
	ActivityLog  ActivityLog
	TxManager    TransactionManager
}

func NewReceivingService(deps *ReceivingDeps) *ReceivingService {
	return &ReceivingService{
		location:     deps.Location,
		repo:         deps.Repo,
		responsible:  deps.Responsible,
		notification: deps.Notification,
		most:         deps.Most,
		activityLog:  deps.ActivityLog,
		txManager:    deps.TxManager,
	}
}

type Receiving interface {
	ReceivingFromApp(ctx context.Context, dto *models.ReceivingDTO) error
	ReceivingFromChannel(ctx context.Context, dto *models.DialogResponse) error
	ReceivingDialogOpen(ctx context.Context, dto *models.PostAction) error
	ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error
	ForcedReceiptAll(ctx context.Context) error
}

// Получение инструментов. Запрос прилетает из веб-приложения
func (s *ReceivingService) ReceivingFromApp(ctx context.Context, dto *models.ReceivingDTO) error {
	responsible, err := s.responsible.GetBySSOId(ctx, dto.UserId)
	if err != nil {
		return err
	}

	logger.Debug("receiving", logger.StringAttr("status", dto.Status), logger.AnyAttr("responsible", responsible))
	if dto.Status == constants.LocationStatusUsed {
		if len(responsible) == 0 {
			return models.ErrNotResponsible
		}

		ids := []string{}
		for _, r := range responsible {
			ids = append(ids, r.DepartmentId)
		}

		// при подтверждении инструментов я отфильтровываю те что не находятся в том же подразделении что и пользователь
		filtered, err := s.location.SelectByDepartments(ctx, &models.SelectByDepsDTO{
			InstrumentIds: dto.InstrumentIds,
			DepartmentIds: ids,
			Status:        constants.LocationStatusMoved,
		})
		if err != nil {
			return err
		}

		isFull := len(filtered) == len(dto.InstrumentIds)
		if !isFull {
			dto.InstrumentIds = filtered
		}
	}

	if len(dto.InstrumentIds) == 0 {
		return models.ErrCannotConfirmReceipt
	}

	if err := s.repo.Receiving(ctx, dto); err != nil {
		return fmt.Errorf("failed to receiving si location. error: %w", err)
	}
	return nil
}

// Получение инструментов. Запрос прилетает из канала mattermost
func (s *ReceivingService) ReceivingFromChannel(ctx context.Context, dto *models.DialogResponse) error {
	status := constants.LocationStatusUsed
	instrumentIds := []string{}

	state := strings.Split(dto.State, "&")
	for _, s := range state {
		arr := strings.SplitN(s, ":", 2)
		if arr[0] == "Status" {
			status = arr[1]
		}
	}

	for k, v := range dto.Submission {
		if v {
			instrumentIds = append(instrumentIds, k)
		}
	}

	location := &models.ReceivingDTO{
		InstrumentIds: instrumentIds,
		Status:        status,
		HasConfirmed:  true,
	}
	if err := s.repo.Receiving(ctx, location); err != nil {
		return fmt.Errorf("failed to receiving si location. error: %w", err)
	}

	if err := s.notification.CheckReceiving(ctx, dto); err != nil {
		return err
	}

	return nil
}
func (s *ReceivingService) ReceivingDialogOpen(ctx context.Context, dto *models.PostAction) error {
	if err := s.most.Open(ctx, dto); err != nil {
		return err
	}
	return nil
}

func (s *ReceivingService) ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error {
	// Получаем старые данные для логирования
	oldData, err := s.location.GetById(ctx, &models.GetLocationDTO{InstrumentId: dto.InstrumentId})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return fmt.Errorf("failed to get old location data. error: %w", err)
	}

	if err := s.repo.ForcedReceipt(ctx, dto); err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return models.ErrInstrumentReceived
		}
		return fmt.Errorf("failed to forced receipt si. error: %w", err)
	}

	// Логирование
	if oldData != nil {
		recordName := fmt.Sprintf("%s (%s)", oldData.Place, oldData.Status)
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "locations",
			RecordId:   oldData.Id,
			RecordName: recordName,
			Action:     "FORCED_RECEIPT",
			UserId:     dto.Actor.ID,
			UserName:   dto.Actor.Name,
			OldValue:   oldData,
			NewValue:   &models.Location{Status: constants.LocationStatusUsed, DateOfReceiving: oldData.DateOfReceiving},
		})
	}

	return nil
}

func (s *ReceivingService) ForcedReceiptAll(ctx context.Context) error {
	logger.Info("Forced receipt SI")
	if err := s.repo.ForcedReceiptAll(ctx); err != nil {
		return fmt.Errorf("failed to forced receipt all si. error: %w", err)
	}
	return nil
}
