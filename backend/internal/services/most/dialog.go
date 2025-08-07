package most

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/mattermost/mattermost-server/v6/model"
)

type DialogService struct {
	client *model.Client4
}

func NewDialogService(client *model.Client4) *DialogService {
	return &DialogService{
		client: client,
	}
}

type Dialog interface {
	Open(ctx context.Context, action *models.PostAction) error
}

func (s *DialogService) Open(ctx context.Context, action *models.PostAction) error {
	state := []string{
		fmt.Sprintf("PostId:%s", action.PostId),
		action.Context.State,
	}

	dialogData := model.Dialog{
		CallbackId:       action.Context.CallbackId,
		Title:            action.Context.Title,
		IntroductionText: action.Context.Description,
		Elements:         action.Context.Fields,
		State:            strings.Join(state, "&"),
	}

	dialog := model.OpenDialogRequest{
		TriggerId: action.TriggerId,
		URL:       action.Context.Url,
		Dialog:    dialogData,
	}

	_, err := s.client.OpenInteractiveDialog(dialog)
	if err != nil {
		return fmt.Errorf("failed to open dialog. error: %w", err)
	}
	return nil
}
