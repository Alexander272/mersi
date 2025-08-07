package most

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/mattermost/mattermost-server/v6/model"
)

type PostService struct {
	client  *model.Client4
	channel Channel
}

func NewPostService(client *model.Client4, channel Channel) *PostService {
	return &PostService{
		client:  client,
		channel: channel,
	}
}

type Post interface {
	Create(ctx context.Context, dto *models.CreatePostDTO) error
	Update(ctx context.Context, dto *models.UpdatePostDTO) error
	Upload(ctx context.Context, dto *models.UploadFileDTO) (string, error)
}

func (s *PostService) Create(ctx context.Context, dto *models.CreatePostDTO) error {
	post := &model.Post{
		ChannelId: dto.ChannelId,
		Message:   dto.Message,
		IsPinned:  dto.IsPinned,
		FileIds:   dto.FileIds,
	}

	if dto.UserId != "" && post.ChannelId == "" {
		channelId, err := s.channel.Create(ctx, dto.UserId)
		if err != nil {
			return err
		}
		post.ChannelId = channelId
	}

	if post.ChannelId == "" {
		return fmt.Errorf("ChannelId is empty")
	}

	if dto.Actions != nil {
		attachment := &model.SlackAttachment{
			Actions: dto.Actions,
		}
		post.AddProp("attachments", []*model.SlackAttachment{attachment})
	}
	if dto.Attachments != nil {
		post.AddProp("attachments", dto.Attachments)
	}

	for _, p := range dto.Props {
		post.AddProp(p.Key, p.Value)
	}

	if dto.IsPinned {
		dataId := post.GetProp("data_id")
		dataType := post.GetProp("data_type")

		if dataId != nil {
			req := &models.GetPost{
				ChannelId: post.ChannelId,
				DataType:  dataType.(string),
				DataId:    dataId.(string),
			}
			if err := s.findDuplicate(req); err != nil {
				return fmt.Errorf("failed to find duplicate. error: %w", err)
			}
		}
	}

	//TODO
	// можно передавать ID поста. выполнять поиск в канале этого ID, если он есть удалять сообщение и отправлять новое
	// благодаря такой схеме можно избежать варианта, когда в канале три одинаковых сообщений (для получения инструмента) с разной датой, а изменяться по нажатию будет только последнее
	// если я задаю ID метод выдает ошибку (только какого хрена он позволяет его задавать)
	// если записывать ID поста в props и закреплять сообщение, а потом получать все закрепленные и искать в props нужный мне ID
	// в ID поста надо учитывать ID пользователя из-за которого этот пост создавался или что-то подобное (у метролога может быть много всего)

	_, _, err := s.client.CreatePost(post)
	if err != nil {
		return fmt.Errorf("failed to create post. error: %w", err)
	}
	return nil
}

func (s *PostService) Update(ctx context.Context, dto *models.UpdatePostDTO) error {
	// тут тогда надо будет откреплять пост (скорее всего, можно будет передавать флаг)
	post := &model.Post{
		Id:       dto.PostId,
		Message:  dto.Message,
		IsPinned: false,
	}

	if dto.Actions != nil {
		attachment := &model.SlackAttachment{
			Actions: dto.Actions,
		}
		post.AddProp("attachments", []*model.SlackAttachment{attachment})
	}
	if dto.Attachments != nil {
		post.AddProp("attachments", dto.Attachments)
	}

	for _, p := range dto.Props {
		post.AddProp(p.Key, p.Value)
	}

	_, _, err := s.client.UpdatePost(dto.PostId, post)
	if err != nil {
		return fmt.Errorf("failed to update post. error: %w", err)
	}
	return nil
}

func (s *PostService) Upload(ctx context.Context, dto *models.UploadFileDTO) (string, error) {
	res, _, err := s.client.UploadFile(dto.Data.Bytes(), dto.ChannelId, dto.Filename)
	if err != nil {
		return "", fmt.Errorf("failed to upload file. error: %w", err)
	}

	if len(res.FileInfos) == 0 {
		return "", fmt.Errorf("failed to upload file. len(res.FileInfos) == 0")
	}
	return res.FileInfos[0].Id, nil
}

func (s *PostService) findDuplicate(req *models.GetPost) error {
	posts, _, err := s.client.GetPinnedPosts(req.ChannelId, "")
	if err != nil {
		return fmt.Errorf("failed to get pinned posts. error: %w", err)
	}

	dataIds := models.Universe{}
	if req.DataType == "array" {
		dataIds = models.NewUniverse(strings.Split(req.DataId, ","))
	}

	for _, p := range posts.Posts {
		if req.DataType == "array" {
			if p.GetProp("data_id") == nil {
				continue
			}

			tmp := p.GetProp("data_id").(string)
			data := models.NewUniverse(strings.Split(tmp, ","))
			ok := false
			if len(data) > len(dataIds) {
				ok = data.ContainSet(strings.Split(req.DataId, ","))
			} else {
				ok = dataIds.ContainSet(strings.Split(tmp, ","))
			}

			if !ok {
				continue
			}
			_, err := s.client.DeletePost(p.Id)
			if err != nil {
				return fmt.Errorf("failed to delete post. error: %w", err)
			}
			break
		} else {
			if p.GetProp("data_id") == req.DataId {
				_, err := s.client.DeletePost(p.Id)
				if err != nil {
					return fmt.Errorf("failed to delete post. error: %w", err)
				}
				break
			}
		}
	}

	return nil
}
