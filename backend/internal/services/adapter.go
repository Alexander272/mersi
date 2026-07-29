package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/casbin/casbin/v3/model"
	"github.com/casbin/casbin/v3/persist"
)

type adapterService struct {
	perms         AdapterPermissions
	roleHierarchy RoleHierarchy
	users         UserRepository
}

type AdapterDeps struct {
	Permissions   AdapterPermissions
	RoleHierarchy RoleHierarchy
	Users         UserRepository
}

type AdapterPermissions interface {
	LoadPolicy(ctx context.Context) ([]*models.Permission, error)
}

type UserRepository interface {
	LoadPolicy(ctx context.Context) ([]*models.UserRolePolicy, error)
}

func NewAdapter(deps *AdapterDeps) *adapterService {
	return &adapterService{
		perms:         deps.Permissions,
		roleHierarchy: deps.RoleHierarchy,
		users:         deps.Users,
	}
}

type Adapter interface {
	LoadPolicy(model model.Model) error
	SavePolicy(model model.Model) error
	AddPolicy(sec string, ptype string, rule []string) error
	RemovePolicy(sec string, ptype string, rule []string) error
	RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error
}

func (s *adapterService) LoadPolicy(model model.Model) error {
	logger.Info("load policy")

	rootPolicy := "p, root, *, *, *"
	if err := persist.LoadPolicyLine(rootPolicy, model); err != nil {
		return fmt.Errorf("failed to load root policy: %w", err)
	}

	// load permissions
	permissions, err := s.perms.LoadPolicy(context.Background())
	if err != nil {
		return err
	}
	for _, p := range permissions {
		line := fmt.Sprintf("p, %s, %s, %s, %s", p.Role, p.Realm, p.Object, p.Action)
		logger.Debug("permissions", logger.StringAttr("item", line))
		if err := persist.LoadPolicyLine(line, model); err != nil {
			return fmt.Errorf("failed to load permissions policy: %w", err)
		}
	}

	// load role hierarchy
	roles, err := s.roleHierarchy.LoadPolicy(context.Background())
	if err != nil {
		return err
	}
	for _, r := range roles {
		line := fmt.Sprintf("g, %s, %s, %s", r.ParentRole, r.Role, r.Realm)
		logger.Debug("permissions", logger.StringAttr("group", line))
		if err := persist.LoadPolicyLine(line, model); err != nil {
			return fmt.Errorf("failed to load group policy: %w", err)
		}
	}

	// load user roles
	users, err := s.users.LoadPolicy(context.Background())
	if err != nil {
		return err
	}
	for _, u := range users {
		line := fmt.Sprintf("g, %s, %s, %s", u.SSO_ID, u.RoleSlug, u.RealmId)
		logger.Debug("permissions", logger.StringAttr("group", line))
		if err := persist.LoadPolicyLine(line, model); err != nil {
			return fmt.Errorf("failed to load group policy: %w", err)
		}
	}

	return nil
}

func (s *adapterService) SavePolicy(model model.Model) error {
	return nil
}

func (s *adapterService) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (s *adapterService) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

func (s *adapterService) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
