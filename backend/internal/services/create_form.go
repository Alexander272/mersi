package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type CreateFormService struct {
	repo repository.CreateForm
}

func NewCreateFormService(repo repository.CreateForm) *CreateFormService {
	return &CreateFormService{
		repo: repo,
	}
}

type CreateForm interface {
	Get(ctx context.Context, req *models.GetCreateFormDTO) ([]*models.CreateFormStep, error)
	Create(ctx context.Context, dto *models.CreateFormFieldDTO) error
	Update(ctx context.Context, dto *models.CreateFormFieldDTO) error
	UpdateSeveral(ctx context.Context, dto []*models.CreateFormFieldDTO) error
	Delete(ctx context.Context, dto *models.DeleteCreateFormFieldDTO) error
}

func (s *CreateFormService) Get(ctx context.Context, req *models.GetCreateFormDTO) ([]*models.CreateFormStep, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf(`failed to get 'create form' steps. error: %w`, err)
	}
	return data, nil
}

func (s *CreateFormService) Create(ctx context.Context, dto *models.CreateFormFieldDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create 'create form' field. error: %w", err)
	}
	return nil
}

func (s *CreateFormService) Update(ctx context.Context, dto *models.CreateFormFieldDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update 'create form' field. error: %w", err)
	}
	return nil
}

func (s *CreateFormService) UpdateSeveral(ctx context.Context, dto []*models.CreateFormFieldDTO) error {
	if err := s.repo.UpdateSeveral(ctx, dto); err != nil {
		return fmt.Errorf("failed to update several fields to create form. error: %w", err)
	}
	return nil
}

func (s *CreateFormService) Delete(ctx context.Context, dto *models.DeleteCreateFormFieldDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete 'create form' field. error: %w", err)
	}
	return nil
}
