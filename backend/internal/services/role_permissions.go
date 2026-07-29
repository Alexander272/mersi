package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type RolePermissionService struct {
	repo repository.Permissions
}

func NewRolePermissionService(repo repository.Permissions) *RolePermissionService {
	return &RolePermissionService{
		repo: repo,
	}
}

type RolePermissions interface {
	GetAll(ctx context.Context) ([]*models.Permission, error)
	GetRolePermissions(ctx context.Context, roleId string) (map[string]bool, error)
	GetInherited(ctx context.Context, roleId string) (map[string]bool, error)
	CountForAll(ctx context.Context, descendants map[string][]string) (map[string]models.PermsWithCount, error)
	ReplacePermissions(ctx context.Context, tx postgres.Tx, roleId string, permissionIds []string) error
}

func (s *RolePermissionService) GetAll(ctx context.Context) ([]*models.Permission, error) {
	return s.repo.GetAll(ctx)
}

func (s *RolePermissionService) GetRolePermissions(ctx context.Context, roleId string) (map[string]bool, error) {
	return s.repo.GetRolePermissionsMap(ctx, roleId)
}

func (s *RolePermissionService) GetInherited(ctx context.Context, roleId string) (map[string]bool, error) {
	data, err := s.repo.GetInheritedByRole(ctx, roleId)
	if err != nil {
		return nil, err
	}
	result := make(map[string]bool, len(data))
	for k := range data {
		result[k] = true
	}
	return result, nil
}

func (s *RolePermissionService) CountForAll(ctx context.Context, descendants map[string][]string) (map[string]models.PermsWithCount, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *RolePermissionService) ReplacePermissions(ctx context.Context, tx postgres.Tx, roleId string, permissionIds []string) error {
	return s.repo.ReplacePermissions(ctx, tx, roleId, permissionIds)
}
