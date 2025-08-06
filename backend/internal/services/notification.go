package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/pkg/logger"
)

type NotificationService struct {
	si   SI
	file File
}

type NotificationDeps struct {
	SI   SI
	File File
}

func NewNotificationService(deps *NotificationDeps) *NotificationService {
	return &NotificationService{
		si:   deps.SI,
		file: deps.File,
	}
}

type Notification interface{}

func (s *NotificationService) CheckVerification() error {
	logger.Info("Check verification")

	now := time.Now()
	if time.Now().Day() != 5 {
		return nil
	}

	dto := &models.Period{
		StartAt:  time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()).Unix(),
		FinishAt: time.Date(now.Year(), now.Month()+2, 0, 0, 0, 0, 0, now.Location()).Unix(),
	}

	data, err := s.si.GetVerification(context.Background(), dto)
	if err != nil {
		return err
	}

	//TODO надо определять нужен ли мне docx

	for _, d := range data {
		switch d.BidType {
		case "ointo_si":
			// doc, err := s.file.MakeDocSchedule(context.Background(), d.SI)
			// if err != nil {
			// 	return err
			// }
			// post := &models.CreatePostDTO{
			// 	ChannelId: d.NotificationChannel,
			// 	Message: "В следующем месяце подходит срок поверки у следующих инструментов",
			// 	Attachments: []*most_model.SlackAttachment{

			// 	},
			// }

		default:
		}
	}
	// doc, err := s.file.MakeDocSchedule(context.Background(), si)

	return fmt.Errorf("not implemented")
}
