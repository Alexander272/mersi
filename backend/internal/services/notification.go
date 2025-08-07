package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/services/most"
	"github.com/Alexander272/mersi/backend/pkg/logger"
)

type NotificationService struct {
	si   SI
	file File
	most *most.MostService
}

type NotificationDeps struct {
	SI   SI
	File File
	Most *most.MostService
}

func NewNotificationService(deps *NotificationDeps) *NotificationService {
	return &NotificationService{
		si:   deps.SI,
		file: deps.File,
		most: deps.Most,
	}
}

type Notification interface {
	CheckVerification() error
}

func (s *NotificationService) CheckVerification() error {
	logger.Info("Check verification")

	now := time.Now()
	if time.Now().Day() != 5 {
		return nil
	}

	dto := &models.Period{
		StartAt:  time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix(),
		FinishAt: time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Unix(),
	}

	data, err := s.si.GetVerification(context.Background(), dto)
	if err != nil {
		return err
	}

	for _, d := range data {
		switch d.BidType {
		case "ointo_si":
			if err := s.sendFile(d); err != nil {
				return err
			}

		default:
			if err := s.sendVerificationTable(d); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *NotificationService) sendFile(dto *models.SiVerification) error {
	doc, err := s.file.MakeDocSchedule(context.Background(), dto.SI)
	if err != nil {
		return err
	}
	if doc == nil {
		return fmt.Errorf("doc is nil")
	}

	uploadDTO := &models.UploadFileDTO{
		Data:      doc,
		ChannelId: dto.NotificationChannel,
		Filename:  "Поверка.docx",
	}
	fileId, err := s.most.Post.Upload(context.Background(), uploadDTO)
	if err != nil {
		return err
	}
	logger.Debug("Check verification", logger.StringAttr("fileId", fileId))

	post := &models.CreatePostDTO{
		ChannelId: dto.NotificationChannel,
		Message:   "#### В следующем месяце подходит срок поверки у следующих инструментов",
		FileIds:   []string{fileId},
	}
	if err := s.most.Post.Create(context.Background(), post); err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) sendVerificationTable(dto *models.SiVerification) error {
	tableHeader := "| Наименование СИ | Вид (тип, марка) СИ | зав.№ | Дата следующей поверки |"
	title := "В следующем месяце подходит срок поверки у следующих инструментов"

	if dto.BidType == "ointo_eq" {
		tableHeader = "| Наименование ИО | Марка, тип| зав.№ | Следующая аттестация |"
		title = "В следующем месяце подходит срок поверки у следующего оборудования"
	}

	table := []string{
		tableHeader,
		"|:--|:--|:--|:--|",
	}

	for _, d := range dto.SI {
		table = append(table, fmt.Sprintf("| %s| %s | %s | %s |",
			d.Name, d.Type, d.FactoryNumber, time.Unix(d.NextVerificationDate, 0).Format("02.01.2006")),
		)
	}

	post := &models.CreatePostDTO{
		ChannelId: dto.NotificationChannel,
		Message:   fmt.Sprintf("#### %s\n%s", title, strings.Join(table, "\n")),
	}
	if err := s.most.Post.Create(context.Background(), post); err != nil {
		return err
	}

	return nil
}
