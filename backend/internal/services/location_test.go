package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
)

func newLocationService(repo *fakeLocationRepo, resp *fakeResponsibleSvc) *LocationService {
	return NewLocationService(&LocationDeps{
		Repo:        repo,
		TxManager:   &fakeTxManager{},
		Responsible: resp,
		ActivityLog: &fakeActivityLogSvc{},
	})
}

func locationDTO() *models.LocationDTO {
	return &models.LocationDTO{
		Id:           "l1",
		InstrumentId: "i1",
		DepartmentId: "d1",
		Status:       constants.LocationStatusMoved,
		Actor:        &models.Actor{ID: "u1", Name: "User"},
	}
}

func TestExecuteCreate_NoResponsible(t *testing.T) {
	svc := newLocationService(&fakeLocationRepo{}, &fakeResponsibleSvc{})
	err := svc.Create(context.Background(), &fakeTx{}, locationDTO())
	if !errors.Is(err, models.ErrNoResponsible) {
		t.Fatalf("expected ErrNoResponsible, got %v", err)
	}
}

func TestExecuteCreate_NoChannel(t *testing.T) {
	svc := newLocationService(&fakeLocationRepo{}, &fakeResponsibleSvc{
		getWithChannelFn: func(ctx context.Context, req *models.GetResponsibleDTO) ([]*models.ResponsibleWithChannel, error) {
			if req.DepartmentId != "d1" {
				t.Fatalf("expected department d1, got %s", req.DepartmentId)
			}
			return []*models.ResponsibleWithChannel{{Id: "r1", DepartmentId: "d1", ChannelId: ""}}, nil
		},
	})
	err := svc.Create(context.Background(), &fakeTx{}, locationDTO())
	if !errors.Is(err, models.ErrNoChannel) {
		t.Fatalf("expected ErrNoChannel, got %v", err)
	}
}

func TestExecuteCreate_Success(t *testing.T) {
	repo := &fakeLocationRepo{}
	activity := &fakeActivityLogSvc{}
	svc := NewLocationService(&LocationDeps{
		Repo:        repo,
		TxManager:   &fakeTxManager{},
		Responsible: &fakeResponsibleSvc{
			getWithChannelFn: func(ctx context.Context, req *models.GetResponsibleDTO) ([]*models.ResponsibleWithChannel, error) {
				return []*models.ResponsibleWithChannel{{Id: "r1", DepartmentId: "d1", ChannelId: "ch1"}}, nil
			},
		},
		ActivityLog: activity,
	})

	dto := locationDTO()
	if err := svc.Create(context.Background(), &fakeTx{}, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(activity.logs) != 1 || activity.logs[0].Action != "CREATE" {
		t.Fatalf("expected CREATE activity log, got %+v", activity.logs)
	}
}

func TestExecuteCreate_OnlyChecksResponsibleForMoved(t *testing.T) {
	repo := &fakeLocationRepo{}
	activity := &fakeActivityLogSvc{}
	svc := NewLocationService(&LocationDeps{
		Repo:        repo,
		TxManager:   &fakeTxManager{},
		Responsible: &fakeResponsibleSvc{},
		ActivityLog: activity,
	})

	dto := locationDTO()
	dto.Status = constants.LocationStatusUsed
	dto.DepartmentId = ""
	if err := svc.Create(context.Background(), &fakeTx{}, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSeveral_Empty(t *testing.T) {
	svc := newLocationService(&fakeLocationRepo{}, &fakeResponsibleSvc{})
	isFull, err := svc.CreateSeveral(context.Background(), nil)
	if isFull {
		t.Fatal("expected isFull=false for empty input")
	}
	if !errors.Is(err, models.ErrCannotMoveInstrument) {
		t.Fatalf("expected ErrCannotMoveInstrument, got %v", err)
	}
}

func TestCreateSeveral_NotResponsible(t *testing.T) {
	svc := newLocationService(&fakeLocationRepo{}, &fakeResponsibleSvc{})
	dto := []*models.LocationDTO{{
		InstrumentId: "i1",
		PersonId:     "",
		NeedConfirm:  true,
		UserId:       "u1",
	}}
	isFull, err := svc.CreateSeveral(context.Background(), dto)
	if isFull {
		t.Fatal("expected isFull=false")
	}
	if !errors.Is(err, models.ErrNotResponsible) {
		t.Fatalf("expected ErrNotResponsible, got %v", err)
	}
}

func TestCreateSeveral_FiltersByDepartments(t *testing.T) {
	repo := &fakeLocationRepo{
		selectByDeptFn: func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
			return []string{"i1"}, nil
		},
	}
	svc := newLocationService(repo, &fakeResponsibleSvc{
		getBySSOIdFn: func(ctx context.Context, id string) ([]*models.Responsible, error) {
			return []*models.Responsible{{Id: "r1", DepartmentId: "d1"}}, nil
		},
	})

	dto := []*models.LocationDTO{
		{InstrumentId: "i1", PersonId: "", NeedConfirm: true, UserId: "u1"},
		{InstrumentId: "i2", PersonId: "", NeedConfirm: true, UserId: "u1"},
	}
	isFull, err := svc.CreateSeveral(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isFull {
		t.Fatal("expected isFull=false when only some instruments match")
	}
	if len(repo.createSeveralCalled) != 1 {
		t.Fatalf("expected 1 CreateSeveral call, got %d", len(repo.createSeveralCalled))
	}
	got := repo.createSeveralCalled[0]
	if len(got) != 1 || got[0].InstrumentId != "i1" {
		t.Fatalf("expected only i1 to be created, got %+v", got)
	}
}

func TestCreateSeveral_AllFilteredOut(t *testing.T) {
	svc := newLocationService(&fakeLocationRepo{
		selectByDeptFn: func(ctx context.Context, dto *models.SelectByDepsDTO) ([]string, error) {
			return nil, nil
		},
	}, &fakeResponsibleSvc{
		getBySSOIdFn: func(ctx context.Context, id string) ([]*models.Responsible, error) {
			return []*models.Responsible{{Id: "r1", DepartmentId: "d1"}}, nil
		},
	})

	dto := []*models.LocationDTO{
		{InstrumentId: "i1", PersonId: "", NeedConfirm: true, UserId: "u1"},
	}
	isFull, err := svc.CreateSeveral(context.Background(), dto)
	if isFull {
		t.Fatal("expected isFull=false")
	}
	if !errors.Is(err, models.ErrCannotMoveInstrument) {
		t.Fatalf("expected ErrCannotMoveInstrument, got %v", err)
	}
}

func TestCreateSeveral_SkipsFilterWithoutConfirm(t *testing.T) {
	repo := &fakeLocationRepo{}
	svc := newLocationService(repo, &fakeResponsibleSvc{})

	dto := []*models.LocationDTO{
		{InstrumentId: "i1", PersonId: "p1", NeedConfirm: false, UserId: "u1"},
	}
	isFull, err := svc.CreateSeveral(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isFull {
		t.Fatal("expected isFull=true when no filtering happens")
	}
	if len(repo.createSeveralCalled) != 1 {
		t.Fatalf("expected 1 CreateSeveral call, got %d", len(repo.createSeveralCalled))
	}
}
