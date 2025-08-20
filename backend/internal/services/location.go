package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/services/most"
	"github.com/Alexander272/mersi/backend/pkg/logger"
)

type LocationService struct {
	repo         repository.Location
	responsible  Responsible
	notification Notification
	most         *most.MostService
}

type LocationDeps struct {
	Repo         repository.Location
	Responsible  Responsible
	Notification Notification
	Most         *most.MostService
}

func NewLocationService(deps *LocationDeps) *LocationService {
	return &LocationService{
		repo:         deps.Repo,
		responsible:  deps.Responsible,
		notification: deps.Notification,
		most:         deps.Most,
	}
}

type Location interface {
	Get(ctx context.Context, dto *models.GetLocationDTO) ([]*models.Location, error)
	GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	GetSeveralLast(ctx context.Context, req *models.GetSeveralLocationsDTO) ([]*models.Location, error)
	GetUsedByHolder(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error)
	GetUsedByDepartment(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error)
	SelectByDepartments(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error)
	Create(ctx context.Context, dto *models.LocationDTO) error
	CreateSeveral(ctx context.Context, dto []*models.LocationDTO) (bool, error)
	Update(ctx context.Context, dto *models.LocationDTO) error
	SetPerson(ctx context.Context, personId string) error
	SetDepartment(ctx context.Context, departmentId string) error
	ReceivingFromApp(ctx context.Context, dto *models.ReceivingDTO) error
	ReceivingFromChannel(ctx context.Context, dto *models.DialogResponse) error
	ReceivingDialogOpen(ctx context.Context, dto *models.PostAction) error
	ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error
	ForcedReceiptAll(ctx context.Context) error
	Delete(ctx context.Context, dto *models.DeleteLocationDTO) error
}

