package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

type VerificationService struct {
	repo       repository.Verification
	txManager  TransactionManager
	verDocs    VerificationDoc
	instrument Instrument
	docs       Document
}

type VerificationDeps struct {
	Repo       repository.Verification
	TxManager  TransactionManager
	VerDocs    VerificationDoc
	Instrument Instrument
	Docs       Document
}

func NewVerificationService(deps *VerificationDeps) *VerificationService {
	return &VerificationService{
		repo:       deps.Repo,
		txManager:  deps.TxManager,
		verDocs:    deps.VerDocs,
		instrument: deps.Instrument,
		docs:       deps.Docs,
	}
}

type Verification interface {
	Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error)
	GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error)
	Create(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	// Create(ctx context.Context, dto *models.VerificationDTO) error
	CreateSeveral(ctx context.Context, dto []*models.VerificationDTO) error
	Update(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error
	Delete(ctx context.Context, dto *models.DeleteVerificationDTO) error
}

func (s *VerificationService) Get(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error) {
	data, err := s.repo.Get(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get verification by instrument. error: %w", err)
	}

	docs, err := s.verDocs.GetGrouped(ctx, &models.GetGroupedVerificationDocsDTO{InstrumentId: req.InstrumentId})
	if err != nil {
		return nil, err
	}

	for i := range data {
		group, exists := docs.Groups[data[i].Id]
		if exists {
			data[i].Docs = group.Docs
		}
	}
	return data, nil
}

func (s *VerificationService) GetLast(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
	data, err := s.repo.GetLast(ctx, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get last verification. error: %w", err)
	}

	docs, err := s.verDocs.Get(ctx, &models.GetVerificationDocsDTO{VerificationId: data.Id})
	if err != nil {
		return nil, err
	}
	data.Docs = docs

	return data, nil
}

func (s *VerificationService) Create(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.ExecuteInTx(ctx, func(newTx postgres.Tx) error {
			return s.executeCreate(ctx, newTx, dto)
		})
	}

	// Если транзакция передана, используем её
	return s.executeCreate(ctx, tx, dto)
}
func (s *VerificationService) executeCreate(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if err := s.repo.CreateInTx(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to create verification. error: %w", err)
	}

	for i := range dto.Docs {
		dto.Docs[i].VerificationId = dto.Id
	}

	if err := s.verDocs.CreateSeveral(ctx, tx, dto.Docs); err != nil {
		s.Delete(ctx, &models.DeleteVerificationDTO{Id: dto.Id})
		return err
	}

	if len(dto.Docs) > 0 {
		pathDTO := &models.PathParts{
			InstrumentId: dto.InstrumentId,
			Group:        "verifications",
			UserId:       dto.UserId,
		}
		if err := s.docs.ChangePath(ctx, pathDTO); err != nil {
			return err
		}
	}

	instDTO := &models.UpdateStatus{
		Id:     dto.InstrumentId,
		Status: models.InstrumentStatus(dto.Status),
	}

	if err := s.instrument.ChangeStatus(ctx, tx, instDTO); err != nil {
		return err
	}

	return nil
}

// func (s *VerificationService) Create(ctx context.Context, dto *models.VerificationDTO) error {
// 	if err := s.repo.Create(ctx, dto); err != nil {
// 		return fmt.Errorf("failed to create verification. error: %w", err)
// 	}

// 	for i := range dto.Docs {
// 		dto.Docs[i].VerificationId = dto.Id
// 	}
// 	if err := s.verDocs.CreateSeveral(ctx, nil, dto.Docs); err != nil {
// 		s.Delete(ctx, &models.DeleteVerificationDTO{Id: dto.Id})
// 		return err
// 	}

// 	if len(dto.Docs) > 0 {
// 		pathDTO := &models.PathParts{
// 			InstrumentId: dto.InstrumentId,
// 			Group:        "verifications",
// 			UserId:       dto.UserId,
// 		}
// 		if err := s.docs.ChangePath(ctx, pathDTO); err != nil {
// 			return err
// 		}
// 	}

// 	instDTO := &models.UpdateStatus{
// 		Id:     dto.InstrumentId,
// 		Status: models.InstrumentStatus(dto.Status),
// 	}
// 	if err := s.instrument.ChangeStatus(ctx, nil, instDTO); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (s *VerificationService) CreateSeveral(ctx context.Context, dto []*models.VerificationDTO) error {

	if len(dto) == 0 {
		return nil
	}
	return s.txManager.ExecuteInTx(ctx, func(tx postgres.Tx) error {
		if err := s.repo.CreateSeveral(ctx, tx, dto); err != nil {
			return fmt.Errorf("failed to create verifications. error: %w", err)
		}

		docs := []*models.VerificationDocDTO{}
		for i := range dto {
			for j := range dto[i].Docs {
				dto[i].Docs[j].VerificationId = dto[i].Id
			}
			docs = append(docs, dto[i].Docs...)
		}
		if err := s.verDocs.CreateSeveral(ctx, tx, docs); err != nil {
			return err
		}

		return nil
	})
}

func (s *VerificationService) Update(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if tx == nil {
		// Если транзакция не передана, создаем новую
		return s.txManager.ExecuteInTx(ctx, func(newTx postgres.Tx) error {
			return s.executeUpdate(ctx, newTx, dto)
		})
	}

	// Если транзакция передана, используем её
	return s.executeUpdate(ctx, tx, dto)
}
func (s *VerificationService) executeUpdate(ctx context.Context, tx postgres.Tx, dto *models.VerificationDTO) error {
	if err := s.repo.Update(ctx, tx, dto); err != nil {
		return fmt.Errorf("failed to update verification. error: %w", err)
	}

	newDocs := []*models.VerificationDocDTO{}
	updatedDocs := []*models.VerificationDocDTO{}
	for i := range dto.Docs {
		if dto.Docs[i].Id == "" {
			dto.Docs[i].VerificationId = dto.Id
			newDocs = append(newDocs, dto.Docs[i])
		} else {
			updatedDocs = append(updatedDocs, dto.Docs[i])
		}
	}

	if len(newDocs) > 0 {
		if err := s.verDocs.CreateSeveral(ctx, tx, newDocs); err != nil {
			return err
		}
	}
	if len(updatedDocs) > 0 {
		if err := s.verDocs.UpdateSeveral(ctx, tx, updatedDocs); err != nil {
			return err
		}
	}
	return nil
}

func (s *VerificationService) Delete(ctx context.Context, dto *models.DeleteVerificationDTO) error {
	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete verification. error: %w", err)
	}
	return nil
}
