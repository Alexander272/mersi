package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type LocationService struct {
	repo        repository.Location
	txManager   TransactionManager
	responsible Responsible
	activityLog ActivityLog
}

type LocationDeps struct {
	Repo        repository.Location
	TxManager   TransactionManager
	Responsible Responsible
	ActivityLog ActivityLog
}

func NewLocationService(deps *LocationDeps) *LocationService {
	return &LocationService{
		repo:        deps.Repo,
		txManager:   deps.TxManager,
		responsible: deps.Responsible,
		activityLog: deps.ActivityLog,
	}
}

type Location interface {
	Get(ctx context.Context, dto *models.GetLocationDTO) ([]*models.Location, error)
	GetById(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error)
	GetSeveralLast(ctx context.Context, req *models.GetSeveralLocationsDTO) ([]*models.Location, error)
	GetUsedByHolder(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error)
	GetUsedByDepartment(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error)
	SelectByDepartments(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error)
	Create(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error
	// Create(ctx context.Context, dto *models.LocationDTO) error
	CreateSeveral(ctx context.Context, dto []*models.LocationDTO) (bool, error)
	Update(ctx context.Context, dto *models.LocationDTO) error
	SetPerson(ctx context.Context, personId string) error
	SetDepartment(ctx context.Context, departmentId string) error
	Delete(ctx context.Context, dto *models.DeleteLocationDTO) error
}

func (s *LocationService) Get(ctx context.Context, dto *models.GetLocationDTO) ([]*models.Location, error) {
	data, err := s.repo.Get(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get locations. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) GetById(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	data, err := s.repo.GetById(ctx, dto)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get location by id. error: %w", err)
	}
	return data, nil
}

func (s *LocationService) GetLast(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
	data, err := s.repo.GetLast(ctx, dto)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
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

func (s *LocationService) Create(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error {
	if tx == nil {
		return s.txManager.ExecuteInTx(ctx, func(newTx postgres.Tx) error {
			return s.executeCreate(ctx, newTx, dto)
		})
	}
	return s.executeCreate(ctx, tx, dto)
}

func (s *LocationService) executeCreate(ctx context.Context, tx postgres.Tx, dto *models.LocationDTO) error {
	if dto.Status == constants.LocationStatusMoved && dto.DepartmentId != "" {
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

	if err := s.repo.CreateInTx(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create location. error: %w", err)
	}

	recordName := fmt.Sprintf("%s (%s)", dto.DepartmentId, dto.Status)
	// Логирование создания
	s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
		TableName:  "locations",
		RecordId:   dto.Id,
		RecordName: recordName,
		Action:     "CREATE",
		UserId:     dto.Actor.ID,
		UserName:   dto.Actor.Name,
		NewValue:   dto,
	})

	return nil
}

func (s *LocationService) CreateSeveral(ctx context.Context, dto []*models.LocationDTO) (bool, error) {
	if len(dto) == 0 {
		return false, models.ErrCannotMoveInstrument
	}

	responsible, err := s.responsible.GetBySSOId(ctx, dto[0].UserId)
	if err != nil {
		return false, err
	}

	isFull := true
	if dto[0].PersonId == "" && dto[0].NeedConfirm {
		if len(responsible) == 0 {
			return false, models.ErrNotResponsible
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
		depDTO := &models.SelectByDepsDTO{InstrumentIds: instrumentIds, DepartmentIds: ids}
		filtered, err := s.SelectByDepartments(ctx, depDTO)
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
		return isFull, models.ErrCannotMoveInstrument
	}
	if err := s.repo.CreateSeveral(ctx, dto); err != nil {
		return false, fmt.Errorf("failed to create several locations. error: %w", err)
	}
	return isFull, nil
}

func (s *LocationService) Update(ctx context.Context, dto *models.LocationDTO) error {
	// Получаем старые данные для логирования
	oldData, err := s.repo.GetById(ctx, &models.GetLocationDTO{Id: dto.Id})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return fmt.Errorf("failed to get old location data. error: %w", err)
	}

	if err := s.repo.Update(ctx, dto); err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("location %s not found", dto.Id)
		}
		return fmt.Errorf("failed to update location. error: %w", err)
	}

	recordName := fmt.Sprintf("%s (%s)", dto.DepartmentId, dto.Status)
	// Логирование обновления
	if oldData != nil {
		s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
			TableName:  "locations",
			RecordId:   dto.Id,
			RecordName: recordName,
			Action:     "UPDATE",
			UserId:     dto.Actor.ID,
			UserName:   dto.Actor.Name,
			OldValue:   oldData,
			NewValue:   dto,
		})
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

func (s *LocationService) Delete(ctx context.Context, dto *models.DeleteLocationDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		// Получаем старые данные для логирования
		oldData, err := s.repo.GetById(ctx, &models.GetLocationDTO{Id: dto.Id})
		if err != nil && !errors.Is(err, models.ErrNoRows) {
			return fmt.Errorf("failed to get old location data. error: %w", err)
		}

		if err := s.repo.Delete(ctx, tx, dto); err != nil {
			if errors.Is(err, models.ErrNoRows) {
				return models.ErrSingleLocationDelete
			}
			return fmt.Errorf("failed to delete location. error: %w", err)
		}

		// Логирование удаления
		if oldData != nil {
			recordName := fmt.Sprintf("%s (%s)", oldData.Place, oldData.Status)
			s.activityLog.LogActivity(ctx, &models.CreateActivityLogDTO{
				TableName:  "locations",
				RecordId:   dto.Id,
				RecordName: recordName,
				Action:     "DELETE",
				UserId:     dto.Actor.ID,
				UserName:   dto.Actor.Name,
				OldValue:   oldData,
			})
		}
		return nil
	})
}
