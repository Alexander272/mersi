package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/models"
)

func TestUserGetRoles_NoRolesReturnsForbidden(t *testing.T) {
	role := &fakeRoleSvc{
		getWithRealmFn: func(ctx context.Context, dto *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error) {
			return []*models.RoleWithRealm{}, nil
		},
	}
	svc := &UserService{role: role}

	_, err := svc.GetRoles(context.Background(), &models.GetUserInfoDTO{UserID: "user-1"})
	if !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestUserGetRoles_EmptyRealmUsesHighestRole(t *testing.T) {
	role := &fakeRoleSvc{
		getWithRealmFn: func(ctx context.Context, dto *models.GetRoleByRealmDTO) ([]*models.RoleWithRealm, error) {
			return []*models.RoleWithRealm{
				{Name: "admin", RealmId: "realm-1"},
				{Name: "viewer", RealmId: "realm-2"},
			}, nil
		},
		getFn: func(ctx context.Context, name string) (*models.Role, error) {
			if name != "admin" {
				t.Fatalf("expected role admin, got %q", name)
			}
			return &models.Role{Name: "admin", Rules: []string{"b", "a"}}, nil
		},
	}
	svc := &UserService{role: role}

	user, err := svc.GetRoles(context.Background(), &models.GetUserInfoDTO{UserID: "user-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Role != "admin" {
		t.Fatalf("expected role admin, got %q", user.Role)
	}
	if len(user.Permissions) != 2 || user.Permissions[0] != "a" || user.Permissions[1] != "b" {
		t.Fatalf("expected sorted permissions, got %v", user.Permissions)
	}
}
