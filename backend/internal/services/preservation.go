package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type PreservationService struct {
	repo repository.Preservation
}

func NewPreservationService(repo repository.Preservation) *PreservationService {
	return &PreservationService{
		repo: repo,
	}
}

type Preservation interface {
	Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error)
	GetLast(ctx context.Context, req *models.GetPreservationsDTO) (*models.Preservation, error)
	Create(ctx context.Context, dto *models.PreservationDTO) error
	Update(ctx context.Context, dto *models.PreservationDTO) error
	Delete(ctx context.Context, dto *models.DeletePreservationDTO) error
}

func (s *PreservationService) Get(ctx context.Context, req *models.GetPreservationsDTO) ([]*models.Preservation, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get preservations by instrument id. error: %w", err)
	}
	return data, nil
}

func (s *PreservationService) GetLast(ctx context.Context, req *models.GetPreservationsDTO) (*models.Preservation, error) {
	data, err := s.repo.GetLast(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last preservation. error: %w", err)
	}
	return data, nil
}

func (s *PreservationService) Create(ctx context.Context, dto *models.PreservationDTO) error {
	candidate, err := s.GetLast(ctx, &models.GetPreservationsDTO{InstrumentId: dto.InstrumentId})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return err
	}
	if candidate != nil && candidate.DateEnd > dto.DateStart {
		return models.ErrNotValid
	}
	//TODO наверное надо еще как-то статус менять или делать что-то подобное, чтобы можно было легко и просто понять нужна ли поверка для инструмента

	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create preservation. error: %w", err)
	}
	return nil
}

func (s *PreservationService) Update(ctx context.Context, dto *models.PreservationDTO) error {
	if dto.DateEnd < dto.DateStart {
		return models.ErrNotValid
	}

	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update preservation. error: %w", err)
	}
	return nil
}

func (s *PreservationService) Delete(ctx context.Context, dto *models.DeletePreservationDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete preservation. error: %w", err)
	}
	return nil
}
