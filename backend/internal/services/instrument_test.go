package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

func newInstrumentServiceFixture() (*InstrumentService, *fakeInstrumentRepo, *fakeDocumentSvc, *fakeActivityLogSvc, *fakeTxManager) {
	repo := &fakeInstrumentRepo{}
	docs := &fakeDocumentSvc{}
	log := &fakeActivityLogSvc{}
	txManager := &fakeTxManager{}
	return NewInstrumentService(repo, docs, log, txManager), repo, docs, log, txManager
}

func testInstrumentDTO() *models.InstrumentDTO {
	return &models.InstrumentDTO{
		Id:              "ins-1",
		Name:            "SI-1",
		UserId:          "user-1",
		ActOfEnteringId: "",
		Actor:           &models.Actor{ID: "user-1", Name: "User"},
	}
}

func TestInstrumentCreate_WithoutTx_ExecutesInTransaction(t *testing.T) {
	svc, repo, docs, log, txManager := newInstrumentServiceFixture()
	dto := testInstrumentDTO()

	txUsed := false
	txManager.execute = func(ctx context.Context, fn func(tx postgres.Tx) error) error {
		txUsed = true
		return fn(&fakeTx{})
	}

	err := svc.Create(context.Background(), nil, dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !txUsed {
		t.Fatal("expected transaction manager to be used for nil tx")
	}
	if len(repo.created) != 1 {
		t.Fatalf("expected repo.CreateInTx to be called once, got %d", len(repo.created))
	}
	if len(log.logs) != 1 || log.logs[0].Action != "CREATE" {
		t.Fatalf("expected CREATE activity log, got %+v", log.logs)
	}
	if len(docs.pathChanges) != 0 {
		t.Fatalf("ChangePath must not be called when ActOfEnteringId is empty, got %+v", docs.pathChanges)
	}
}

func TestInstrumentCreate_WithTx_SkipsTransactionManager(t *testing.T) {
	svc, repo, _, log, _ := newInstrumentServiceFixture()
	dto := testInstrumentDTO()

	created := &fakeTx{}
	if err := svc.Create(context.Background(), created, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.created) != 1 {
		t.Fatalf("expected repo.CreateInTx to be called once, got %d", len(repo.created))
	}
	if created.committed {
		t.Fatal("direct tx path must not commit through manager")
	}
	if len(log.logs) != 1 || log.logs[0].TableName != "instruments" {
		t.Fatalf("expected instruments activity log, got %+v", log.logs)
	}
}

func TestInstrumentCreate_CallsChangePathWhenActOfEntering(t *testing.T) {
	svc, repo, docs, _, _ := newInstrumentServiceFixture()
	dto := testInstrumentDTO()
	dto.ActOfEnteringId = "act-1"

	if err := svc.Create(context.Background(), nil, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.created) != 1 {
		t.Fatalf("expected repo.CreateInTx to be called once, got %d", len(repo.created))
	}
	if len(docs.pathChanges) != 1 {
		t.Fatalf("expected ChangePath to be called, got %d", len(docs.pathChanges))
	}
	p := docs.pathChanges[0]
	if p.InstrumentId != dto.Id || p.Group != "act" || p.UserId != dto.UserId || !p.IdWasEmpty {
		t.Fatalf("unexpected path parts: %+v", p)
	}
}

func TestInstrumentCreate_PropagatesChangePathError(t *testing.T) {
	svc, _, docs, _, _ := newInstrumentServiceFixture()
	dto := testInstrumentDTO()
	dto.ActOfEnteringId = "act-1"
	docs.changePathFn = func(ctx context.Context, d *models.PathParts) error {
		return models.ErrNotValid
	}

	if err := svc.Create(context.Background(), nil, dto); !errors.Is(err, models.ErrNotValid) {
		t.Fatalf("expected ErrNotValid, got %v", err)
	}
}

func TestInstrumentCreate_PropagatesRepoError(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	repo.createInTxFn = func(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
		return models.ErrNoRows
	}

	if err := svc.Create(context.Background(), nil, testInstrumentDTO()); !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestInstrumentCreateSeveral_EmptyIsNoop(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()

	if err := svc.CreateSeveral(context.Background(), []*models.InstrumentDTO{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.createdSeveral) != 0 {
		t.Fatalf("expected no repo calls, got %d", len(repo.createdSeveral))
	}
}

func TestInstrumentCreateSeveral_PassesThrough(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	dto := []*models.InstrumentDTO{{Name: "SI-1"}, {Name: "SI-2"}}

	if err := svc.CreateSeveral(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.createdSeveral) != 1 || len(repo.createdSeveral[0]) != 2 {
		t.Fatalf("expected 2 instruments to be created, got %+v", repo.createdSeveral)
	}
}

func TestInstrumentUpdate_LogsWithOldData(t *testing.T) {
	svc, repo, _, log, _ := newInstrumentServiceFixture()
	old := &models.Instrument{Id: "ins-1", Name: "Old name"}
	repo.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return old, nil
	}
	dto := testInstrumentDTO()
	dto.Name = "New name"

	if err := svc.Update(context.Background(), nil, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.updated) != 1 || repo.updated[0] != dto {
		t.Fatalf("expected repo.Update with dto, got %+v", repo.updated)
	}
	if len(log.logs) != 1 || log.logs[0].Action != "UPDATE" {
		t.Fatalf("expected UPDATE activity log, got %+v", log.logs)
	}
	if log.logs[0].OldValue != old {
		t.Fatalf("expected old value in log, got %+v", log.logs[0])
	}
}

func TestInstrumentUpdate_PropagatesGetByIdError(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	repo.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return nil, models.ErrForbidden
	}

	if err := svc.Update(context.Background(), nil, testInstrumentDTO()); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestInstrumentDelete_LogsDelete(t *testing.T) {
	svc, repo, _, log, _ := newInstrumentServiceFixture()
	repo.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return &models.Instrument{Id: "ins-1", Name: "SI-1"}, nil
	}

	if err := svc.Delete(context.Background(), &models.DeleteSiDTO{Id: "ins-1", Actor: &models.Actor{ID: "user-1"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.deletedIds) != 1 || repo.deletedIds[0] != "ins-1" {
		t.Fatalf("expected repo.Delete with ins-1, got %+v", repo.deletedIds)
	}
	if len(log.logs) != 1 || log.logs[0].Action != "DELETE" {
		t.Fatalf("expected DELETE activity log, got %+v", log.logs)
	}
}

func TestInstrumentDelete_ToleratesErrNoRowsOnGetById(t *testing.T) {
	svc, repo, _, log, _ := newInstrumentServiceFixture()
	repo.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return nil, models.ErrNoRows
	}

	if err := svc.Delete(context.Background(), &models.DeleteSiDTO{Id: "ins-1", Actor: &models.Actor{ID: "user-1"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.deletedIds) != 1 {
		t.Fatalf("expected delete to proceed, got %+v", repo.deletedIds)
	}
	if len(log.logs) != 0 {
		t.Fatalf("expected no activity log when instrument not found, got %+v", log.logs)
	}
}

func TestInstrumentDelete_PropagatesGetByIdError(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	repo.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return nil, models.ErrForbidden
	}

	if err := svc.Delete(context.Background(), &models.DeleteSiDTO{Id: "ins-1", Actor: &models.Actor{ID: "user-1"}}); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if len(repo.deletedIds) != 0 {
		t.Fatalf("expected no delete, got %+v", repo.deletedIds)
	}
}

func TestInstrumentGetById_ReturnsErrNoRows(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	repo.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return nil, models.ErrNoRows
	}

	if _, err := svc.GetById(context.Background(), &models.GetInstrumentByIdDTO{Id: "ins-1"}); !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestInstrumentChangePosition_PassesThrough(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	changePositionCalled := false
	repo.changePositionFn = func(ctx context.Context, dto *models.ChangePositionDTO) error {
		changePositionCalled = true
		if dto.NewPosition != 2 {
			t.Fatalf("expected new position 2, got %d", dto.NewPosition)
		}
		return nil
	}

	if err := svc.ChangePosition(context.Background(), &models.ChangePositionDTO{SectionId: "sec-1", NewPosition: 2, OldPosition: 1}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changePositionCalled {
		t.Fatal("expected repo.ChangePosition to be called")
	}
}

func TestInstrumentChangeStatus_PassesThrough(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	statusCalled := false
	repo.changeStatusFn = func(ctx context.Context, tx postgres.Tx, dto *models.UpdateStatus) error {
		statusCalled = true
		if dto.Status != models.InstrumentStatusRepair {
			t.Fatalf("expected repair status, got %s", dto.Status)
		}
		return nil
	}

	if err := svc.ChangeStatus(context.Background(), nil, &models.UpdateStatus{Id: "ins-1", Status: models.InstrumentStatusRepair}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !statusCalled {
		t.Fatal("expected repo.ChangeStatus to be called")
	}
}

func TestInstrumentChangeSeveralStatuses_PassesThrough(t *testing.T) {
	svc, repo, _, _, _ := newInstrumentServiceFixture()
	severalCalled := false
	repo.changeSeveralStatusFn = func(ctx context.Context, dto []*models.UpdateStatus) error {
		severalCalled = true
		if len(dto) != 1 {
			t.Fatalf("expected 1 status update, got %d", len(dto))
		}
		return nil
	}

	if err := svc.ChangeSeveralStatuses(context.Background(), []*models.UpdateStatus{{Id: "ins-1"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !severalCalled {
		t.Fatal("expected repo.ChangeSeveralStatuses to be called")
	}
}
