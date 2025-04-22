package services

// type PermissionService struct {
// 	enforcer casbin.IEnforcer
// }

// type Permission interface {
// 	Register(confPath string, menu Menu, role Role) error
// 	Enforce(params ...interface{}) (bool, error)
// }

// func NewPermissionService(confPath string, menu Menu, role Role) *PermissionService {
// 	permission := &PermissionService{}
// 	if err := permission.Register(confPath, menu, role); err != nil {
// 		log.Fatalf("failed to initialize permission service. error: %s", err.Error())
// 	}
// 	return permission
// }

// func (s *PermissionService) Register(path string, menu Menu, role Role) error {
// 	var err error
// 	adapter := NewPolicyAdapter(menu, role)

// 	s.enforcer, err = casbin.NewEnforcer(path, adapter)
// 	if err != nil {
// 		return fmt.Errorf("failed to create enforcer. error: %w", err)
// 	}

// 	if err = s.enforcer.LoadPolicy(); err != nil {
// 		return fmt.Errorf("failed to load policy. error: %w", err)
// 	}

// 	return nil
// }

// func (s *PermissionService) Enforce(params ...interface{}) (bool, error) {
// 	return s.enforcer.Enforce(params...)
// }

// type PolicyAdapter struct {
// 	menu Menu
// 	role Role
// }

// func NewPolicyAdapter(menu Menu, role Role) *PolicyAdapter {
// 	return &PolicyAdapter{
// 		menu: menu,
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
// 	menu, err := s.menu.GetAll(context.Background())
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

// 	for _, m := range menu {
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
