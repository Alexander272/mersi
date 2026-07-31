package services

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestClaimsToUser_ExtractsRoleByServicePrefix(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":               "user-1",
		"preferred_username": "alice",
		"realm_access": map[string]interface{}{
			"roles": []interface{}{"app_user", "sia_admin"},
		},
	}

	user, err := claimsToUser(claims, "sia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "user-1" {
		t.Fatalf("expected id user-1, got %s", user.ID)
	}
	if user.Name != "alice" {
		t.Fatalf("expected name alice, got %s", user.Name)
	}
	if user.Role != "admin" {
		t.Fatalf("expected role admin (from sia_admin), got %s", user.Role)
	}
}

func TestClaimsToUser_PicksFirstMatchingRole(t *testing.T) {
	claims := jwt.MapClaims{
		"sub": "user-1",
		"realm_access": map[string]interface{}{
			"roles": []interface{}{"sia_metrologist", "sia_admin"},
		},
	}

	user, err := claimsToUser(claims, "sia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Role != "metrologist" {
		t.Fatalf("expected first matching role metrologist, got %s", user.Role)
	}
}

func TestClaimsToUser_NoServiceRoleLeavesRoleEmpty(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":               "user-1",
		"preferred_username": "alice",
		"realm_access": map[string]interface{}{
			"roles": []interface{}{"app_user"},
		},
	}

	user, err := claimsToUser(claims, "sia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Role != "" {
		t.Fatalf("expected empty role, got %s", user.Role)
	}
}

func TestClaimsToUser_MissingRealmAccess(t *testing.T) {
	user, err := claimsToUser(jwt.MapClaims{"sub": "user-1"}, "sia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Role != "" || user.ID != "user-1" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestClaimsToUser_WrongClaimTypesIgnored(t *testing.T) {
	claims := jwt.MapClaims{
		"sub":               123,
		"preferred_username": 456,
		"realm_access":      "not-a-map",
	}

	user, err := claimsToUser(claims, "sia")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != "" || user.Name != "" || user.Role != "" {
		t.Fatalf("expected empty user fields, got %+v", user)
	}
}
