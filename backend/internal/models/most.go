package models

import "github.com/mattermost/mattermost-server/v6/model"

type CreatePostDTO struct {
	UserId      string                   `json:"userId"`
	ChannelId   string                   `json:"channelId"`
	Message     string                   `json:"message" binding:"required"`
	IsPinned    bool                     `json:"isPinned"`
	Props       []*Props                 `json:"props"`
	Actions     []*model.PostAction      `json:"actions"`
	Attachments []*model.SlackAttachment `json:"attachments"`
}

type UpdatePostDTO struct {
	PostId      string                   `json:"postId" binding:"required"`
	Message     string                   `json:"message" binding:"required"`
	Props       []*Props                 `json:"props"`
	Actions     []*model.PostAction      `json:"actions"`
	Attachments []*model.SlackAttachment `json:"attachments"`
}

type Props struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
