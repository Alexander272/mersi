package services

import (
	"context"

	"github.com/Alexander272/mersi/backend/internal/models"
)

type UserService struct {
	role Role
	// filter   DefaultFilter
}

func NewUserService(role Role) *UserService {
	return &UserService{
		role: role,
	}
}

type User interface {
	GetRoles(ctx context.Context, req *models.GetUserInfoDTO) (*models.User, error)
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

	// get menu
	rule, err := s.role.Get(ctx, user.Role)
	if err != nil {
		return nil, err
	}

	// get default filters
	// filters, err := s.filter.Get(ctx, &models.GetFilterDTO{SSOId: user.Id, RealmId: req.Realm})
	// if err != nil {
	// 	return nil, err
	// }

	user.Permissions = rule.Rules
	return user, nil
}
