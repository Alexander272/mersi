package services

import (
	"context"
	"fmt"
	"log"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
)

type PermissionService struct {
	enforcer casbin.IEnforcer

	rule     Rule
	role     Role
	realm    Realm
	accesses Accesses
}

type PermissionDeps struct {
	ConfPath string
	Adapter  persist.Adapter
	Rule     Rule
	Role     Role
	Realm    Realm
	Accesses Accesses
}

type Permission interface {
	Register(confPath string, adapter persist.Adapter) error
	Enforce(params ...interface{}) (bool, error)
	ReloadPolicies(ctx context.Context) error
}

func NewPermissionService(deps *PermissionDeps) *PermissionService {
	permission := &PermissionService{
		rule:     deps.Rule,
		role:     deps.Role,
		realm:    deps.Realm,
		accesses: deps.Accesses,
	}

	if err := permission.Register(deps.ConfPath, deps.Adapter); err != nil {
		log.Fatalf("failed to initialize permission service. error: %s", err.Error())
	}
	return permission
}

func (s *PermissionService) Register(path string, adapter persist.Adapter) error {
	//? можно попробовать наследовать конкретных пользователь от роли (пример: g, alice, user, domain1 -
	// вместо alice будет id пользователя, вместо user будет его роль, а вместо domain1 будет id realm)
	// только похоже мне придется роли дублировать для каждого realm что не очень хорошо
	// при таком сценарии я могу добавить id realm в правила

	/*
		пример правил:
		p, reader, domain1, data1, read
		p, user, domain1, data1, write
		p, admin, domain2, data2, read
		p, admin, domain2, data2, write

		g, reader, ,
		g, user, reader, domain1
		g, alice, user, domain1
		g, bob, admin, domain2

		как было раньше:
		p, admin, roles, read
		p, admin, roles, write
		p, admin, realms, write
		p, admin, users, write
		g, reader,
		g, user, reader
		g, metrologist, reader
		g, editor, reader
	*/

	var err error
	s.enforcer, err = casbin.NewEnforcer(path, adapter)
	if err != nil {
		return fmt.Errorf("failed to create enforcer. error: %w", err)
	}

	if err := s.ReloadPolicies(context.Background()); err != nil {
		return fmt.Errorf("failed to prepare policies. error: %w", err)
	}

	// if err = s.enforcer.LoadPolicy(); err != nil {
	// 	return fmt.Errorf("failed to load policy. error: %w", err)
	// }

	return nil
}

func (s *PermissionService) Enforce(params ...interface{}) (bool, error) {
	return s.enforcer.Enforce(params...)
}

func (s *PermissionService) ReloadPolicies(ctx context.Context) error {
	s.enforcer.ClearPolicy()
	if err := s.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy. err: %w", err)
	}

	rules, err := s.rule.GetAll(ctx)
	if err != nil {
		return err
	}
	realms, err := s.realm.Get(ctx, &models.GetRealmsDTO{})
	if err != nil {
		return err
	}
	roles, err := s.role.GetAllWithNames(ctx, &models.GetRolesDTO{})
	if err != nil {
		return err
	}
	accesses, err := s.accesses.GetOriginal(ctx)
	if err != nil {
		return err
	}

	for _, r := range realms {
		for _, m := range rules {
			if _, err := s.enforcer.AddNamedPolicy("p", m.RoleName, r.ID, m.ItemName, m.ItemMethod); err != nil {
				return fmt.Errorf("failed to add policy. error: %w", err)
			}
		}
	}

	for _, r := range roles {
		if len(r.Extends) == 0 {
			if _, err := s.enforcer.AddGroupingPolicy(r.Name, "-", "-"); err != nil {
				return fmt.Errorf("failed to add group policy. error: %w", err)
			}
		}

		for _, v := range r.Extends {
			for _, rm := range realms {
				if _, err := s.enforcer.AddGroupingPolicy(r.Name, v, rm.ID); err != nil {
					return fmt.Errorf("failed to add extended group policy. error: %w", err)
				}
			}
		}
	}

	for _, a := range accesses {
		role := ""
		for _, r := range roles {
			if r.ID == a.RoleID {
				role = r.Name
				break
			}
		}
		if role == "" {
			continue
		}
		if _, err := s.enforcer.AddGroupingPolicy(a.SSO_ID, role, a.RealmID); err != nil {
			return fmt.Errorf("failed to add group (access) policy. error: %w", err)
		}
	}

	if err := s.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy. err: %w", err)
	}
	if err = s.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load policy. error: %w", err)
	}
	return nil
}

func (s *PermissionService) PreparePolicies() error {
	// TODO загрузить в базу правила (добавить все роли и всех пользователей)
	// 	rules, err := s.rule.GetAll(context.Background())
	// 	if err != nil {
	// 		return err
	// 	}

	//TODO получить список realm
	// 	for _, m := range rules {
	// 		line := fmt.Sprintf("p, %s, %s, %s", m.RoleName, m.ItemName, m.ItemMethod)
	// 		logger.Debug("permissions", logger.StringAttr("menu item", line))
	// 		if err := persist.LoadPolicyLine(line, model); err != nil {
	// 			return fmt.Errorf("failed to load policy. error: %w", err)
	// 		}
	// 	}

	return fmt.Errorf("not implemented")
}

