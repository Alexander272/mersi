package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type RoleHierarchyService struct {
	repo repository.RoleHierarchy
}

func NewRoleHierarchyService(repo repository.RoleHierarchy) *RoleHierarchyService {
	return &RoleHierarchyService{
		repo: repo,
	}
}

type RoleHierarchy interface {
	LoadPolicy(ctx context.Context) ([]*models.SyncRoleInheritance, error)
	GetRoleDescendants(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error)
	GetDirectChildren(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error)
	AddInheritance(ctx context.Context, tx postgres.Tx, dto *models.RoleHierarchyDTO) error
	AddInheritances(ctx context.Context, tx postgres.Tx, roleId string, parentRoleIds []string) error
	RemoveInheritance(ctx context.Context, tx postgres.Tx, dto *models.RoleHierarchyDTO) error
	RemoveInheritances(ctx context.Context, tx postgres.Tx, roleId string, parentRoleIds []string) error
}

func (s *RoleHierarchyService) LoadPolicy(ctx context.Context) ([]*models.SyncRoleInheritance, error) {
	data, err := s.repo.LoadPolicy(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load role policy: %w", err)
	}
	return data, nil
}

func (s *RoleHierarchyService) GetDirectChildren(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error) {
	data, err := s.repo.GetDirectChildren(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get direct children: %w", err)
	}
	return data, nil
}

func (s *RoleHierarchyService) GetRoleDescendants(ctx context.Context, req *models.GetRolesInheritance) (map[string][]string, error) {
	data, err := s.repo.GetRoleDescendants(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get role descendants: %w", err)
	}
	return data, nil
}

func (s *RoleHierarchyService) AddInheritance(ctx context.Context, tx postgres.Tx, dto *models.RoleHierarchyDTO) error {
	if dto.ParentRoleId == dto.RoleId {
		return models.ErrCannotInheritFromSelf
	}

	if err := s.repo.AddInheritance(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to add inheritance: %w", err)
	}
	return nil
}

func (s *RoleHierarchyService) RemoveInheritance(ctx context.Context, tx postgres.Tx, dto *models.RoleHierarchyDTO) error {
	if err := s.repo.RemoveInheritance(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to remove inheritance: %w", err)
	}
	return nil
}

func (s *RoleHierarchyService) AddInheritances(ctx context.Context, tx postgres.Tx, roleId string, parentRoleIds []string) error {
	for _, parentId := range parentRoleIds {
		if roleId == parentId {
			return models.ErrCannotInheritFromSelf
		}
	}

	if err := s.repo.AddInheritances(ctx, tx, roleId, parentRoleIds); err != nil {
		return fmt.Errorf("failed to add inheritances: %w", err)
	}
	return nil
}

func (s *RoleHierarchyService) RemoveInheritances(ctx context.Context, tx postgres.Tx, roleId string, parentRoleIds []string) error {
	if err := s.repo.RemoveInheritances(ctx, tx, roleId, parentRoleIds); err != nil {
		return fmt.Errorf("failed to remove inheritances: %w", err)
	}
	return nil
}
