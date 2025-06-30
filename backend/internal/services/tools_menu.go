package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type ToolsMenuService struct {
	repo    repository.ToolsMenu
	context CustomContextMenu
	roles   Role
}

func NewToolsMenuService(repo repository.ToolsMenu, context CustomContextMenu, roles Role) *ToolsMenuService {
	return &ToolsMenuService{
		repo:    repo,
		context: context,
		roles:   roles,
	}
}

type ToolsMenu interface {
	Get(ctx context.Context, req *models.GetToolsMenuDTO) ([]*models.ToolsMenu, error)
	ToggleFavorite(ctx context.Context, dto *models.ChangeFavoriteDTO) error
	Create(ctx context.Context, dto *models.ToolsMenuDTO) error
	Update(ctx context.Context, dto *models.ToolsMenuDTO) error
	Delete(ctx context.Context, dto *models.DeleteToolsMenuDTO) error
}

func (s *ToolsMenuService) Get(ctx context.Context, req *models.GetToolsMenuDTO) ([]*models.ToolsMenu, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get tools menu. error: %s", err)
	}

	role, err := s.roles.Get(ctx, req.Role)
	if err != nil {
		return nil, err
	}

	res := []*models.ToolsMenu{}
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

func (s *ToolsMenuService) ToggleFavorite(ctx context.Context, dto *models.ChangeFavoriteDTO) error {
	contextDTO := &models.CustomContextMenuDTO{UserId: dto.UserId, ToolsMenuId: dto.Id}
	if dto.Favorite {
		if err := s.context.Create(ctx, contextDTO); err != nil {
			return err
		}
	} else {
		if err := s.context.Delete(ctx, contextDTO); err != nil {
			return err
		}
	}

	// if err := s.repo.ToggleFavorite(ctx, dto); err != nil {
	// 	s.context.Delete(ctx, contextDTO)
	// 	return fmt.Errorf("failed to toggle favorite menu item. error: %w", err)
	// }
	return nil
}

func (s *ToolsMenuService) Create(ctx context.Context, dto *models.ToolsMenuDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create tools menu item. error: %w", err)
	}
	return nil
}

func (s *ToolsMenuService) Update(ctx context.Context, dto *models.ToolsMenuDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update tools menu item. error: %w", err)
	}
	return nil
}

func (s *ToolsMenuService) Delete(ctx context.Context, dto *models.DeleteToolsMenuDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete tools menu item. error: %w", err)
	}
	return nil
}
