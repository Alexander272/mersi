package services

import (
	"context"

	"github.com/Alexander272/mersi/backend/internal/models"
)

type SIService struct {
	instrument Instrument
}

type SiDeps struct {
	Instrument Instrument
}

func NewSiService(deps *SiDeps) *SIService {
	return &SIService{
		instrument: deps.Instrument,
	}
}

type SI interface {
	Create(ctx context.Context, dto *models.SiDTO) error
}

func (s *SIService) Create(ctx context.Context, dto *models.SiDTO) error {
	if err := s.instrument.Create(ctx, dto.Instrument); err != nil {
		return err
	}
	if dto.Verification != nil {
		dto.Verification.InstrumentId = dto.Instrument.Id
	}

	return nil
}
