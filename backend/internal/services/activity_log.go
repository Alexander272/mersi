package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/goccy/go-json"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/logger"
)

type ActivityLogService struct {
	repo repository.ActivityLog
	wg   sync.WaitGroup
}

func NewActivityLogService(repo repository.ActivityLog) *ActivityLogService {
	return &ActivityLogService{repo: repo}
}

type ActivityLog interface {
	Create(ctx context.Context, dto *models.CreateActivityLogDTO)
	CreateSeveral(ctx context.Context, dto []*models.CreateActivityLogDTO)
	GetByRecord(ctx context.Context, dto *models.GetActivityLogDTO) ([]*models.ActivityLog, error)
	GetAll(ctx context.Context, dto *models.ActivityLogFilter) ([]*models.ActivityLog, error)
	LogActivity(ctx context.Context, dto *models.CreateActivityLogDTO)
	Wait()
}

func (s *ActivityLogService) Create(ctx context.Context, dto *models.CreateActivityLogDTO) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.repo.Create(ctx, dto); err != nil {
			logger.Error("failed to create activity log async", logger.ErrAttr(err))
			error_bot.Send(nil, fmt.Sprintf("failed to create activity log async: %v", err), dto)
		}
	}()
}

func (s *ActivityLogService) CreateSeveral(ctx context.Context, dto []*models.CreateActivityLogDTO) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.repo.CreateSeveral(ctx, dto); err != nil {
			logger.Error("failed to create several activity log async", logger.ErrAttr(err))
			error_bot.Send(nil, fmt.Sprintf("failed to create several activity log async: %v", err), dto)
		}
	}()
}

func (s *ActivityLogService) Wait() {
	s.wg.Wait()
}

func (s *ActivityLogService) GetByRecord(ctx context.Context, dto *models.GetActivityLogDTO) ([]*models.ActivityLog, error) {
	data, err := s.repo.GetByRecord(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity log by record. error: %w", err)
	}
	return data, nil
}

func (s *ActivityLogService) GetAll(ctx context.Context, dto *models.ActivityLogFilter) ([]*models.ActivityLog, error) {
	data, err := s.repo.GetAll(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get activity log. error: %w", err)
	}
	return data, nil
}

func (s *ActivityLogService) LogActivity(ctx context.Context, dto *models.CreateActivityLogDTO) {
	if dto == nil {
		return
	}
	var err error
	// Если OldValue или NewValue не []byte (json.RawMessage), то маршалим
	if dto.OldValue != nil {
		if _, ok := dto.OldValue.([]byte); !ok {
			dto.OldValue, err = json.Marshal(dto.OldValue)
			if err != nil {
				logger.Error("failed to marshal old value", logger.ErrAttr(err))
				error_bot.Send(nil, fmt.Sprintf("failed to marshal old value: %v", err), dto.OldValue)
			}
		}
	}
	if dto.NewValue != nil {
		if _, ok := dto.NewValue.([]byte); !ok {
			dto.NewValue, err = json.Marshal(dto.NewValue)
			if err != nil {
				logger.Error("failed to marshal new value", logger.ErrAttr(err))
				error_bot.Send(nil, fmt.Sprintf("failed to marshal new value: %v", err), dto.NewValue)
			}
		}
	}

	s.Create(ctx, dto)
}
