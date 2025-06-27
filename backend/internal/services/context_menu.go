package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type ContextService struct {
	repo repository.ContextMenu
	role Role
}

func NewContextService(repo repository.ContextMenu, role Role) *ContextService {
	return &ContextService{
		repo: repo,
		role: role,
	}
}

type ContextMenu interface {
	Get(ctx context.Context, req *models.GetContextMenuDTO) ([]*models.ContextMenu, error)
	Create(ctx context.Context, dto *models.ContextMenuDTO) error
	Update(ctx context.Context, dto *models.ContextMenuDTO) error
	Delete(ctx context.Context, dto *models.DeleteContextMenuDTO) error
}

func (s *ContextService) Get(ctx context.Context, req *models.GetContextMenuDTO) ([]*models.ContextMenu, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get context menu. error: %w", err)
	}

	role, err := s.role.Get(ctx, req.Role)
	if err != nil {
		return nil, err
	}

	res := []*models.ContextMenu{}
	filterMap := make(map[string]bool)

	for _, v := range role.Rules {
		filterMap[v] = true
	}

	for _, v := range data {
		if _, ok := filterMap[v.Rule]; ok {
			res = append(res, v)
		}
	}

	return res, nil
}

func (s *ContextService) Create(ctx context.Context, dto *models.ContextMenuDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create context menu. error: %w", err)
	}
	return nil
}

func (s *ContextService) Update(ctx context.Context, dto *models.ContextMenuDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update context menu. error: %w", err)
	}
	return nil
}

func (s *ContextService) Delete(ctx context.Context, dto *models.DeleteContextMenuDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete context menu. error: %w", err)
	}
	return nil
}
