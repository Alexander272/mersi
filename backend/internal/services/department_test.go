package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/models"
)

func newDepartmentFixture() (*DepartmentService, *fakeDepartmentRepo, *fakeLocationSvc) {
	repo := &fakeDepartmentRepo{}
	location := &fakeLocationSvc{}
	return NewDepartmentService(repo, location), repo, location
}

func TestDepartmentGetAll_NilBecomesEmpty(t *testing.T) {
	svc, repo, _ := newDepartmentFixture()
	repo.getAllFn = func(ctx context.Context, req *models.GetDepartmentsDTO) ([]*models.Department, error) {
		return nil, nil
	}

	data, err := svc.GetAll(context.Background(), &models.GetDepartmentsDTO{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil || len(data) != 0 {
		t.Fatalf("expected empty slice, got %+v", data)
	}
}

func TestDepartmentGetById_ReturnsErrNoRows(t *testing.T) {
	svc, _, _ := newDepartmentFixture()

	_, err := svc.GetById(context.Background(), &models.GetDepartmentByIdDTO{Id: "dept-1"})
	if !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestDepartmentDelete_BlockedWhenInstrumentsUsed(t *testing.T) {
	svc, repo, location := newDepartmentFixture()

	location.getUsedByDeptFn = func(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error) {
		return []*models.Location{{Id: "loc-1"}}, nil
	}

	if err := svc.Delete(context.Background(), "dept-1"); !errors.Is(err, models.ErrDeleteDepartmentHasInstrument) {
		t.Fatalf("expected ErrDeleteDepartmentHasInstrument, got %v", err)
	}
	if len(location.setDepartmentCalls) != 0 {
		t.Fatalf("expected no SetDepartment, got %+v", location.setDepartmentCalls)
	}
	if len(repo.deleted) != 0 {
		t.Fatalf("expected no repo.Delete, got %+v", repo.deleted)
	}
}

func TestDepartmentDelete_OrderedSetDepartmentThenDelete(t *testing.T) {
	svc, repo, location := newDepartmentFixture()

	if err := svc.Delete(context.Background(), "dept-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(location.setDepartmentCalls) != 1 || location.setDepartmentCalls[0] != "dept-1" {
		t.Fatalf("expected SetDepartment with dept-1, got %+v", location.setDepartmentCalls)
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "dept-1" {
		t.Fatalf("expected repo.Delete with dept-1, got %+v", repo.deleted)
	}
}

func TestDepartmentDelete_PropagatesLocationError(t *testing.T) {
	svc, repo, location := newDepartmentFixture()

	location.getUsedByDeptFn = func(ctx context.Context, dto *models.GetLocationByDepartmentDTO) ([]*models.Location, error) {
		return nil, models.ErrForbidden
	}

	if err := svc.Delete(context.Background(), "dept-1"); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if len(repo.deleted) != 0 {
		t.Fatalf("expected no repo.Delete, got %+v", repo.deleted)
	}
}

func TestDepartmentDelete_PropagatesRepoError(t *testing.T) {
	svc, repo, _ := newDepartmentFixture()

	repo.deleteFn = func(ctx context.Context, id string) error {
		return models.ErrForbidden
	}

	if err := svc.Delete(context.Background(), "dept-1"); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDepartmentCreate_PassesThrough(t *testing.T) {
	svc, repo, _ := newDepartmentFixture()
	repo.createFn = func(ctx context.Context, dto *models.DepartmentDTO) (string, error) {
		if dto.Name != "Department" {
			t.Fatalf("unexpected dto: %+v", dto)
		}
		return "dept-1", nil
	}

	id, err := svc.Create(context.Background(), &models.DepartmentDTO{Name: "Department"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != "dept-1" {
		t.Fatalf("unexpected id: %s", id)
	}
}
