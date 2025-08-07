package most

import (
	"context"
	"fmt"

	"github.com/mattermost/mattermost-server/v6/model"
)

type ChannelService struct {
	client  *model.Client4
	botName string
}

func NewChannelService(client *model.Client4, botName string) *ChannelService {
	return &ChannelService{
		client:  client,
		botName: botName,
	}
}

type Channel interface {
	Create(ctx context.Context, userId string) (string, error)
}

// create new direct chanel
func (s *ChannelService) Create(ctx context.Context, userId string) (string, error) {
	bot, _, err := s.client.GetUserByUsername(s.botName, "")
	if err != nil {
		return "", fmt.Errorf("failed to get bot. error: %w", err)
	}

	channel, _, err := s.client.CreateDirectChannel(bot.Id, userId)
	if err != nil {
		return "", fmt.Errorf("failed to create direct channel. error: %w", err)
	}

	return channel.Id, nil
}
