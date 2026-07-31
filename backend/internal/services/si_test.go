package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

func newSIServiceFixture() (*SIService, *fakeInstrumentSvc, *fakeVerificationSvc, *fakeLocationSvc) {
	instrument := &fakeInstrumentSvc{}
	verification := &fakeVerificationSvc{}
	location := &fakeLocationSvc{}
	svc := NewSiService(&SiDeps{
		Repo:         &fakeSIRepo{},
		TxManager:    &fakeTxManager{},
		Instrument:   instrument,
		Verification: verification,
		Location:     location,
	})
	return svc, instrument, verification, location
}

func testSiDTO() *models.SiDTO {
	return &models.SiDTO{
		Instrument: &models.InstrumentDTO{
			Id:      "ins-1",
			Name:    "SI-1",
			UserId:  "user-1",
			Actor:   &models.Actor{ID: "user-1", Name: "User"},
		},
	}
}

func TestSIGetById_ReturnsFullBaseSI(t *testing.T) {
	svc, instrument, verification, location := newSIServiceFixture()

	instrument.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return &models.Instrument{Id: "ins-1", Name: "SI-1"}, nil
	}
	verification.getLastFn = func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
		return &models.Verification{Id: "ver-1", InstrumentId: "ins-1"}, nil
	}
	location.getLastFn = func(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
		return &models.Location{Id: "loc-1", InstrumentId: "ins-1", Status: "used"}, nil
	}

	data, err := svc.GetById(context.Background(), &models.GetSiByIdDTO{Id: "ins-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Instrument == nil || data.Instrument.Id != "ins-1" {
		t.Fatalf("unexpected instrument: %+v", data.Instrument)
	}
	if data.Verification == nil || data.Verification.Id != "ver-1" {
		t.Fatalf("unexpected verification: %+v", data.Verification)
	}
	if data.Location == nil || data.Location.Id != "loc-1" {
		t.Fatalf("unexpected location: %+v", data.Location)
	}
}

func TestSIGetById_ToleratesMissingVerification(t *testing.T) {
	svc, instrument, _, location := newSIServiceFixture()

	instrument.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return &models.Instrument{Id: "ins-1"}, nil
	}
	location.getLastFn = func(ctx context.Context, dto *models.GetLocationDTO) (*models.Location, error) {
		return &models.Location{Id: "loc-1"}, nil
	}

	data, err := svc.GetById(context.Background(), &models.GetSiByIdDTO{Id: "ins-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Verification != nil {
		t.Fatalf("expected nil verification, got %+v", data.Verification)
	}
}

func TestSIGetById_ToleratesMissingLocation(t *testing.T) {
	svc, instrument, verification, _ := newSIServiceFixture()

	instrument.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return &models.Instrument{Id: "ins-1"}, nil
	}
	verification.getLastFn = func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
		return &models.Verification{Id: "ver-1"}, nil
	}

	data, err := svc.GetById(context.Background(), &models.GetSiByIdDTO{Id: "ins-1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data.Location != nil {
		t.Fatalf("expected nil location, got %+v", data.Location)
	}
}

func TestSIGetById_PropagatesVerificationError(t *testing.T) {
	svc, instrument, verification, _ := newSIServiceFixture()

	instrument.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return &models.Instrument{Id: "ins-1"}, nil
	}
	verification.getLastFn = func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) {
		return nil, models.ErrForbidden
	}

	if _, err := svc.GetById(context.Background(), &models.GetSiByIdDTO{Id: "ins-1"}); !errors.Is(err, models.ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestSIGetById_PropagatesInstrumentError(t *testing.T) {
	svc, instrument, _, _ := newSIServiceFixture()

	instrument.getByIdFn = func(ctx context.Context, req *models.GetInstrumentByIdDTO) (*models.Instrument, error) {
		return nil, models.ErrNoRows
	}

	if _, err := svc.GetById(context.Background(), &models.GetSiByIdDTO{Id: "ins-1"}); !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows, got %v", err)
	}
}

func TestSICreate_WithVerificationAndLocation(t *testing.T) {
	svc, instrument, verification, location := newSIServiceFixture()

	dto := testSiDTO()
	dto.Verification = &models.VerificationDTO{Id: "ver-1"}
	dto.Location = &models.LocationDTO{Id: "loc-1"}

	if err := svc.Create(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instrument.created) != 1 {
		t.Fatalf("expected instrument.Create, got %d", len(instrument.created))
	}
	if len(verification.created) != 1 {
		t.Fatalf("expected verification.Create, got %d", len(verification.created))
	}
	v := verification.created[0]
	if v.InstrumentId != dto.Instrument.Id || v.UserId != dto.Instrument.UserId || v.Status != string(models.InstrumentStatusWork) {
		t.Fatalf("unexpected verification injection: %+v", v)
	}
	if len(location.locationsCreated) != 1 {
		t.Fatalf("expected location.Create, got %d", len(location.locationsCreated))
	}
	if location.locationsCreated[0].InstrumentId != dto.Instrument.Id {
		t.Fatalf("unexpected location injection: %+v", location.locationsCreated[0])
	}
}

