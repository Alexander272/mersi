package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
	"github.com/Alexander272/mersi/backend/pkg/auth"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/Nerzal/gocloak/v13"
)

type UserService struct {
	repo      repository.Users
	txManager TransactionManager
	keycloak  *auth.KeycloakClient
	role      Role
}

type UsersDeps struct {
	Repo      repository.Users
	TxManager TransactionManager
	Keycloak  *auth.KeycloakClient
	Role      Role
}

func NewUserService(deps *UsersDeps) *UserService {
	return &UserService{
		repo:      deps.Repo,
		txManager: deps.TxManager,
		keycloak:  deps.Keycloak,
		role:      deps.Role,
	}
}

type Users interface {
	GetAll(ctx context.Context) ([]*models.UserData, error)
	GetByAccess(ctx context.Context, req *models.GetByAccessDTO) ([]*models.UserData, error)
	GetByRealm(ctx context.Context, req *models.GetByRealmDTO) ([]*models.UserData, error)
	GetById(ctx context.Context, id string) (*models.UserData, error)
	GetBySSOId(ctx context.Context, id string) (*models.UserData, error)
	GetRoles(ctx context.Context, req *models.GetUserInfoDTO) (*models.User, error)
	GetPermissions(ctx context.Context, userId string) (map[string][]string, error)
	LoadPolicy(ctx context.Context) ([]*models.UserRolePolicy, error)
	Sync(ctx context.Context) error
	Create(ctx context.Context, dto *models.UserData) error
	CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserData) error
	Update(ctx context.Context, dto *models.UserData) error
	UpdateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserData) error
	Delete(ctx context.Context, id string) error
	DeleteSeveral(ctx context.Context, tx postgres.Tx, ids []string) error
}

func (s *UserService) GetAll(ctx context.Context) ([]*models.UserData, error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users. error: %w", err)
	}
	return data, nil
}

func (s *UserService) GetByAccess(ctx context.Context, req *models.GetByAccessDTO) ([]*models.UserData, error) {
	data, err := s.repo.GetByAccess(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by access. error: %w", err)
	}
	return data, nil
}

func (s *UserService) GetByRealm(ctx context.Context, req *models.GetByRealmDTO) ([]*models.UserData, error) {
	data, err := s.repo.GetByRealm(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get users by realm. error: %w", err)
	}
	return data, nil
}

func (s *UserService) GetById(ctx context.Context, id string) (*models.UserData, error) {
	data, err := s.repo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by id. error: %w", err)
	}
	return data, nil
}

func (s *UserService) GetBySSOId(ctx context.Context, id string) (*models.UserData, error) {
	data, err := s.repo.GetBySSOId(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by sso_id. error: %w", err)
	}
	return data, nil
}

func (s *UserService) GetRoles(ctx context.Context, req *models.GetUserInfoDTO) (*models.User, error) {
	roles, err := s.role.GetWithRealm(ctx, &models.GetRoleByRealmDTO{UserID: req.UserID})
	if err != nil {
		return nil, err
	}
	user := &models.User{ID: req.UserID}

	user.Roles = roles
	if req.Realm == "" {
		user.Role = roles[0].Name
		req.Realm = roles[0].RealmId
	} else {
		for _, r := range roles {
			if r.RealmId == req.Realm {
				user.Role = r.Name
				break
			}
		}
	}

	return user, nil
}

