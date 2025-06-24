package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type RuleItemService struct {
	repo repository.RuleItem
}

func NewRuleItemService(repo repository.RuleItem) *RuleItemService {
	return &RuleItemService{
		repo: repo,
	}
}

type RuleItem interface {
	GetAll(context.Context) ([]*models.RuleItem, error)
	Create(context.Context, *models.RuleItemDTO) error
	Update(context.Context, *models.RuleItemDTO) error
	Delete(context.Context, string) error
}

func (s *RuleItemService) GetAll(ctx context.Context) ([]*models.RuleItem, error) {
	RuleItems, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all menu items. error: %w", err)
	}
	return RuleItems, nil
}

func (s *RuleItemService) Create(ctx context.Context, menu *models.RuleItemDTO) error {
	if err := s.repo.Create(ctx, menu); err != nil {
		return fmt.Errorf("failed to create menu item. error: %w", err)
	}
	return nil
}

func (s *RuleItemService) Update(ctx context.Context, menu *models.RuleItemDTO) error {
	if err := s.repo.Update(ctx, menu); err != nil {
		return fmt.Errorf("failed to update menu item. error: %w", err)
	}
	return nil
}

func (s *RuleItemService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete menu item. error: %w", err)
	}
	return nil
}