func TestSICreate_WithoutVerificationAndLocation(t *testing.T) {
	svc, instrument, verification, location := newSIServiceFixture()

	if err := svc.Create(context.Background(), testSiDTO()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instrument.created) != 1 {
		t.Fatalf("expected instrument.Create, got %d", len(instrument.created))
	}
	if len(verification.created) != 0 {
		t.Fatalf("expected no verification.Create, got %d", len(verification.created))
	}
	if len(location.locationsCreated) != 0 {
		t.Fatalf("expected no location.Create, got %d", len(location.locationsCreated))
	}
}

func TestSICreate_PropagatesError(t *testing.T) {
	svc, instrument, _, _ := newSIServiceFixture()

	instrument.createFn = func(ctx context.Context, tx postgres.Tx, dto *models.InstrumentDTO) error {
		return models.ErrNotValid
	}

	if err := svc.Create(context.Background(), testSiDTO()); !errors.Is(err, models.ErrNotValid) {
		t.Fatalf("expected ErrNotValid, got %v", err)
	}
}

func TestSIUpdate_WithEmptyVerificationId_CreatesVerification(t *testing.T) {
	svc, instrument, verification, _ := newSIServiceFixture()

	dto := testSiDTO()
	dto.Verification = &models.VerificationDTO{}

	if err := svc.Update(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instrument.updated) != 1 {
		t.Fatalf("expected instrument.Update, got %d", len(instrument.updated))
	}
	if len(verification.created) != 1 {
		t.Fatalf("expected verification.Create, got %d", len(verification.created))
	}
	if len(verification.updated) != 0 {
		t.Fatalf("expected no verification.Update, got %d", len(verification.updated))
	}
}

func TestSIUpdate_WithExistingVerificationId_UpdatesVerification(t *testing.T) {
	svc, instrument, verification, _ := newSIServiceFixture()

	dto := testSiDTO()
	dto.Verification = &models.VerificationDTO{Id: "ver-1"}

	if err := svc.Update(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instrument.updated) != 1 {
		t.Fatalf("expected instrument.Update, got %d", len(instrument.updated))
	}
	if len(verification.created) != 0 {
		t.Fatalf("expected no verification.Create, got %d", len(verification.created))
	}
	if len(verification.updated) != 1 {
		t.Fatalf("expected verification.Update, got %d", len(verification.updated))
	}
}

func TestSIUpdate_WithoutVerification(t *testing.T) {
	svc, instrument, verification, _ := newSIServiceFixture()

	if err := svc.Update(context.Background(), testSiDTO()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instrument.updated) != 1 {
		t.Fatalf("expected instrument.Update, got %d", len(instrument.updated))
	}
	if len(verification.created) != 0 || len(verification.updated) != 0 {
		t.Fatalf("expected no verification calls, got created=%d updated=%d", len(verification.created), len(verification.updated))
	}
}

func TestSIDelete_MapsErrNoRows(t *testing.T) {
	svc, instrument, _, _ := newSIServiceFixture()

	instrument.deleteFn = func(ctx context.Context, dto *models.DeleteSiDTO) error {
		return models.ErrNoRows
	}

	if err := svc.Delete(context.Background(), &models.DeleteSiDTO{Id: "ins-1"}); !errors.Is(err, models.ErrDeleteInstrumentAtHolder) {
		t.Fatalf("expected ErrDeleteInstrumentAtHolder, got %v", err)
	}
}

func TestSIDelete_Success(t *testing.T) {
	svc, instrument, _, _ := newSIServiceFixture()

	if err := svc.Delete(context.Background(), &models.DeleteSiDTO{Id: "ins-1", Actor: &models.Actor{ID: "user-1"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instrument.deleted) != 1 || instrument.deleted[0].Id != "ins-1" {
		t.Fatalf("expected instrument.Delete with ins-1, got %+v", instrument.deleted)
	}
}

func TestSIGet_PassesThrough(t *testing.T) {
	svc, _, _, _ := newSIServiceFixture()
	svc.repo = &fakeSIRepo{getFn: func(ctx context.Context, req *models.GetSiDTO) ([]*models.SI, error) {
		return []*models.SI{{Id: "si-1"}}, nil
	}}

	data, err := svc.Get(context.Background(), &models.GetSiDTO{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 || data[0].Id != "si-1" {
		t.Fatalf("unexpected data: %+v", data)
	}
}