func (s *PermissionService) SavePolicy() error {
	if err := s.enforcer.SavePolicy(); err != nil {
		return fmt.Errorf("failed to save policy. err: %w", err)
	}
	return nil
}

func (s *PermissionService) AddPolicy(ptype string, params ...interface{}) error {
	_, err := s.enforcer.AddNamedPolicy(ptype, params...)
	if err != nil {
		return fmt.Errorf("failed to add policy. err: %w", err)
	}
	return nil
}

func (s *PermissionService) RemovePolicy(ptype string, params ...interface{}) error {
	_, err := s.enforcer.RemoveNamedPolicy(ptype, params...)
	if err != nil {
		return fmt.Errorf("failed to remove policy. err: %w", err)
	}
	return nil
}

// type PolicyAdapter struct {
// 	rule Rule
// 	role Role
// }

// func NewPolicyAdapter(rule Rule, role Role) *PolicyAdapter {
// 	return &PolicyAdapter{
// 		rule: rule,
// 		role: role,
// 	}
// }

// type Adapter interface {
// 	LoadPolicy(model model.Model) error
// 	SavePolicy(model model.Model) error
// 	AddPolicy(sec string, ptype string, rule []string) error
// 	RemovePolicy(sec string, ptype string, rule []string) error
// 	RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error
// }

// func (s *PolicyAdapter) LoadPolicy(model model.Model) error {
// 	rules, err := s.rule.GetAll(context.Background())
// 	if err != nil {
// 		return err
// 	}

// 	roles, err := s.role.GetAllWithNames(context.Background(), &models.GetRolesDTO{})
// 	if err != nil {
// 		return err
// 	}

// 	// for _, m := range menu {
// 	// 	for _, mi := range m.MenuItems {
// 	// 		line := fmt.Sprintf("p, %s, %s, %s", m.Role.Name, mi.Name, mi.Method)
// 	// 		logger.Debug("permissions", logger.StringAttr("menu item", line))
// 	// 		if err := persist.LoadPolicyLine(line, model); err != nil {
// 	// 			return fmt.Errorf("failed to load policy. error: %w", err)
// 	// 		}
// 	// 	}

// 	// 	if len(m.Role.Extends) == 0 {
// 	// 		line := fmt.Sprintf("g, %s, ", m.Role.Name)
// 	// 		logger.Debug("permissions", logger.StringAttr("role", line))
// 	// 		if err := persist.LoadPolicyLine(line, model); err != nil {
// 	// 			return fmt.Errorf("failed to load group policy. error: %w", err)
// 	// 		}
// 	// 	}

// 	// 	for _, v := range m.Role.Extends {
// 	// 		line := fmt.Sprintf("g, %s, %s", m.Role.Name, v)
// 	// 		logger.Debug("permissions", logger.StringAttr("extends", line))
// 	// 		if err := persist.LoadPolicyLine(line, model); err != nil {
// 	// 			return fmt.Errorf("failed to load group policy. error: %w", err)
// 	// 		}
// 	// 	}
// 	// }

// 	for _, m := range rules {
// 		line := fmt.Sprintf("p, %s, %s, %s", m.RoleName, m.ItemName, m.ItemMethod)
// 		logger.Debug("permissions", logger.StringAttr("menu item", line))
// 		if err := persist.LoadPolicyLine(line, model); err != nil {
// 			return fmt.Errorf("failed to load policy. error: %w", err)
// 		}
// 	}

// 	for _, r := range roles {
// 		if len(r.Extends) == 0 {
// 			line := fmt.Sprintf("g, %s, ", r.Name)
// 			logger.Debug("permissions", logger.StringAttr("role", line))
// 			if err := persist.LoadPolicyLine(line, model); err != nil {
// 				return fmt.Errorf("failed to load group policy. error: %w", err)
// 			}
// 		}
// 		for _, v := range r.Extends {
// 			line := fmt.Sprintf("g, %s, %s", r.Name, v)
// 			logger.Debug("permissions", logger.StringAttr("extends", line))
// 			if err := persist.LoadPolicyLine(line, model); err != nil {
// 				return fmt.Errorf("failed to load group policy. error: %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }

// // SavePolicy saves all policy rules to the storage.
// func (s *PolicyAdapter) SavePolicy(model model.Model) error {
// 	return nil
// }

// // AddPolicy adds a policy rule to the storage.
// // This is part of the Auto-Save feature.
// func (s *PolicyAdapter) AddPolicy(sec string, ptype string, rule []string) error {
// 	return nil
// }

// // RemovePolicy removes a policy rule from the storage.
// // This is part of the Auto-Save feature.
// func (s *PolicyAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
// 	return nil
// }

// // RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// // This is part of the Auto-Save feature.
// func (a *PolicyAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
// 	return nil
// }
