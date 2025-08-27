package services

import (
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
)

type ChannelService struct {
	repo repository.Channel
}

func NewChannelService(repo repository.Channel) *ChannelService {
	return &ChannelService{
		repo: repo,
	}
}

type Channel interface {
	GetAll(ctx context.Context) ([]*models.Channel, error)
	Create(ctx context.Context, channel *models.Channel) error
	Update(ctx context.Context, channel *models.Channel) error
	Delete(ctx context.Context, id string) error
}

func (s *ChannelService) GetAll(ctx context.Context) ([]*models.Channel, error) {
	data, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all channels. error: %w", err)
	}
	return data, nil
}

func (s *ChannelService) Create(ctx context.Context, channel *models.Channel) error {
	if err := s.repo.Create(ctx, channel); err != nil {
		return fmt.Errorf("failed to create channel. error: %w", err)
	}
	return nil
}

func (s *ChannelService) Update(ctx context.Context, channel *models.Channel) error {
	if err := s.repo.Update(ctx, channel); err != nil {
		return fmt.Errorf("failed to update channel. error: %w", err)
	}
	return nil
}

func (s *ChannelService) Delete(ctx context.Context, id string) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete channel. error: %w", err)
	}
	return nil
}
