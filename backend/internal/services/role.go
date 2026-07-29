package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/access"
	"github.com/Alexander272/mersi/backend/internal/events"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"golang.org/x/sync/errgroup"
)

type RoleService struct {
	repo      repository.Role
	hierarchy RoleHierarchy
	perms     RolePermissions
	eventBus  *events.PolicyEventManager
	tm        TransactionManager
}

type RoleDeps struct {
	Repo        repository.Role
	Hierarchy   RoleHierarchy
	Permissions RolePermissions
	EventBus    *events.PolicyEventManager
	TM          TransactionManager
}

func NewRoleService(deps *RoleDeps) *RoleService {
	return &RoleService{
		repo:      deps.Repo,
		hierarchy: deps.Hierarchy,
		perms:     deps.Permissions,
		eventBus:  deps.EventBus,
		tm:        deps.TM,
	}
}

type Role interface {
	GetAll(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	GetAllWithNames(context.Context, *models.GetRolesDTO) ([]*models.RoleFull, error)
	Get(context.Context, string) (*models.Role, error)
	GetRules(ctx context.Context, roleName string) ([]string, error)
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

func (s *RoleService) GetRules(ctx context.Context, roleName string) ([]string, error) {
	return s.repo.GetRules(ctx, roleName)
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

// New permission system methods
func (s *RoleService) GetOne(ctx context.Context, req *models.GetRoleDTO) (*models.Role, error) {
	data, err := s.repo.GetOne(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return data, nil
}

func (s *RoleService) GetAllRoles(ctx context.Context) ([]*models.Role, error) {
	data, err := s.repo.GetAllRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}
	return data, nil
}

func (s *RoleService) GetOneWithPermissions(ctx context.Context, req *models.GetRoleDTO) (*models.RoleWithPerms, error) {
	data, err := s.GetOne(ctx, req)
	if err != nil {
		return nil, err
	}

	children, err := s.hierarchy.GetDirectChildren(ctx, &models.GetRolesInheritance{Roles: []string{data.Slug}})
	if err != nil {
		return nil, err
	}

	perms, err := s.GetPermissionsGrouped(ctx, req)
	if err != nil {
		return nil, err
	}

	res := &models.RoleWithPerms{
		Role:     *data,
		Inherits: children[data.Slug],
		Perms:    perms,
	}
	return res, nil
}

func (s *RoleService) GetWithStats(ctx context.Context) ([]*models.RoleWithStats, error) {
	roles, err := s.GetAllRoles(ctx)
	if err != nil {
		return nil, err
	}

	roleIDs := make([]string, 0, len(roles))
	slugs := make([]string, 0, len(roles))
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
		slugs = append(slugs, role.Slug)
	}

	var (
		userCounts  map[string]int
		permsCounts map[string]models.PermsWithCount
		children    map[string][]string
		descendants map[string][]string
	)

	g, asyncCtx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		var err error
		children, err = s.hierarchy.GetDirectChildren(asyncCtx, &models.GetRolesInheritance{Roles: slugs})
		return err
	})
	g.Go(func() error {
		var err error
		descendants, err = s.hierarchy.GetRoleDescendants(asyncCtx, &models.GetRolesInheritance{Roles: slugs})
		return err
	})
	g.Go(func() error {
		var err error
		userCounts, err = s.repo.GetUserCount(asyncCtx, roleIDs)
		if err != nil {
			return fmt.Errorf("failed to get user count: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	if len(descendants) > 0 {
		permsCounts, err = s.perms.CountForAll(ctx, descendants)
		if err != nil {
			return nil, err
		}
	}

	result := make([]*models.RoleWithStats, 0, len(roles))
	for _, r := range roles {
		result = append(result, &models.RoleWithStats{
			Role:       *r,
			Children:   children[r.Slug],
			UserCount:  userCounts[r.ID],
			PermsCount: permsCounts[r.Slug],
		})
	}

	return result, nil
}

func (s *RoleService) GetPermissionsGrouped(ctx context.Context, req *models.GetRoleDTO) ([]*models.RolePermissionsGrouped, error) {
	allPerms, err := s.perms.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	assigned, err := s.perms.GetRolePermissions(ctx, 	req.ID)
	if err != nil {
		return nil, err
	}

	inherited, err := s.perms.GetInherited(ctx, 	req.ID)
	if err != nil {
		return nil, err
	}
	logger.Debug("permissions", logger.IntAttr("assigned", len(assigned)), logger.IntAttr("inherited", len(inherited)))

	grouped := make(map[string][]*models.RolePermissionItem)
	for _, p := range allPerms {
		item := &models.RolePermissionItem{
			PermissionId: p.Id,
			Object:       p.Object,
			Action:       p.Action,
			IsAssigned:   assigned[p.Id],
			IsInherited:  inherited[p.Id],
		}
		grouped[p.Object] = append(grouped[p.Object], item)
	}

	result := make([]*models.RolePermissionsGrouped, 0, len(grouped))
	for group, resources := range grouped {
		res, _ := access.Reg.GetBySlug(access.ResourceSlug(group))
		result = append(result, &models.RolePermissionsGrouped{
			Group:     group,
			Title:     res.Name,
			Resources: resources,
		})
	}

	return result, nil
}

func (s *RoleService) IsExists(ctx context.Context, roleName string) (bool, error) {
	data, err := s.repo.IsExists(ctx, roleName)
	if err != nil {
		return false, fmt.Errorf("failed to check if role exists: %w", err)
	}
	return data, nil
}

func (s *RoleService) CreateRole(ctx context.Context, dto *models.RoleDTO) error {
	return s.tm.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		err := s.repo.CreateRole(ctx, tx, dto)
		if err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}

		inheritIDs := make([]string, 0, len(dto.Inherits))
		if len(dto.Inherits) > 0 {
			roleIDs, err := s.repo.GetIdsBySlugs(ctx, dto.Inherits)
			if err != nil {
				return fmt.Errorf("failed to get role IDs: %w", err)
			}
			for _, slug := range dto.Inherits {
				if id, ok := roleIDs[slug]; ok {
					inheritIDs = append(inheritIDs, id)
				} else {
					return fmt.Errorf("parent role not found: %s", slug)
				}
			}

			if err := s.hierarchy.AddInheritances(ctx, tx, 	dto.ID, inheritIDs); err != nil {
				return fmt.Errorf("failed to add inheritances: %w", err)
			}
		}

		if err := s.repo.AssignPermissions(ctx, tx, 	dto.ID, dto.Permissions); err != nil {
			return fmt.Errorf("failed to assign permissions: %w", err)
		}

		s.eventBus.Notify(events.PolicyEvent{})
		return nil
	})
}

