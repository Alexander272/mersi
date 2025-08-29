package services

import (
	"context"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/go-co-op/gocron/v2"
)

type SchedulerService struct {
	cron         gocron.Scheduler
	notification Notification
	user         User
	location     Location
	documents    Document
}

type SchedulerDeps struct {
	Notification Notification
	User         User
	Location     Location
	Documents    Document
}

func NewSchedulerService(deps *SchedulerDeps) *SchedulerService {
	cron, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create new scheduler. error: %s", err.Error())
	}

	return &SchedulerService{
		cron:         cron,
		notification: deps.Notification,
		user:         deps.User,
		location:     deps.Location,
		documents:    deps.Documents,
	}
}

type Scheduler interface {
	Start(conf *config.SchedulerConfig) error
	Stop() error
}

// запуск заданий в cron
func (s *SchedulerService) Start(conf *config.SchedulerConfig) error {
	now := time.Now()

	hours := int(conf.StartTime.Hours())
	minutes := int(math.Round(math.Mod(conf.StartTime.Hours(), 1) * 60))

	jobStart := time.Date(now.Year(), now.Month(), now.Day(), hours, minutes, 0, 0, now.Location())
	if now.Hour() > hours || (now.Hour() == hours && now.Minute() >= minutes) {
		jobStart = jobStart.Add(24 * time.Hour)
	}
	// // вернуть нормальное время запуска
	// jobStart := now.Add(1 * time.Minute)
	logger.Info("starting jobs time " + jobStart.Format("02.01.2006 15:04:05"))

	job := gocron.DurationJob(conf.Interval)
	task := gocron.NewTask(s.job)
	jobStartAt := gocron.WithStartAt(gocron.WithStartDateTime(jobStart))

	_, err := s.cron.NewJob(job, task, jobStartAt)
	if err != nil {
		return fmt.Errorf("failed to create new job. error: %w", err)
	}

	//? запуск крона через интервал (по умолчанию день)
	s.cron.Start()
	return nil
}

// остановка заданий в cron
func (s *SchedulerService) Stop() error {
	if err := s.cron.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown cron scheduler. error: %w", err)
	}
	return nil
}

func (s *SchedulerService) job() {
	logger.Info("job was started")

	// принудительное получение инструмента, который не принимают больше 20 дней
	if err := s.location.ForcedReceiptAll(context.Background()); err != nil {
		logger.Error("location forced receipt error:", logger.ErrAttr(err))
		error_bot.Send(nil, err.Error(), nil)
		return
	}

	// проверка необходимости сдать использующиеся инструменты на поверку
	if err := s.notification.CheckUsed(); err != nil {
		logger.Error("notification check used error:", logger.ErrAttr(err))
		error_bot.Send(nil, err.Error(), nil)
	}

	// проверка отправленных инструментов и рассылка уведомлений для подтверждения получения
	if err := s.notification.CheckSent(); err != nil {
		logger.Error("notification check sent error:", logger.ErrAttr(err))
		error_bot.Send(nil, err.Error(), nil)
	}

	// проверка инструментов на необходимость поверки
	if err := s.notification.CheckVerification(); err != nil {
		logger.Error("notification check verification error:", logger.ErrAttr(err))
		error_bot.Send(nil, err.Error(), nil)
	}

	// Синхронизация пользователей с keycloak
	if err := s.user.Sync(context.Background()); err != nil {
		logger.Error("user sync error:", logger.ErrAttr(err))
		error_bot.Send(nil, err.Error(), nil)
		return
	}

	// Удаление пустых папок
	if err := s.documents.RemoveEmptyFolders(context.Background()); err != nil {
		logger.Error("delete empty folders error:", logger.ErrAttr(err))
		error_bot.Send(nil, err.Error(), nil)
		return
	}
}
