package most

import "github.com/Alexander272/mersi/backend/pkg/mattermost"

type MostService struct {
	Channel
	Dialog
	Post
}

type MostDeps struct {
	Client *mattermost.Client
}

// type Most interface {
// 	Channel
// 	Dialog
// 	Post
// }

func NewMostService(deps MostDeps) *MostService {
	channel := NewChannelService(deps.Client.Http, deps.Client.Name)
	dialog := NewDialogService(deps.Client.Http)
	post := NewPostService(deps.Client.Http, channel)

	return &MostService{
		Channel: channel,
		Dialog:  dialog,
		Post:    post,
	}
}
