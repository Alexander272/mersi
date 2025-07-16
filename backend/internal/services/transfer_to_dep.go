package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type TransferToDepService struct {
	repo repository.TransferToDepartment
}

func NewTransferToDepService(repo repository.TransferToDepartment) *TransferToDepService {
	return &TransferToDepService{repo: repo}
}

type TransferToDepartment interface {
	Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error)
	Create(ctx context.Context, dto *models.TransferToDepartmentDTO) error
	Update(ctx context.Context, dto *models.TransferToDepartmentDTO) error
	Delete(ctx context.Context, dto *models.DeleteTransferToDepDTO) error
}

func (s *TransferToDepService) Get(ctx context.Context, req *models.GetTransferToDepDTO) ([]*models.TransferToDepartment, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transfers to department. error: %w", err)
	}
	return data, nil
}

func (s *TransferToDepService) Create(ctx context.Context, dto *models.TransferToDepartmentDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create transfer to department. error: %w", err)
	}
	//TODO надо еще как-то статус менять или делать что-то подобное
	return nil
}

func (s *TransferToDepService) Update(ctx context.Context, dto *models.TransferToDepartmentDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update transfer to department. error: %w", err)
	}
	return nil
}

func (s *TransferToDepService) Delete(ctx context.Context, dto *models.DeleteTransferToDepDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete transfer to department. error: %w", err)
	}
	return nil
}
