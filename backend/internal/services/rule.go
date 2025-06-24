package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type RuleService struct {
	repo repository.Rule
	item RuleItem
}

func NewRuleService(repo repository.Rule, item RuleItem) *RuleService {
	return &RuleService{
		repo: repo,
		item: item,
	}
}

type Rule interface {
	// GetAll(context.Context) ([]*models.RuleFull, error)
	GetAll(context.Context) ([]*models.Rule, error)
	Create(context.Context, *models.RuleDTO) error
	Update(context.Context, *models.RuleDTO) error
	Delete(context.Context, string) error
}

// func (s *RuleService) GetAll(ctx context.Context) ([]*models.RuleFull, error) {
// 	Rule, err := s.repo.GetAll(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to get all Rule. error: %w", err)
// 	}

// 	items, err := s.item.GetAll(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	RuleFull := []*models.RuleFull{}

// 	for i, m := range Rule {
// 		RuleItem := &models.RuleItem{}
// 		for _, item := range items {
// 			if m.RuleItemId == item.Id {
// 				RuleItem = item
// 				break
// 			}
// 		}

// 		logger.Debug("get Rule", logger.AnyAttr("Rule", m))

// 		if i == 0 || RuleFull[len(RuleFull)-1].Id != m.RoleId {
// 			RuleFull = append(RuleFull, &models.RuleFull{
// 				Id: m.RoleId,
// 				Role: &models.RoleFull{
// 					Id:      m.RoleId,
// 					Name:    m.RoleName,
// 					Level:   m.RoleLevel,
// 					Extends: m.RoleExtends,
// 				},
// 				RuleItems: []*models.RuleItem{RuleItem},
// 			})
// 		} else {
// 			RuleFull[len(RuleFull)-1].RuleItems = append(RuleFull[len(RuleFull)-1].RuleItems, RuleItem)
// 		}

// 	}

// 	return RuleFull, nil
// }

func (s *RuleService) GetAll(ctx context.Context) ([]*models.Rule, error) {
	Rule, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all Rule. error: %w", err)
	}
	return Rule, nil
}

func (s *RuleService) Create(ctx context.Context, Rule *models.RuleDTO) error {
	if err := s.repo.Create(ctx, Rule); err != nil {
		return fmt.Errorf("failed to create Rule. error: %w", err)
	}
	return nil
}

func (s *RuleService) Update(ctx context.Context, Rule *models.RuleDTO) error {
	if err := s.repo.Update(ctx, Rule); err != nil {
		return fmt.Errorf("failed to update Rule. error: %w", err)
	}
	return nil
}

func (s *RuleService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete Rule. error: %w", err)
	}
	return nil
}