func (s *UserService) Sync(ctx context.Context) error {
	logger.Info("Sync users started")

	token, err := s.keycloak.Login(ctx)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	// 1. Быстрый поиск ID группы
	groups, err := s.keycloak.Client.GetGroups(ctx, token.AccessToken, s.keycloak.Realm, gocloak.GetGroupsParams{
		Search: gocloak.StringP("mersi"), // Фильтруем на стороне Keycloak
	})
	if err != nil {
		return fmt.Errorf("failed to get groups: %w", err)
	}

	var groupID string
	for _, g := range groups {
		if g.Name != nil && *g.Name == "mersi" {
			groupID = *g.ID
			break
		}
	}
	if groupID == "" {
		return fmt.Errorf("group 'mersi' not found")
	}

	// 2. Получаем активных пользователей из Keycloak
	keycloakUsers, err := s.keycloak.Client.GetGroupMembers(ctx, token.AccessToken, s.keycloak.Realm, groupID, gocloak.GetGroupsParams{Max: gocloak.IntP(1000)})
	if err != nil {
		return fmt.Errorf("failed to get group members: %w", err)
	}

	if len(keycloakUsers) == 0 {
		return fmt.Errorf("group 'mersi' is empty")
	}

	// Пред-аллокация для пачки данных из Keycloak
	kcDataMap := make(map[string]*models.UserData, len(keycloakUsers))
	for _, u := range keycloakUsers {
		if u.Enabled != nil && !*u.Enabled {
			continue
		}

		userData := s.mapToUserData(u)
		kcDataMap[userData.SSO_ID] = userData
	}

	// 3. Получаем текущих пользователей из нашей БД
	dbUsers, err := s.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch DB users: %w", err)
	}

	toCreate := make([]*models.UserData, 0)
	toUpdate := make([]*models.UserData, 0)
	toDelete := make([]string, 0)

	// 4. Основной цикл синхронизации
	dbUserMap := make(map[string]*models.UserData, len(dbUsers))

	for _, dbU := range dbUsers {
		dbUserMap[dbU.SSO_ID] = dbU

		if kcData, exists := kcDataMap[dbU.SSO_ID]; exists {
			// Проверяем, нужно ли реально обновлять (DeepEqual или по полям)
			if s.isChanged(dbU, kcData) {
				toUpdate = append(toUpdate, kcData)
			}
			// Удаляем из мапы Keycloak, чтобы там остались только "новые"
			delete(kcDataMap, dbU.SSO_ID)
		} else {
			// Если в Keycloak нет, а в БД есть — на удаление
			toDelete = append(toDelete, dbU.SSO_ID)
		}
	}

	// Все, кто остались в kcDataMap — новые
	for _, newU := range kcDataMap {
		toCreate = append(toCreate, newU)
	}

	// 5. Выполнение операций (Batch processing)
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if len(toCreate) > 0 {
			if err := s.CreateSeveral(ctx, tx, toCreate); err != nil {
				return err
			}
		}
		if len(toUpdate) > 0 {
			if err := s.UpdateSeveral(ctx, tx, toUpdate); err != nil {
				return err
			}
		}
		if len(toDelete) > 0 {
			if err := s.DeleteSeveral(ctx, tx, toDelete); err != nil {
				return err
			}
		}

		logger.Info("Sync finished",
			"created", len(toCreate),
			"updated", len(toUpdate),
			"deleted", len(toDelete))
		return nil
	})
}

// Вспомогательная функция для маппинга (убирает дублирование nil-проверок)
func (s *UserService) mapToUserData(u *gocloak.User) *models.UserData {
	return &models.UserData{
		SSO_ID:    s.nonNil(u.ID),
		Username:  s.nonNil(u.Username),
		Email:     s.nonNil(u.Email),
		FirstName: s.nonNil(u.FirstName),
		LastName:  s.nonNil(u.LastName),
	}
}

func (s *UserService) nonNil(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

// Функция проверки изменений, чтобы не дергать БД зря
func (s *UserService) isChanged(old, new *models.UserData) bool {
	return old.Username != new.Username ||
		old.Email != new.Email ||
		old.FirstName != new.FirstName ||
		old.LastName != new.LastName
}

func (s *UserService) Create(ctx context.Context, dto *models.UserData) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create user. error: %w", err)
	}
	return nil
}

func (s *UserService) CreateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserData) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.CreateSeveral(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create few users. error: %w", err)
	}
	return nil
}

func (s *UserService) Update(ctx context.Context, dto *models.UserData) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update user. error: %w", err)
	}
	return nil
}

func (s *UserService) UpdateSeveral(ctx context.Context, tx postgres.Tx, dto []*models.UserData) error {
	if len(dto) == 0 {
		return nil
	}
	if err := s.repo.UpdateSeveral(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to update few users. error: %w", err)
	}
	return nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user. error: %w", err)
	}
	return nil
}

func (s *UserService) DeleteSeveral(ctx context.Context, tx postgres.Tx, ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.DeleteSeveral(ctx, tx, ids); err != nil {
		return fmt.Errorf("failed to delete few users. error: %w", err)
	}
	return nil
}

func (s *UserService) LoadPolicy(ctx context.Context) ([]*models.UserRolePolicy, error) {
	data, err := s.repo.LoadPolicy(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load user policy: %w", err)
	}
	return data, nil
}

func (s *UserService) GetPermissions(ctx context.Context, userId string) (map[string][]string, error) {
	roles, err := s.role.GetWithRealm(ctx, &models.GetRoleByRealmDTO{UserID: userId})
	if err != nil {
		return nil, err
	}

	permissions := make(map[string][]string)
	for _, r := range roles {
		rules, err := s.role.GetRules(ctx, r.Name)
		if err != nil {
			continue
		}
		permissions[r.RealmId] = rules
	}
	return permissions, nil
}
