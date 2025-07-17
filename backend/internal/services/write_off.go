package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type WriteOffService struct {
	repo repository.WriteOff
}

func NewWriteOffService(repo repository.WriteOff) *WriteOffService {
	return &WriteOffService{
		repo: repo,
	}
}

type WriteOff interface {
	Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error)
	Create(ctx context.Context, dto *models.WriteOffDTO) error
	Update(ctx context.Context, dto *models.WriteOffDTO) error
	Delete(ctx context.Context, dto *models.DeleteWriteOffDTO) error
}

func (s *WriteOffService) Get(ctx context.Context, req *models.GetWriteOffDTO) ([]*models.WriteOff, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get write off. error: %w", err)
	}
	return data, nil
}

func (s *WriteOffService) Create(ctx context.Context, dto *models.WriteOffDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create write off. error: %w", err)
	}
	//TODO надо еще как-то статус менять или делать что-то подобное
	return nil
}

func (s *WriteOffService) Update(ctx context.Context, dto *models.WriteOffDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update write off. error: %w", err)
	}
	return nil
}

func (s *WriteOffService) Delete(ctx context.Context, dto *models.DeleteWriteOffDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete write off. error: %w", err)
	}
	return nil
}
