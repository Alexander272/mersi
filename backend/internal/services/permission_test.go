package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/casbin/casbin/v2/persist/file-adapter"
)

func TestReloadPolicies_EnforcesAccess(t *testing.T) {
	dir := t.TempDir()
	policyFile := filepath.Join(dir, "policy.csv")
	if err := os.WriteFile(policyFile, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	adapter := fileadapter.NewAdapter(policyFile)
	modelPath := filepath.Join("..", "..", "configs", "privacy.conf")

	svc := NewPermissionService(&PermissionDeps{
		ConfPath: modelPath,
		Adapter:  adapter,
		Rule: &fakeRuleSvc{
			getAllFn: func(ctx context.Context) ([]*models.Rule, error) {
				return []*models.Rule{{RoleName: "reader", ItemName: "data1", ItemMethod: "read"}}, nil
			},
		},
		Realm: &fakeRealmSvc{
			getFn: func(ctx context.Context, dto *models.GetRealmsDTO) ([]*models.Realm, error) {
				return []*models.Realm{{ID: "r1"}, {ID: "r2"}}, nil
			},
		},
		Role: &fakeRoleSvc{
			getAllWithNamesFn: func(ctx context.Context, dto *models.GetRolesDTO) ([]*models.RoleFull, error) {
				return []*models.RoleFull{
					{ID: "reader-id", Name: "reader", Extends: nil},
					{ID: "user-id", Name: "user", Extends: []string{"reader"}},
				}, nil
			},
		},
		Accesses: &fakeAccessesSvc{
			getOriginalFn: func(ctx context.Context) ([]*models.AccessesDTO, error) {
				return []*models.AccessesDTO{{SSO_ID: "alice", RoleID: "user-id", RealmID: "r1"}}, nil
			},
		},
	})

	cases := []struct {
		name   string
		sub    string
		dom    string
		obj    string
		act    string
		expect bool
	}{
		{name: "alice reader in realm r1", sub: "alice", dom: "r1", obj: "data1", act: "read", expect: true},
		{name: "alice denied write", sub: "alice", dom: "r1", obj: "data1", act: "write", expect: false},
		{name: "alice denied in realm r2", sub: "alice", dom: "r2", obj: "data1", act: "read", expect: false},
		{name: "reader role in r1", sub: "reader", dom: "r1", obj: "data1", act: "read", expect: true},
		{name: "unknown user denied", sub: "bob", dom: "r1", obj: "data1", act: "read", expect: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := svc.Enforce(tc.sub, tc.dom, tc.obj, tc.act)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ok != tc.expect {
				t.Fatalf("expected %v, got %v", tc.expect, ok)
			}
		})
	}
}

func TestReloadPolicies_ReloadsFromFreshData(t *testing.T) {
	dir := t.TempDir()
	policyFile := filepath.Join(dir, "policy.csv")
	if err := os.WriteFile(policyFile, []byte{}, 0o600); err != nil {
		t.Fatal(err)
	}
	adapter := fileadapter.NewAdapter(policyFile)
	modelPath := filepath.Join("..", "..", "configs", "privacy.conf")

	data := []*models.AccessesDTO{{SSO_ID: "alice", RoleID: "r-id", RealmID: "r1"}}
	accesses := &fakeAccessesSvc{
		getOriginalFn: func(ctx context.Context) ([]*models.AccessesDTO, error) { return data, nil },
	}

	svc := NewPermissionService(&PermissionDeps{
		ConfPath: modelPath,
		Adapter:  adapter,
		Rule: &fakeRuleSvc{
			getAllFn: func(ctx context.Context) ([]*models.Rule, error) {
				return []*models.Rule{{RoleName: "reader", ItemName: "data1", ItemMethod: "read"}}, nil
			},
		},
		Realm: &fakeRealmSvc{
			getFn: func(ctx context.Context, dto *models.GetRealmsDTO) ([]*models.Realm, error) {
				return []*models.Realm{{ID: "r1"}}, nil
			},
		},
		Role: &fakeRoleSvc{
			getAllWithNamesFn: func(ctx context.Context, dto *models.GetRolesDTO) ([]*models.RoleFull, error) {
				return []*models.RoleFull{{ID: "r-id", Name: "reader", Extends: nil}}, nil
			},
		},
		Accesses: accesses,
	})

	ok, err := svc.Enforce("alice", "r1", "data1", "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected access granted")
	}

	data[0].SSO_ID = "bob"
	if err := svc.ReloadPolicies(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok, _ := svc.Enforce("alice", "r1", "data1", "read"); ok {
		t.Fatal("expected alice access revoked after reload")
	}
	if ok, _ := svc.Enforce("bob", "r1", "data1", "read"); !ok {
		t.Fatal("expected bob access granted after reload")
	}
}
