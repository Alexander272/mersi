package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
)

func newReceivingService(loc *fakeLocationSvc, repo *fakeLocationRepo, resp *fakeResponsibleSvc, notif *fakeNotificationSvc, dialog *fakeDialog) *ReceivingService {
	return NewReceivingService(&ReceivingDeps{
		Location:     loc,
		Repo:         repo,
		Responsible:  resp,
		Notification: notif,
		Most:         dialog,
		ActivityLog:  &fakeActivityLogSvc{},
		TxManager:    &fakeTxManager{},
	})
}

func TestReceivingFromApp_NotResponsible(t *testing.T) {
	svc := newReceivingService(
		&fakeLocationSvc{},
		&fakeLocationRepo{},
		&fakeResponsibleSvc{},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	dto := &models.ReceivingDTO{InstrumentIds: []string{"i1"}, Status: constants.LocationStatusUsed, UserId: "u1"}
	err := svc.ReceivingFromApp(context.Background(), dto)
	if !errors.Is(err, models.ErrNotResponsible) {
		t.Fatalf("expected ErrNotResponsible, got %v", err)
	}
}

func TestReceivingFromApp_FilterByDepartments(t *testing.T) {
	repo := &fakeLocationRepo{}
	svc := newReceivingService(
		&fakeLocationSvc{
			selectByDepartmentsFn: func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
				if len(dto.DepartmentIds) != 2 {
					t.Fatalf("expected 2 department ids, got %v", dto.DepartmentIds)
				}
				if dto.Status != constants.LocationStatusMoved {
					t.Fatalf("expected status moved, got %s", dto.Status)
				}
				return []string{"i1"}, nil
			},
		},
		repo,
		&fakeResponsibleSvc{
			getBySSOIdFn: func(ctx context.Context, id string) ([]*models.Responsible, error) {
				return []*models.Responsible{{Id: "r1", DepartmentId: "d1"}, {Id: "r2", DepartmentId: "d2"}}, nil
			},
		},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	dto := &models.ReceivingDTO{InstrumentIds: []string{"i1", "i2", "i3"}, Status: constants.LocationStatusUsed, UserId: "u1"}
	if err := svc.ReceivingFromApp(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivingCalled) != 1 {
		t.Fatalf("expected 1 Receiving call, got %d", len(repo.receivingCalled))
	}
	got := repo.receivingCalled[0]
	if len(got.InstrumentIds) != 1 || got.InstrumentIds[0] != "i1" {
		t.Fatalf("expected filtered ids [i1], got %v", got.InstrumentIds)
	}
}

func TestReceivingFromApp_AllInDepartmentsKeepsIds(t *testing.T) {
	repo := &fakeLocationRepo{}
	svc := newReceivingService(
		&fakeLocationSvc{
			selectByDepartmentsFn: func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
				return []string{"i1", "i2", "i3"}, nil
			},
		},
		repo,
		&fakeResponsibleSvc{
			getBySSOIdFn: func(ctx context.Context, id string) ([]*models.Responsible, error) {
				return []*models.Responsible{{Id: "r1", DepartmentId: "d1"}}, nil
			},
		},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	dto := &models.ReceivingDTO{InstrumentIds: []string{"i1", "i2", "i3"}, Status: constants.LocationStatusUsed, UserId: "u1"}
	if err := svc.ReceivingFromApp(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := repo.receivingCalled[0]
	if len(got.InstrumentIds) != 3 {
		t.Fatalf("expected all 3 ids kept, got %v", got.InstrumentIds)
	}
}

func TestReceivingFromApp_EmptyAfterFilter(t *testing.T) {
	svc := newReceivingService(
		&fakeLocationSvc{
			selectByDepartmentsFn: func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
				return nil, nil
			},
		},
		&fakeLocationRepo{},
		&fakeResponsibleSvc{
			getBySSOIdFn: func(ctx context.Context, id string) ([]*models.Responsible, error) {
				return []*models.Responsible{{Id: "r1", DepartmentId: "d1"}}, nil
			},
		},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	dto := &models.ReceivingDTO{InstrumentIds: []string{"i1", "i2"}, Status: constants.LocationStatusUsed, UserId: "u1"}
	err := svc.ReceivingFromApp(context.Background(), dto)
	if !errors.Is(err, models.ErrCannotConfirmReceipt) {
		t.Fatalf("expected ErrCannotConfirmReceipt, got %v", err)
	}
}

func TestReceivingFromApp_EmptyIdsWithoutFilter(t *testing.T) {
	svc := newReceivingService(
		&fakeLocationSvc{},
		&fakeLocationRepo{},
		&fakeResponsibleSvc{},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	dto := &models.ReceivingDTO{InstrumentIds: nil, Status: constants.LocationStatusReserve, UserId: "u1"}
	err := svc.ReceivingFromApp(context.Background(), dto)
	if !errors.Is(err, models.ErrCannotConfirmReceipt) {
		t.Fatalf("expected ErrCannotConfirmReceipt, got %v", err)
	}
}

func TestReceivingFromApp_RepoError(t *testing.T) {
	repoErr := errors.New("db is down")
	repo := &fakeLocationRepo{receivingFn: func(ctx context.Context, dto *models.ReceivingDTO) error { return repoErr }}
	svc := newReceivingService(
		&fakeLocationSvc{},
		repo,
		&fakeResponsibleSvc{},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	dto := &models.ReceivingDTO{InstrumentIds: []string{"i1"}, Status: constants.LocationStatusReserve, UserId: "u1"}
	err := svc.ReceivingFromApp(context.Background(), dto)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped repo error, got %v", err)
	}
}

func TestReceivingFromChannel_ParsesState(t *testing.T) {
	repo := &fakeLocationRepo{}
	notif := &fakeNotificationSvc{
		checkReceivingFn: func(ctx context.Context, dto *models.DialogResponse) error { return nil },
	}
	svc := newReceivingService(&fakeLocationSvc{}, repo, &fakeResponsibleSvc{}, notif, &fakeDialog{})

	dto := &models.DialogResponse{
		State:      "PostId:abc&Status:moved",
		Submission: map[string]bool{"i1": true, "i2": false, "i3": true},
	}
	if err := svc.ReceivingFromChannel(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.receivingCalled) != 1 {
		t.Fatalf("expected 1 Receiving call, got %d", len(repo.receivingCalled))
	}
	got := repo.receivingCalled[0]
	if got.Status != constants.LocationStatusMoved {
		t.Fatalf("expected status moved, got %s", got.Status)
	}
	if !got.HasConfirmed {
		t.Fatal("expected HasConfirmed to be true")
	}
	if len(got.InstrumentIds) != 2 {
		t.Fatalf("expected 2 accepted instruments, got %v", got.InstrumentIds)
	}
	seen := make(map[string]bool)
	for _, id := range got.InstrumentIds {
		seen[id] = true
	}
	if !seen["i1"] || !seen["i3"] || seen["i2"] {
		t.Fatalf("expected only i1 and i3, got %v", got.InstrumentIds)
	}
}

func TestReceivingFromChannel_DefaultStatus(t *testing.T) {
	repo := &fakeLocationRepo{}
	svc := newReceivingService(&fakeLocationSvc{}, repo, &fakeResponsibleSvc{}, &fakeNotificationSvc{}, &fakeDialog{})

	dto := &models.DialogResponse{
		State:      "PostId:abc",
		Submission: map[string]bool{"i1": true},
	}
	if err := svc.ReceivingFromChannel(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.receivingCalled[0].Status != constants.LocationStatusUsed {
		t.Fatalf("expected default status used, got %s", repo.receivingCalled[0].Status)
	}
}

func TestReceivingDialogOpen(t *testing.T) {
	dialog := &fakeDialog{}
	svc := newReceivingService(&fakeLocationSvc{}, &fakeLocationRepo{}, &fakeResponsibleSvc{}, &fakeNotificationSvc{}, dialog)

	action := &models.PostAction{TriggerId: "t1", PostId: "p1"}
	if err := svc.ReceivingDialogOpen(context.Background(), action); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dialog.action != action {
		t.Fatal("expected action to be passed to dialog")
	}
}

func TestForcedReceipt_InstrumentReceived(t *testing.T) {
	repo := &fakeLocationRepo{
		forcedReceiptFn: func(ctx context.Context, dto *models.ForcedReceiptDTO) error {
			return models.ErrNoRows
		},
	}
	svc := newReceivingService(&fakeLocationSvc{}, repo, &fakeResponsibleSvc{}, &fakeNotificationSvc{}, &fakeDialog{})

	err := svc.ForcedReceipt(context.Background(), &models.ForcedReceiptDTO{InstrumentId: "i1", Actor: &models.Actor{ID: "u1", Name: "User"}})
	if !errors.Is(err, models.ErrInstrumentReceived) {
		t.Fatalf("expected ErrInstrumentReceived, got %v", err)
	}
}

func TestForcedReceipt_SuccessLogsActivity(t *testing.T) {
	repo := &fakeLocationRepo{}
	activity := &fakeActivityLogSvc{}
	svc := newReceivingService(
		&fakeLocationSvc{
			getByIdFn: func(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
				return &models.Location{Id: "l1", Place: "dep1", Status: constants.LocationStatusMoved}, nil
			},
		},
		repo,
		&fakeResponsibleSvc{},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)
	svc.activityLog = activity

	dto := &models.ForcedReceiptDTO{InstrumentId: "i1", Actor: &models.Actor{ID: "u1", Name: "User"}}
	if err := svc.ForcedReceipt(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(repo.forcedCalled) != 1 {
		t.Fatalf("expected 1 ForcedReceipt call, got %d", len(repo.forcedCalled))
	}
	if len(activity.logs) != 1 {
		t.Fatalf("expected 1 activity log, got %d", len(activity.logs))
	}
	if activity.logs[0].Action != "FORCED_RECEIPT" {
		t.Fatalf("expected FORCED_RECEIPT action, got %s", activity.logs[0].Action)
	}
}

func TestForcedReceipt_GetOldDataError(t *testing.T) {
	repoErr := errors.New("select failed")
	svc := newReceivingService(
		&fakeLocationSvc{
			getByIdFn: func(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
				return nil, repoErr
			},
		},
		&fakeLocationRepo{},
		&fakeResponsibleSvc{},
		&fakeNotificationSvc{},
		&fakeDialog{},
	)

	err := svc.ForcedReceipt(context.Background(), &models.ForcedReceiptDTO{InstrumentId: "i1"})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestForcedReceiptAll(t *testing.T) {
	repo := &fakeLocationRepo{}
	svc := newReceivingService(&fakeLocationSvc{}, repo, &fakeResponsibleSvc{}, &fakeNotificationSvc{}, &fakeDialog{})
	if err := svc.ForcedReceiptAll(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
