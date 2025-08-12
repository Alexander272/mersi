package models

import (
	"bytes"

	"github.com/mattermost/mattermost-server/v6/model"
)

type CreatePostDTO struct {
	UserId      string                   `json:"userId"`
	ChannelId   string                   `json:"channelId"`
	Message     string                   `json:"message" binding:"required"`
	IsPinned    bool                     `json:"isPinned"`
	FileIds     []string                 `json:"fileIds"`
	Props       []*Props                 `json:"props"`
	Actions     []*model.PostAction      `json:"actions"`
	Attachments []*model.SlackAttachment `json:"attachments"`
}

type UpdatePostDTO struct {
	PostId      string                   `json:"postId" binding:"required"`
	Message     string                   `json:"message" binding:"required"`
	IsPinned    bool                     `json:"isPinned"`
	FileIds     []string                 `json:"fileIds"`
	Props       []*Props                 `json:"props"`
	Actions     []*model.PostAction      `json:"actions"`
	Attachments []*model.SlackAttachment `json:"attachments"`
}

type GetPost struct {
	ChannelId string `json:"channelId"`
	DataType  string `json:"dataType"`
	DataId    string `json:"dataId"`
}

type UploadFileDTO struct {
	ChannelId string
	Filename  string
	Data      *bytes.Buffer
}

type Props struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DialogResponse struct {
	Type       string          `json:"type"`
	CallbackID string          `json:"callback_id"`
	State      string          `json:"state"`
	UserID     string          `json:"user_id"`
	ChannelID  string          `json:"channel_id"`
	TeamID     string          `json:"team_id"`
	Submission map[string]bool `json:"submission"`
	Cancelled  bool            `json:"cancelled"`
}

type PostAction struct {
	UserId      string            `json:"user_id"`
	UserName    string            `json:"user_name"`
	ChannelId   string            `json:"channel_id"`
	ChannelName string            `json:"channel_name"`
	TeamId      string            `json:"team_id"`
	TeamName    string            `json:"team_domain"`
	PostId      string            `json:"post_id"`
	TriggerId   string            `json:"trigger_id" binding:"required"`
	Type        string            `json:"type"`
	DataSource  string            `json:"data_source"`
	Context     PostActionContext `json:"context,omitempty" binding:"required"`
}

type PostActionContext struct {
	Url         string                `json:"url"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	CallbackId  string                `json:"callbackId"`
	SubmitLabel string                `json:"submitLabel"`
	State       string                `json:"state"`
	Fields      []model.DialogElement `json:"fields"`
}

type Universe map[string]bool

func NewUniverse(s []string) Universe {
	u := make(Universe)
	for _, i := range s {
		u[i] = true
	}
	return u
}

func (u Universe) ContainSet(s []string) bool {
	for _, i := range s {
		if !u[i] {
			return false
		}
	}
	return true
}
