package services

import (
	"context"
	"errors"
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
	GetById(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error)
	GetUniqueData(ctx context.Context, req *models.GetUniqueDTO) ([]string, error)
	Create(ctx context.Context, dto *models.InstrumentDTO) error
	Update(ctx context.Context, dto *models.InstrumentDTO) error
	ChangeStatus(ctx context.Context, dto *models.UpdateStatus) error
	Delete(ctx context.Context, id string) error
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

func (s *InstrumentService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete instrument. error: %w", err)
	}
	return nil
}