func (s *RoleService) UpdateRole(ctx context.Context, dto *models.RoleDTO) error {
	oldRole, err := s.GetOne(ctx, &models.GetRoleDTO{ID: dto.ID})
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	if !oldRole.IsEditable {
		return fmt.Errorf("error while updating role %s: %w", dto.Name, models.ErrRoleNotEditable)
	}

	return s.tm.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		err := s.repo.UpdateRole(ctx, tx, dto)
		if err != nil {
			return fmt.Errorf("failed to update role: %w", err)
		}

		if len(dto.Inherits) > 0 {
			currentInherits, err := s.hierarchy.GetRoleDescendants(ctx, &models.GetRolesInheritance{Roles: []string{dto.Slug}})
			if err != nil {
				return fmt.Errorf("failed to get current inherits: %w", err)
			}

			currentSlugs := currentInherits[dto.Slug]
			toAdd := make([]string, 0, len(dto.Inherits))
			toRemove := make([]string, 0, len(currentSlugs))

			currentMap := make(map[string]bool)
			for _, s := range currentSlugs {
				currentMap[s] = true
			}

			newMap := make(map[string]bool)
			for _, s := range dto.Inherits {
				newMap[s] = true
			}

			for _, s := range dto.Inherits {
				if !currentMap[s] {
					toAdd = append(toAdd, s)
				}
			}
			for _, s := range currentSlugs {
				if !newMap[s] {
					toRemove = append(toRemove, s)
				}
			}

			if len(toAdd) > 0 {
				addIDs, err := s.repo.GetIdsBySlugs(ctx, toAdd)
				if err != nil {
					return fmt.Errorf("failed to get role IDs: %w", err)
				}
				inheritIDs := make([]string, 0, len(toAdd))
				for _, s := range toAdd {
					if id, ok := addIDs[s]; ok {
						inheritIDs = append(inheritIDs, id)
					}
				}

				logger.Debug("inheritances",
					logger.IntAttr("count", len(inheritIDs)),
					logger.StringAttr("role", 	dto.ID),
					logger.StringAttr("parents", strings.Join(dto.Inherits, ", ")))

				if err := s.hierarchy.AddInheritances(ctx, tx, 	dto.ID, inheritIDs); err != nil {
					return err
				}
			}

			if len(toRemove) > 0 {
				removeIDs, err := s.repo.GetIdsBySlugs(ctx, toRemove)
				if err != nil {
					return fmt.Errorf("failed to get role IDs: %w", err)
				}
				parentIDs := make([]string, 0, len(toRemove))
				for _, s := range toRemove {
					if id, ok := removeIDs[s]; ok {
						parentIDs = append(parentIDs, id)
					}
				}
				if err := s.hierarchy.RemoveInheritances(ctx, tx, 	dto.ID, parentIDs); err != nil {
					return err
				}
			}
		}

		if len(dto.Permissions) > 0 {
			if err := s.perms.ReplacePermissions(ctx, tx, 	dto.ID, dto.Permissions); err != nil {
				return err
			}
		}

		s.eventBus.Notify(events.PolicyEvent{})
		return nil
	})
}

func (s *RoleService) DeleteRole(ctx context.Context, dto *models.DeleteRoleDTO) error {
	role, err := s.repo.GetOne(ctx, &models.GetRoleDTO{ID: dto.ID})
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}

	if role.IsSystem || !role.IsEditable {
		return models.ErrReservedRole
	}

	err = s.repo.DeleteRole(ctx, nil, dto)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	s.eventBus.Notify(events.PolicyEvent{})

	return nil
}

func (s *RoleService) AssignPermission(ctx context.Context, dto *models.RolePermissionDTO) error {
	return s.tm.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		err := s.repo.AssignPermission(ctx, tx, dto)
		if err != nil {
			return fmt.Errorf("failed to assign permission: %w", err)
		}
		return nil
	})
}

func (s *RoleService) DeletePermission(ctx context.Context, dto *models.RolePermissionDTO) error {
	return s.tm.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		err := s.repo.DeletePermission(ctx, tx, dto)
		if err != nil {
			return fmt.Errorf("failed to delete permission: %w", err)
		}
		return nil
	})
}

func (s *RoleService) SetPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return s.tm.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.perms.ReplacePermissions(ctx, tx, roleID, permissionIDs); err != nil {
			return fmt.Errorf("failed to replace permissions: %w", err)
		}
		return nil
	})
}
