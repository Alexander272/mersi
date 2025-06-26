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
	GetById(ctx context.Context, req *models.GetSiByIdDTO) (*models.BaseSI, error)
	Create(ctx context.Context, dto *models.SiDTO) error
	Update(ctx context.Context, dto *models.SiDTO) error
	ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error
}

func (s *SIService) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get si. error: %w", err)
	}
	return data, nil
}

func (s *SIService) GetById(ctx context.Context, req *models.GetSiByIdDTO) (*models.BaseSI, error) {
	instrument, err := s.instrument.GetById(ctx, &models.GetInstrumentByIdDTO{Id: req.Id})
	if err != nil {
		return nil, err
	}
	verification, err := s.verification.GetLast(ctx, &models.GetVerificationDTO{InstrumentId: req.Id})
	if err != nil {
		return nil, err
	}
	data := &models.BaseSI{
		Instrument:   instrument,
		Verification: verification,
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

func (s *SIService) Update(ctx context.Context, dto *models.SiDTO) error {
	if err := s.instrument.Update(ctx, dto.Instrument); err != nil {
		return err
	}
	if dto.Verification != nil {
		if err := s.verification.Update(ctx, dto.Verification); err != nil {
			return err
		}
	}
	return nil
}

func (s *SIService) ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error {
	if err := s.instrument.ChangePosition(ctx, dto); err != nil {
		return err
	}
	return nil
}
