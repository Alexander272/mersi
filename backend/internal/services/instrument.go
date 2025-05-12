package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type InstrumentService struct {
	repo repository.Instrument
}

func NewInstrumentService(repo repository.Instrument) *InstrumentService {
	return &InstrumentService{
		repo: repo,
	}
}

type Instrument interface {
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	Create(ctx context.Context, dto *models.InstrumentDTO) error
	Update(ctx context.Context, dto *models.InstrumentDTO) error
	ChangeStatus(ctx context.Context, dto *models.UpdateStatus) error
}

func (s *InstrumentService) GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error) {
	data, err := s.repo.GetUniqueData(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get unique data for field. error: %w", err)
	}
	return data, nil
}

func (s *InstrumentService) Create(ctx context.Context, dto *models.InstrumentDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create instrument. error: %w", err)
	}
	return nil
}

func (s *InstrumentService) Update(ctx context.Context, dto *models.InstrumentDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update instrument. error: %w", err)
	}
	return nil
}

func (s *InstrumentService) ChangeStatus(ctx context.Context, dto *models.UpdateStatus) error {
	if err := s.repo.ChangeStatus(ctx, dto); err != nil {
		return fmt.Errorf("failed to change instrument status. error: %w", err)
	}
	return nil
}
