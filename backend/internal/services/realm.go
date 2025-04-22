package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type RealmService struct {
	repo repository.Realm
}

func NewRealmService(repo repository.Realm) *RealmService {
	return &RealmService{
		repo: repo,
	}
}

type Realm interface {
	Get(ctx context.Context, req *models.GetRealmsDTO) ([]*models.Realm, error)
	GetById(ctx context.Context, req *models.GetRealmByIdDTO) (*models.Realm, error)
	GetByUser(ctx context.Context, req *models.GetRealmByUserDTO) ([]*models.Realm, error)
	Choose(ctx context.Context, dto *models.ChooseRealmDTO) (*models.User, error)
	Create(ctx context.Context, dto *models.RealmDTO) error
	Update(ctx context.Context, dto *models.RealmDTO) error
	Delete(ctx context.Context, dto *models.DeleteRealmDTO) error
}

func (s *RealmService) Get(ctx context.Context, req *models.GetRealmsDTO) ([]*models.Realm, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get realms. error: %w", err)
	}
	return data, nil
}

func (s *RealmService) GetById(ctx context.Context, req *models.GetRealmByIdDTO) (*models.Realm, error) {
	data, err := s.repo.GetById(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get realm by id. error: %w", err)
	}
	return data, nil
}

func (s *RealmService) GetByUser(ctx context.Context, req *models.GetRealmByUserDTO) ([]*models.Realm, error) {
	data, err := s.repo.GetByUser(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get realm by user. error: %w", err)
	}
	return data, nil
}

func (s *RealmService) Choose(ctx context.Context, dto *models.ChooseRealmDTO) (*models.User, error) {
	// get menu
	// menu, err := s.role.Get(ctx, dto.Role)
	// if err != nil {
	// 	return nil, err
	// }

	// // get default filters
	// filters, err := s.filter.Get(ctx, &models.GetFilterDTO{SSOId: dto.UserId, RealmId: dto.RealmId})
	// if err != nil {
	// 	return nil, err
	// }

	// user := &models.User{
	// 	Role:    dto.Role,
	// 	Menu:    menu.Menu,
	// 	Filters: filters,
	// }
	// return user, nil

	return nil, fmt.Errorf("not implemented")
}

func (s *RealmService) Create(ctx context.Context, dto *models.RealmDTO) error {
	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create realm. error: %w", err)
	}
	return nil
}

func (s *RealmService) Update(ctx context.Context, dto *models.RealmDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update realm. error: %w", err)
	}
	return nil
}

func (s *RealmService) Delete(ctx context.Context, dto *models.DeleteRealmDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete realm. error: %w", err)
	}
	return nil
}
