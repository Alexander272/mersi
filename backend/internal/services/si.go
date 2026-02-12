package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type SIService struct {
	repo         repository.SI
	txManager    TransactionManager
	instrument   Instrument
	verification Verification
	location     Location
}

type SiDeps struct {
	Repo         repository.SI
	TxManager    TransactionManager
	Instrument   Instrument
	Verification Verification
	Location     Location
}

func NewSiService(deps *SiDeps) *SIService {
	return &SIService{
		repo:         deps.Repo,
		txManager:    deps.TxManager,
		instrument:   deps.Instrument,
		verification: deps.Verification,
		location:     deps.Location,
	}
}

type SI interface {
	Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error)
	GetById(ctx context.Context, req *models.GetSiByIdDTO) (*models.BaseSI, error)
	GetVerification(ctx context.Context, req *models.Period) ([]*models.SiVerification, error)
	GetSent(ctx context.Context, req *models.GetSiDTO) ([]*models.SiReceiving, error)
	GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error)
	Create(ctx context.Context, dto *models.SiDTO) error
	Update(ctx context.Context, dto *models.SiDTO) error
	ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error
	Delete(ctx context.Context, dto *models.DeleteSiDTO) error
}

func (s *SIService) Get(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get si. error: %w", err)
	}
	return data, nil
}

func (s *SIService) GetById(ctx context.Context, req *models.GetSiByIdDTO) (*models.BaseSI, error) {
	instrument, err := s.instrument.GetById(ctx, &models.GetInstrumentByIdDTO{Id: req.Id})
	if err != nil {
		return nil, err
	}
	verification, err := s.verification.GetLast(ctx, &models.GetVerificationDTO{InstrumentId: req.Id})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return nil, err
	}
	location, err := s.location.GetLast(ctx, &models.GetLocationDTO{InstrumentId: req.Id})
	if err != nil && !errors.Is(err, models.ErrNoRows) {
		return nil, err
	}
	data := &models.BaseSI{
		Instrument:   instrument,
		Verification: verification,
		Location:     location,
	}
	return data, nil
}

func (s *SIService) GetVerification(ctx context.Context, req *models.Period) ([]*models.SiVerification, error) {
	data, err := s.repo.GetVerification(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get si verification. error: %w", err)
	}
	return data, nil
}

func (s *SIService) GetSent(ctx context.Context, req *models.GetSiDTO) ([]*models.SiReceiving, error) {
	data, err := s.repo.GetSent(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent si. error: %w", err)
	}
	return data, nil
}

func (s *SIService) GetUsed(ctx context.Context, req *models.Period) ([]*models.SiReceiving, error) {
	data, err := s.repo.GetUsed(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get used si. error: %w", err)
	}
	return data, nil
}

func (s *SIService) Create(ctx context.Context, dto *models.SiDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.instrument.CreateInTx(ctx, tx, dto.Instrument); err != nil {
			return err
		}

		if dto.Verification != nil {
			dto.Verification.InstrumentId = dto.Instrument.Id
			dto.Verification.UserId = dto.Instrument.UserId
			dto.Verification.Status = string(models.InstrumentStatusWork)
			if err := s.verification.Create(ctx, tx, dto.Verification); err != nil {
				return err
			}
		}

		if dto.Location != nil {
			dto.Location.InstrumentId = dto.Instrument.Id
			if err := s.location.Create(ctx, tx, dto.Location); err != nil {
				return err
			}
		}
		return nil
	})
	// // 1. Начинаем транзакцию
	// tx, err := s.repo.BeginTx(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to begin transaction: %w", err)
	// }

	// // 2. Гарантируем откат при любой ошибке или панике
	// defer func() {
	// 	if p := recover(); p != nil {
	// 		tx.Rollback(ctx)
	// 		panic(p)
	// 	} else if err != nil {
	// 		tx.Rollback(ctx)
	// 	}
	// }()

	// // 3. Выполняем операции, передавая объект транзакции (tx)
	// if err := s.instrument.CreateInTx(ctx, tx, dto.Instrument); err != nil {
	// 	return err
	// }

	// if dto.Verification != nil {
	// 	dto.Verification.InstrumentId = dto.Instrument.Id
	// 	dto.Verification.UserId = dto.Instrument.UserId
	// 	dto.Verification.Status = string(models.InstrumentStatusWork)
	// 	if err := s.verification.CreateInTx(ctx, tx, dto.Verification); err != nil {
	// 		return err
	// 	}
	// }

	// if dto.Location != nil {
	// 	dto.Location.InstrumentId = dto.Instrument.Id
	// 	if err := s.location.CreateInTx(ctx, tx, dto.Location); err != nil {
	// 		return err
	// 	}
	// }

	// // 4. Фиксируем изменения.
	// if err := tx.Commit(ctx); err != nil {
	// 	return fmt.Errorf("failed to commit transaction: %w", err)
	// }
	// return nil

	// if err := s.instrument.Create(ctx, dto.Instrument); err != nil {
	// 	return err
	// }
	// if dto.Verification != nil {
	// 	dto.Verification.InstrumentId = dto.Instrument.Id
	// 	dto.Verification.UserId = dto.Instrument.UserId
	// 	dto.Verification.Status = string(models.InstrumentStatusWork)
	// 	if err := s.verification.Create(ctx, dto.Verification); err != nil {
	// 		s.instrument.Delete(ctx, dto.Instrument.Id)
	// 		return err
	// 	}
	// }
	// if dto.Location != nil {
	// 	dto.Location.InstrumentId = dto.Instrument.Id
	// 	if err := s.location.Create(ctx, dto.Location); err != nil {
	// 		s.instrument.Delete(ctx, dto.Instrument.Id)
	// 		return err
	// 	}
	// }

	// return nil
}

func (s *SIService) Update(ctx context.Context, dto *models.SiDTO) error {
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.instrument.Update(ctx, tx, dto.Instrument); err != nil {
			return err
		}

		if dto.Verification != nil {
			if dto.Verification.Id == "" {
				dto.Verification.InstrumentId = dto.Instrument.Id
				dto.Verification.UserId = dto.Instrument.UserId
				dto.Verification.Status = string(models.InstrumentStatusWork)
				if err := s.verification.Create(ctx, tx, dto.Verification); err != nil {
					return err
				}
			} else {
				if err := s.verification.Update(ctx, tx, dto.Verification); err != nil {
					return err
				}
			}
		}
		// if dto.Location != nil {
		// 	if err := s.location.Update(ctx, dto.Location); err != nil {
		// 		return err
		// 	}
		// }
		return nil
	})
}

func (s *SIService) ChangePosition(ctx context.Context, dto *models.ChangePositionDTO) error {
	if err := s.instrument.ChangePosition(ctx, dto); err != nil {
		return err
	}
	return nil
}

func (s *SIService) Delete(ctx context.Context, dto *models.DeleteSiDTO) error {
	if err := s.instrument.Delete(ctx, dto.Id); err != nil {
		return fmt.Errorf("failed to delete si. error: %w", err)
	}
	return nil
}
