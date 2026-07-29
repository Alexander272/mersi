package services

import (
	"context"
	"fmt"
	"log"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/events"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/casbin/casbin/v3"
)

type accessPoliciesService struct {
	enforcer casbin.IEnforcer
	adapter  Adapter
	cache    SessionCacher
	eventBus *events.PolicyEventManager
	conf     *config.CasbinConfig
}

type PoliciesDeps struct {
	Conf     *config.CasbinConfig
	Adapter  Adapter
	EventBus *events.PolicyEventManager
	Cache    SessionCacher
}

type AccessPolices interface {
	Enforce(params ...interface{}) (bool, error)
	ReloadPolicies() error
	GetPolicies(userId string, realm string) (*models.PolicyResult, error)
}

func NewAccessPoliciesService(deps *PoliciesDeps) *accessPoliciesService {
	s := &accessPoliciesService{
		adapter:  deps.Adapter,
		cache:    deps.Cache,
		eventBus: deps.EventBus,
		conf:     deps.Conf,
	}

	if err := s.register(); err != nil {
		log.Fatalf("failed to initialize permission service. error: %s", err.Error())
	}

	go s.handlePolicyUpdates()

	return s
}

func (s *accessPoliciesService) handlePolicyUpdates() {
	ch := s.eventBus.Subscribe()
	for range ch {
		if err := s.ReloadPolicies(); err != nil {
			logger.Error("failed to reload policies", logger.ErrAttr(err))
		}
		s.cache.Flush(context.Background())
	}
}

func (s *accessPoliciesService) register() error {
	var err error
	s.enforcer, err = casbin.NewEnforcer(s.conf.Path, s.adapter)
	if err != nil {
		return fmt.Errorf("failed to create enforcer: %w", err)
	}

	if err := s.ReloadPolicies(); err != nil {
		return fmt.Errorf("failed to prepare policies: %w", err)
	}

	return nil
}

func (s *accessPoliciesService) Enforce(params ...interface{}) (bool, error) {
	return s.enforcer.Enforce(params...)
}

func (s *accessPoliciesService) ReloadPolicies() error {
	s.enforcer.ClearPolicy()
	if err := s.enforcer.LoadPolicy(); err != nil {
		return fmt.Errorf("failed to load policy: %w", err)
	}
	logger.Info("policies reloaded")
	return nil
}

func (s *accessPoliciesService) GetPolicies(userId string, realm string) (*models.PolicyResult, error) {
	roles := s.enforcer.GetRolesForUserInDomain(userId, realm)

	if len(roles) == 0 {
		return &models.PolicyResult{Perms: []string{}}, nil
	}

	permissions := make(map[string]bool)
	for _, role := range roles {
		if role == "root" {
			return &models.PolicyResult{Perms: []string{"read", "write", "delete"}}, nil
		}
		perms := s.enforcer.GetPermissionsForUserInDomain(role, realm)
		for _, p := range perms {
			if len(p) > 3 {
				permissions[p[3]] = true
			}
		}
	}

	perms := make([]string, 0, len(permissions))
	for k := range permissions {
		perms = append(perms, k)
	}

	return &models.PolicyResult{Perms: perms}, nil
}
