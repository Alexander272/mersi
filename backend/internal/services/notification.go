package services

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/services/most"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/goccy/go-json"
	"github.com/goodsign/monday"
	"github.com/mattermost/mattermost-server/v6/model"
)

type NotificationService struct {
	si        SI
	file      File
	section   Section
	most      most.Post
	conf      config.UsedConfig
	iteration int
	now       func() time.Time
}

type NotificationDeps struct {
	SI      SI
	File    File
	Section Section
	Most    most.Post
	Conf    config.UsedConfig
}

func NewNotificationService(deps *NotificationDeps) *NotificationService {
	return &NotificationService{
		si:        deps.SI,
		file:      deps.File,
		section:   deps.Section,
		most:      deps.Most,
		conf:      deps.Conf,
		iteration: 0,
		now:       time.Now,
	}
}

type Notification interface {
	CheckSent(ctx context.Context) error
	CheckUsed(ctx context.Context) error
	CheckVerification(ctx context.Context) error
	CheckReceiving(ctx context.Context, dto *models.DialogResponse) error
}

func (s *NotificationService) CheckSent(ctx context.Context) error {
	logger.Info("Check sent")

	data, err := s.si.GetSent(ctx, &models.GetSiDTO{})
	if err != nil {
		return err
	}

	for _, d := range data {
		if d.Channel == "" || len(d.SI) == 0 {
			continue
		}

		if err := s.sendInstruments(ctx, d); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) CheckUsed(ctx context.Context) error {
	logger.Info("Check used")
	index := s.iteration % len(s.conf.Times)

	now := s.now()
	monthEnd := time.Date(now.Year(), now.Month()+1, 0, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	if s.iteration >= len(s.conf.Times) {
		monthEnd = time.Date(now.Year(), now.Month()+2, 0, now.Hour(), now.Minute(), now.Second(), now.Nanosecond(), now.Location())
	}

	if monthEnd.Add(-s.conf.Times[len(s.conf.Times)-1]).Before(now) {
		s.iteration = 0
	}

	date := monthEnd.Add(-s.conf.Times[s.iteration])
	if date.After(now) {
		return nil
	}

	startAt := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	finishAt := time.Date(now.Year(), now.Month()+2, 0, 0, 0, 0, 0, now.Location())

	period := &models.Period{
		StartAt:  startAt,
		FinishAt: finishAt,
	}

	params := &models.GetAllSectionsDTO{
		IsActive:        models.NewOptional(true),
		HasReturnNotice: models.NewOptional(true),
	}
	sections, err := s.section.GetAll(ctx, params)
	if err != nil {
		return err
	}

	for _, section := range sections {
		period.SectionId = section.ID

		data, err := s.si.GetUsed(ctx, period)
		if err != nil {
			return err
		}

		for _, d := range data {
			table := []string{
				"| Наименование СИ | зав.№ | Держатель |",
				"|:--|:--|:--|",
			}

			for _, si := range d.SI {
				table = append(table, fmt.Sprintf("|%s|%s|%s|", si.Name, si.FactoryNumber, si.Person))
			}
			term := monday.Format(monthEnd.Add(-s.conf.Times[len(s.conf.Times)-1]), "Mon 2 Jan 2006", monday.LocaleRuRU)
			if s.iteration == len(s.conf.Times)-1 {
				term += " (Сегодня)"
			}

			post := &models.CreatePostDTO{
				ChannelId: d.Channel,
			}
			post.Message = "#### Необходимо сдать инструменты до `" + term + "`\n" + strings.Join(table, "\n")
			post.Props = []*models.Props{
				{Key: "service", Value: "sia"},
			}

			if err := s.most.Create(ctx, post); err != nil {
				return err
			}
		}
	}

	if s.iteration >= len(s.conf.Times) {
		s.iteration = index
	}
	s.iteration = (index + 1)

	return nil
}

func (s *NotificationService) CheckVerification(ctx context.Context) error {
	logger.Info("Check verification")

	params := &models.GetAllSectionsDTO{IsActive: models.NewOptional(true)}
	sections, err := s.section.GetAll(ctx, params)
	if err != nil {
		return err
	}

	// может привязать день к realm или section, только тогда надо как-то получать данные
	// хотя можно получать все section и делать запросы в цикле, если полученный день совпадает с текущим
	// можно еще разрешить отрицательные значения, чтобы можно было считать от конца месяца

	now := s.now()
	for _, item := range sections {
		if item.VerificationDay < 0 {
			if time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location()).Day()-item.VerificationDay != now.Day() {
				return nil
			}
		} else {
			if item.VerificationDay != now.Day() {
				return nil
			}
		}
		logger.Debug("Check verification", logger.AnyAttr("section", item))

		dto := &models.Period{
			StartAt:         time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location()),
			FinishAt:        time.Date(now.Year(), now.Month()+2, 0, 0, 0, 0, 0, now.Location()),
			SectionId:       item.ID,
			ChannelIsOption: false,
		}

		data, err := s.si.GetVerification(ctx, dto)
		if err != nil {
			return err
		}

		for _, d := range data {
			switch d.BidType {
			case "ointo_si":
				if err := s.sendFile(ctx, d); err != nil {
					return err
				}

			default:
				if err := s.sendVerificationTable(ctx, d); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *NotificationService) sendFile(ctx context.Context, dto *models.SiVerification) error {
	doc, err := s.file.MakeDocSchedule(ctx, dto.SI)
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
	fileId, err := s.most.Upload(ctx, uploadDTO)
	if err != nil {
		return err
	}
	logger.Debug("Check verification", logger.StringAttr("fileId", fileId))

	post := &models.CreatePostDTO{
		ChannelId: dto.NotificationChannel,
		Message:   "#### В следующем месяце подходит срок поверки у следующих инструментов",
		FileIds:   []string{fileId},
	}
	if err := s.most.Create(ctx, post); err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) sendVerificationTable(ctx context.Context, dto *models.SiVerification) error {
	tableHeader := "| Наименование СИ | Вид (тип, марка) СИ | зав.№ | Дата следующей поверки |"
	title := "В следующем месяце подходит срок поверки у следующих инструментов"

	if dto.BidType == "ointo_eq" {
		tableHeader = "| Наименование ИО | Марка, тип| зав.№ | Следующая аттестация |"
		title = "В следующем месяце подходит срок аттестации у следующего оборудования"
	}

	table := []string{
		tableHeader,
		"|:--|:--|:--|:--|",
	}

	for _, d := range dto.SI {
		table = append(table, fmt.Sprintf("| %s| %s | %s | %s |",
			d.Name, d.Type, d.FactoryNumber, d.NextVerificationDate.Format("02.01.2006")),
		)
	}

	post := &models.CreatePostDTO{
		ChannelId: dto.NotificationChannel,
		Message:   fmt.Sprintf("#### %s\n%s", title, strings.Join(table, "\n")),
	}
	if err := s.most.Create(ctx, post); err != nil {
		return err
	}

	return nil
}

func (s *NotificationService) CheckReceiving(ctx context.Context, dto *models.DialogResponse) error {
	instrumentsDTO := &models.SiReceiving{
		PostId:  "",
		Status:  "",
		Channel: dto.ChannelID,
	}
	instruments := []*models.SI{}
	accept := []*models.SI{}
	missing := []*models.SI{}

	state := strings.Split(dto.State, "&")
	for _, s := range state {
		arr := strings.SplitN(s, ":", 2)
		switch arr[0] {
		case "PostId":
			instrumentsDTO.PostId = arr[1]
		case "Status":
			instrumentsDTO.Status = arr[1]
		case "SI":
			err := json.Unmarshal([]byte(arr[1]), &instruments)
			if err != nil {
				return fmt.Errorf("failed to json unmarshal. error: %w", err)
			}
		}
	}

	for _, v := range instruments {
		if value, exist := dto.Submission[v.Id]; !exist || !value {
			missing = append(missing, v)
		} else {
			accept = append(accept, v)
		}
	}

	instrumentsDTO.SI = accept
	instrumentsDTO.Place = accept[0].Place
	if err := s.updateInstruments(ctx, instrumentsDTO); err != nil {
		return err
	}

	if len(missing) > 0 {
		instrumentsDTO.SI = missing
		if err := s.sendInstruments(ctx, instrumentsDTO); err != nil {
			return err
		}
	}
	return nil
}

func (s *NotificationService) sendInstruments(ctx context.Context, dto *models.SiReceiving) error {
	post := &models.CreatePostDTO{
		ChannelId: dto.Channel,
		IsPinned:  true,
		// IsPinned:  dto.Status == constants.StatusReceiving,
	}

	columns := []string{"Наименование СИ", "зав.№"}
	if dto.SI[0].Place != "" {
		columns = append(columns, "Держатель")
	}
	table := []string{
		fmt.Sprintf("| %s |", strings.Join(columns, "|")),
		fmt.Sprintf("| %s", strings.Repeat(":--|", len(columns))),
	}

	instrumentIds := []string{}
	fields := []models.FormField{}

	for _, si := range dto.SI {
		instrumentIds = append(instrumentIds, si.Id)
		row := []string{si.Name, si.FactoryNumber}
		if dto.SI[0].Place != "" {
			row = append(row, si.Person)
		}
		table = append(table, fmt.Sprintf("|%s|", strings.Join(row, "|")))

		fields = append(fields, models.FormField{
			Id:       si.Id,
			Title:    fmt.Sprintf("%s (%s)", si.Name, si.FactoryNumber),
			Name:     "Инструмент получен",
			Type:     "bool",
			Default:  "true",
			Optional: true,
		})
	}

	place := ""
	if dto.SI[0].Place != "" {
		place = fmt.Sprintf("(%s)", dto.SI[0].Place)
	}
	post.Message = fmt.Sprintf("#### Подтвердите получение инструментов %s\n%s", place, strings.Join(table, "\n"))

	post.Props = []*models.Props{
		{Key: "service", Value: "sia"},
		{Key: "data_type", Value: "array"},
		{Key: "data_id", Value: strings.Join(instrumentIds, ",")},
	}

	// if dto.Status == constants.StatusReceiving {
	host := os.Getenv("HOST_URL")
	url := host + "/api/v1/si/locations/receiving/dialogs"

	j, err := json.Marshal(dto.SI)
	if err != nil {
		return fmt.Errorf("failed to marshal json. error: %w", err)
	}

	action := &model.PostAction{
		Id:    constants.StatusReceiving,
		Name:  "Получить",
		Style: "primary",
		Integration: &model.PostActionIntegration{
			URL: host + "/api/v1/si/locations/receiving/dialogs/open",
			Context: map[string]interface{}{
				"url":         url,
				"title":       "Получение инструментов",
				"description": "#### Отметьте полученные инструменты",
				"callbackId":  "receiving_form",
				"state":       fmt.Sprintf("Status:%s&SI:%s", dto.Status, string(j)),
				"fields":      fields,
			},
		},
	}

	post.Actions = []*model.PostAction{action}
	// }

	if err := s.most.Create(ctx, post); err != nil {
		return err
	}
	return nil
}

func (s *NotificationService) updateInstruments(ctx context.Context, dto *models.SiReceiving) error {
	post := &models.UpdatePostDTO{
		PostId:   dto.PostId,
		IsPinned: false,
	}

	// lines := []string{
	// 	"| Наименование СИ | зав.№ | Держатель |",
	// 	"|:--|:--|:--|",
	// }
	columns := []string{"Наименование СИ", "зав.№"}
	if dto.SI[0].Place != "" {
		columns = append(columns, "Держатель")
	}
	lines := []string{
		fmt.Sprintf("| %s |", strings.Join(columns, "|")),
		fmt.Sprintf("| %s", strings.Repeat(":--|", len(columns))),
	}

	for _, si := range dto.SI {
		// lines = append(lines, fmt.Sprintf("|%s|%s|%s|", si.Name, si.FactoryNumber, si.Person))
		row := []string{si.Name, si.FactoryNumber}
		if dto.SI[0].Place != "" {
			row = append(row, si.Person)
		}
		lines = append(lines, fmt.Sprintf("|%s|", strings.Join(row, "|")))
	}
	post.Message = "#### Получены инструменты \n" + strings.Join(lines, "\n")
	post.Props = []*models.Props{
		{Key: "service", Value: "sia"},
	}

	if err := s.most.Update(ctx, post); err != nil {
		return err
	}
	return nil
}
