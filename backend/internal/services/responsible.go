package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type ResponsibleService struct {
	repo repository.Responsible
}

func NewResponsibleService(repo repository.Responsible) *ResponsibleService {
	return &ResponsibleService{
		repo: repo,
	}
}

type Responsible interface {
	Get(ctx context.Context, req *models.GetResponsibleDTO) ([]*models.Responsible, error)
	GetBySSOId(ctx context.Context, id string) ([]*models.Responsible, error)
	Change(ctx context.Context, dto *models.ChangeResponsibleDTO) error
	Create(ctx context.Context, dto *models.ResponsibleDTO) error
	CreateSeveral(ctx context.Context, dto []*models.ResponsibleDTO) error
	Update(ctx context.Context, dto *models.ResponsibleDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.ResponsibleDTO) error
	Delete(ctx context.Context, id string) error
	DeleteSeveral(ctx context.Context, ids []string) error
}

func (s *ResponsibleService) Get(ctx context.Context, req *models.GetResponsibleDTO) ([]*models.Responsible, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get responsible. error: %w", err)
	}
	return data, nil
}

func (s *ResponsibleService) GetBySSOId(ctx context.Context, id string) ([]*models.Responsible, error) {
	data, err := s.repo.GetBySSOId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get responsible. error: %w", err)
	}
	return data, nil
}

func (s *ResponsibleService) Change(ctx context.Context, dto *models.ChangeResponsibleDTO) error {
	if err := s.CreateSeveral(ctx, dto.New); err != nil {
		return err
	}
	if err := s.UpdateSeveral(ctx, dto.Updated); err != nil {
		return err
	}
	if err := s.DeleteSeveral(ctx, dto.Deleted); err != nil {
		return err
	}
	return nil
}

func (s *ResponsibleService) Create(ctx context.Context, dto *models.ResponsibleDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create responsible. error: %w", err)
	}
	return nil
}

func (s *ResponsibleService) CreateSeveral(ctx context.Context, dto []*models.ResponsibleDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to create responsible. error: %w", err)
	}
	return nil
}

func (s *ResponsibleService) Update(ctx context.Context, dto *models.ResponsibleDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update responsible. error: %w", err)
	}
	return nil
}

func (s *ResponsibleService) UpdateSeveral(ctx context.Context, dto []*models.ResponsibleDTO) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.UpdateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to update responsible. error: %w", err)
	}
	return nil
}

func (s *ResponsibleService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete responsible. error: %w", err)
	}
	return nil
}

func (s *ResponsibleService) DeleteSeveral(ctx context.Context, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.DeleteSeveral(ctx, ids); err != nil {
		return fmt.Errorf("failed to delete responsible. error: %w", err)
	}
	return nil
}
