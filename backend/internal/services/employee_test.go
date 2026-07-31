package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/models"
)

func newEmployeeFixture() (*EmployeeService, *fakeEmployeeRepo, *fakeLocationSvc) {
	repo := &fakeEmployeeRepo{}
	location := &fakeLocationSvc{}
	return NewEmployeeService(repo, location), repo, location
}

func TestEmployeeGetAll_NilBecomesEmpty(t *testing.T) {
	svc, repo, _ := newEmployeeFixture()
	repo.getAllFn = func(ctx context.Context, req *models.GetEmployeesDTO) ([]*models.Employee, error) {
		return nil, nil
	}

	data, err := svc.GetAll(context.Background(), &models.GetEmployeesDTO{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil || len(data) != 0 {
		t.Fatalf("expected empty slice, got %+v", data)
	}
}

func TestEmployeeGetUnique_NilBecomesEmpty(t *testing.T) {
	svc, repo, _ := newEmployeeFixture()
	repo.getUniqueFn = func(ctx context.Context, dto *models.GetUniqueEmployeeDTO) ([]*models.Employee, error) {
		return nil, nil
	}

	data, err := svc.GetUnique(context.Background(), &models.GetUniqueEmployeeDTO{Realm: "realm-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data == nil || len(data) != 0 {
		t.Fatalf("expected empty slice, got %+v", data)
	}
}

func TestEmployeeGetBySSOId_MapsErrNoRows(t *testing.T) {
	svc, _, _ := newEmployeeFixture()

	_, err := svc.GetBySSOId(context.Background(), "sso-1")
	if !errors.Is(err, models.ErrEmployeeNotFound) {
		t.Fatalf("expected ErrEmployeeNotFound, got %v", err)
	}
}

func TestEmployeeGetById_MapsErrNoRows(t *testing.T) {
	svc, _, _ := newEmployeeFixture()

	_, err := svc.GetById(context.Background(), "emp-1")
	if !errors.Is(err, models.ErrEmployeeNotFound) {
		t.Fatalf("expected ErrEmployeeNotFound, got %v", err)
	}
}

func TestEmployeeCreate_DuplicateName(t *testing.T) {
	svc, repo, _ := newEmployeeFixture()

	repo.getByNameFn = func(ctx context.Context, req *models.GetEmployeeByNameDTO) (*models.Employee, error) {
		return &models.Employee{Id: "emp-1", Name: req.Name}, nil
	}

	err := svc.Create(context.Background(), &models.WriteEmployeeDTO{Name: "Ivan", DepartmentId: "dept-1"})
	if !errors.Is(err, models.ErrEmployeeAlreadyExists) {
		t.Fatalf("expected ErrEmployeeAlreadyExists, got %v", err)
	}
	if len(repo.created) != 0 {
		t.Fatalf("expected no repo.Create, got %+v", repo.created)
	}
}

func TestEmployeeCreate_Success(t *testing.T) {
	svc, repo, _ := newEmployeeFixture()

	if err := svc.Create(context.Background(), &models.WriteEmployeeDTO{Name: "Ivan", DepartmentId: "dept-1"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repo.created) != 1 || repo.created[0].Name != "Ivan" {
		t.Fatalf("expected repo.Create with Ivan, got %+v", repo.created)
	}
}

func TestEmployeeCreate_PropagatesGetByNameError(t *testing.T) {
	svc, repo, _ := newEmployeeFixture()

	repo.getByNameFn = func(ctx context.Context, req *models.GetEmployeeByNameDTO) (*models.Employee, error) {
		return nil, models.ErrForbidden
	}

	if err := svc.Create(context.Background(), &models.WriteEmployeeDTO{Name: "Ivan", DepartmentId: "dept-1"}); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestEmployeeDelete_BlockedWhenInstrumentsUsed(t *testing.T) {
	svc, repo, location := newEmployeeFixture()

	location.getUsedByHolderFn = func(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error) {
		return []*models.Location{{Id: "loc-1"}}, nil
	}

	if err := svc.Delete(context.Background(), "emp-1"); !errors.Is(err, models.ErrDeleteEmployeeHasInstrument) {
		t.Fatalf("expected ErrDeleteEmployeeHasInstrument, got %v", err)
	}
	if len(location.setPersonCalls) != 0 {
		t.Fatalf("expected no SetPerson, got %+v", location.setPersonCalls)
	}
	if len(repo.deleted) != 0 {
		t.Fatalf("expected no repo.Delete, got %+v", repo.deleted)
	}
}

func TestEmployeeDelete_OrderedSetPersonThenDelete(t *testing.T) {
	svc, repo, location := newEmployeeFixture()

	if err := svc.Delete(context.Background(), "emp-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(location.setPersonCalls) != 1 || location.setPersonCalls[0] != "emp-1" {
		t.Fatalf("expected SetPerson with emp-1, got %+v", location.setPersonCalls)
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "emp-1" {
		t.Fatalf("expected repo.Delete with emp-1, got %+v", repo.deleted)
	}
}

func TestEmployeeDelete_PropagatesLocationError(t *testing.T) {
	svc, repo, location := newEmployeeFixture()

	location.getUsedByHolderFn = func(ctx context.Context, dto *models.GetLocationByHolderDTO) ([]*models.Location, error) {
		return nil, models.ErrForbidden
	}

	if err := svc.Delete(context.Background(), "emp-1"); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
	if len(repo.deleted) != 0 {
		t.Fatalf("expected no repo.Delete, got %+v", repo.deleted)
	}
}

func TestEmployeeGetByMostId_PassesThrough(t *testing.T) {
	svc, repo, _ := newEmployeeFixture()
	repo.getByMostFn = func(ctx context.Context, mostId string) (*models.EmployeeData, error) {
		return nil, models.ErrNoRows
	}

	if _, err := svc.GetByMostId(context.Background(), "most-1"); !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}
