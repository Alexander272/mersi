package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type RoleService struct {
	repo repository.Role
}

func NewRoleService(repo repository.Role) *RoleService {
	return &RoleService{
		repo: repo,
	}
}

type Role interface {
	GetAll(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	GetAllWithNames(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	Get(context.Context, string) (*models.Role, error)
	GetByRealm(context.Context, *models.GetRoleByRealmDTO) (*models.RoleFull, error)
	GetWithRealm(context.Context, *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error)
	Create(context.Context, *models.RoleDTO) error
	Update(context.Context, *models.RoleDTO) error
	Delete(context.Context, string) error
}

func (s *RoleService) GetAll(ctx context.Context, req *models.GetRolesDTO) ([]*models.RoleFull, error) {
	roles, err := s.repo.GetAll(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles. error: %w", err)
	}
	return roles, nil
}

func (s *RoleService) GetAllWithNames(ctx context.Context, req *models.GetRolesDTO) ([]*models.RoleFull, error) {
	roles, err := s.repo.GetAllWithNames(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles with names. error: %w", err)
	}
	return roles, nil
}

func (s *RoleService) Get(ctx context.Context, roleName string) (*models.Role, error) {
	role, err := s.repo.Get(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("failed to get role. error: %w", err)
	}
	return role, nil
}

func (s *RoleService) GetByRealm(ctx context.Context, req *models.GetRoleByRealmDTO) (*models.RoleFull, error) {
	data, err := s.repo.GetByRealm(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get role by realm. error: %w", err)
	}
	return data, nil
}

func (s *RoleService) GetWithRealm(ctx context.Context, req *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error) {
	data, err := s.repo.GetWithRealm(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get roles with realm. error: %w", err)
	}
	return data, nil
}

func (s *RoleService) Create(ctx context.Context, role *models.RoleDTO) error {
	if err := s.repo.Create(ctx, role); err != nil {
		return fmt.Errorf("failed to create role. error: %w", err)
	}
	return nil
}

func (s *RoleService) Update(ctx context.Context, role *models.RoleDTO) error {
	if err := s.repo.Update(ctx, role); err != nil {
		return fmt.Errorf("failed to update role. error: %w", err)
	}
	return nil
}

func (s *RoleService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete role. error: %w", err)
	}
	return nil
}