func (s *LocationService) Get(ctx context.Context, dto *models.GetLocationDTO) ([]*models.Location, error) {
	data, err := s.repo.Get(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	data, err := s.repo.GetLast(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get last location. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) GetSeveralLast(ctx context.Context, req *models.GetSeveralLocationsDTO) ([]*models.Location, error) {
	data, err := s.repo.GetSeveralLast(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get several last locations. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) GetUsedByHolder(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error) {
	data, err := s.repo.GetUsedByHolder(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get used by holder. error: %w", err)
	}
	return data, nil
}
func (s *LocationService) GetUsedByDepartment(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error) {
	data, err := s.repo.GetUsedByDepartment(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get used by department. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) SelectByDepartments(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
	data, err := s.repo.SelectByDepartment(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to select by departments. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) Create(ctx context.Context, dto *models.LocationDTO) error {
	if dto.Status == constants.LocationStatusMoved && dto.DepartmentId != "" {
		// department, err := s.department.GetById(ctx, &models.GetDepartmentByIdDTO{Id: dto.DepartmentId})
		// if err != nil {
		// 	return err
		// }

		responsible, err := s.responsible.GetWithChannel(ctx, &models.GetResponsibleDTO{DepartmentId: dto.DepartmentId})
		if err != nil {
			return err
		}
		if len(responsible) == 0 { //? ответственное лицо не задано
			return models.ErrNoResponsible
		}
		if responsible[0].ChannelId == "" { //? канал для уведомлений не задан
			return models.ErrNoChannel
		}
	}

	if err := s.repo.Create(ctx, dto); err != nil {
		return fmt.Errorf("failed to create location. error: %w", err)
	}
	return nil
}

func (s *LocationService) CreateSeveral(ctx context.Context, dto []*models.LocationDTO) (bool, error) {
	if len(dto) == 0 {
		return false, models.ErrNoInstrument
	}

	responsible, err := s.responsible.GetBySSOId(ctx, dto[0].UserId)
	if err != nil {
		return false, err
	}

	isFull := true
	if dto[0].PersonId == "" && dto[0].NeedConfirm {
		if len(responsible) == 0 {
			return false, models.ErrNoResponsible
		}

		instrumentIds := []string{}
		locations := make(map[string]*models.LocationDTO)
		for _, l := range dto {
			instrumentIds = append(instrumentIds, l.InstrumentId)
			locations[l.InstrumentId] = l
		}

		ids := []string{}
		for _, r := range responsible {
			ids = append(ids, r.DepartmentId)
		}

		// при возвращении инструментов я отбрасываю те что не находятся в том же подразделении что и пользователь
		filtered, err := s.SelectByDepartments(ctx, &models.SelectByDepsDTO{InstrumentIds: instrumentIds, DepartmentIds: ids})
		if err != nil {
			return false, err
		}

		isFull = len(filtered) == len(dto)
		if !isFull {
			newLocations := []*models.LocationDTO{}
			for _, l := range filtered {
				newLocations = append(newLocations, locations[l])
			}
			dto = newLocations
		}
	}

	if len(dto) == 0 {
		return isFull, models.ErrNoInstrument
	}
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return false, fmt.Errorf("failed to create several locations. error: %w", err)
	}
	return isFull, nil
}

func (s *LocationService) Update(ctx context.Context, dto *models.LocationDTO) error {
	if err := s.repo.Update(ctx, dto); err != nil {
		return fmt.Errorf("failed to update location. error: %w", err)
	}
	return nil
}

func (s *LocationService) SetPerson(ctx context.Context, personId string) error {
	if err := s.repo.SetPerson(ctx, personId); err != nil {
		return fmt.Errorf("failed to set person. error: %w", err)
	}
	return nil
}
func (s *LocationService) SetDepartment(ctx context.Context, departmentId string) error {
	if err := s.repo.SetDepartment(ctx, departmentId); err != nil {
		return fmt.Errorf("failed to set department. error: %w", err)
	}
	return nil
}

// Получение инструментов. Запрос прилетает из веб-приложения
func (s *LocationService) ReceivingFromApp(ctx context.Context, dto *models.ReceivingDTO) error {
	responsible, err := s.responsible.GetBySSOId(ctx, dto.UserId)
	if err != nil {
		return err
	}

	logger.Debug("receiving", logger.StringAttr("status", dto.Status), logger.AnyAttr("responsible", responsible))
	if dto.Status == constants.LocationStatusUsed {
		if len(responsible) == 0 {
			return models.ErrNoResponsible
		}

		ids := []string{}
		for _, r := range responsible {
			ids = append(ids, r.DepartmentId)
		}

		// при подтверждении инструментов я отфильтровываю те что не находятся в том же подразделении что и пользователь
		filtered, err := s.SelectByDepartments(ctx, &models.SelectByDepsDTO{
			InstrumentIds: dto.InstrumentIds,
			DepartmentIds: ids,
			Status:        constants.LocationStatusMoved,
		})
		if err != nil {
			return err
		}

		isFull := len(filtered) == len(dto.InstrumentIds)
		if !isFull {
			dto.InstrumentIds = filtered
		}
	}

	if len(dto.InstrumentIds) == 0 {
		return models.ErrNoInstrument
	}

	if err := s.repo.Receiving(ctx, dto); err != nil {
		return fmt.Errorf("failed to receiving si location. error: %w", err)
	}
	return nil
}

// Получение инструментов. Запрос прилетает из канала mattermost
func (s *LocationService) ReceivingFromChannel(ctx context.Context, dto *models.DialogResponse) error {
	status := constants.LocationStatusUsed
	instrumentIds := []string{}

	state := strings.Split(dto.State, "&")
	for _, s := range state {
		arr := strings.SplitN(s, ":", 2)
		if arr[0] == "Status" {
			status = arr[1]
		}
	}

	for k, v := range dto.Submission {
		if v {
			instrumentIds = append(instrumentIds, k)
		}
	}

	location := &models.ReceivingDTO{
		InstrumentIds: instrumentIds,
		Status:        status,
		HasConfirmed:  true,
	}
	if err := s.repo.Receiving(ctx, location); err != nil {
		return fmt.Errorf("failed to receiving si location. error: %w", err)
	}

	if err := s.notification.CheckReceiving(ctx, dto); err != nil {
		return err
	}

	return nil
}
func (s *LocationService) ReceivingDialogOpen(ctx context.Context, dto *models.PostAction) error {
	if err := s.most.Open(ctx, dto); err != nil {
		return err
	}
	return nil
}

func (s *LocationService) ForcedReceipt(ctx context.Context, dto *models.ForcedReceiptDTO) error {
	if err := s.repo.ForcedReceipt(ctx, dto); err != nil {
		return fmt.Errorf("failed to forced receipt si. error: %w", err)
	}
	return nil
}

func (s *LocationService) ForcedReceiptAll(ctx context.Context) error {
	logger.Info("Forced receipt SI")
	if err := s.repo.ForcedReceiptAll(ctx); err != nil {
		return fmt.Errorf("failed to forced receipt all si. error: %w", err)
	}
	return nil
}

func (s *LocationService) Delete(ctx context.Context, dto *models.DeleteLocationDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete si location. error: %w", err)
	}
	return nil
}
