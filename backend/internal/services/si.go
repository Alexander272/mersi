package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type SIService struct {
	repo         repository.SI
	instrument   Instrument
	verification Verification
}

type SiDeps struct {
	Repo         repository.SI
	Instrument   Instrument
	Verification Verification
}

func NewSiService(deps *SiDeps) *SIService {
	return &SIService{
		repo:         deps.Repo,
		instrument:   deps.Instrument,
		verification: deps.Verification,
	}
}

type SI interface {
	Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error)
	Create(ctx context.Context, dto *models.SiDTO) error
}

func (s *SIService) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get si. error: %w", err)
	}
	return data, nil
}

func (s *SIService) Create(ctx context.Context, dto *models.SiDTO) error {
	if err := s.instrument.Create(ctx, dto.Instrument); err != nil {
		return err
	}
	if dto.Verification != nil {
		dto.Verification.InstrumentId = dto.Instrument.Id
		dto.Verification.Status = string(models.InstrumentStatusWork)
		if err := s.verification.Create(ctx, dto.Verification); err != nil {
			s.instrument.Delete(ctx, dto.Instrument.Id)
			return err
		}
	}

	return nil
}
