package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository/postgres"
)

func newVerificationService(repo *fakeVerificationRepo, verDoc *fakeVerDocSvc, instrument *fakeInstrumentSvc, docs *fakeDocumentSvc) *VerificationService {
	return NewVerificationService(&VerificationDeps{
		Repo:        repo,
		TxManager:   &fakeTxManager{},
		VerDocs:     verDoc,
		Instrument:  instrument,
		Docs:        docs,
		ActivityLog: &fakeActivityLogSvc{},
	})
}

func verificationDTO() *models.VerificationDTO {
	return &models.VerificationDTO{
		Id:           "v1",
		InstrumentId: "i1",
		Date:         time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
		NextDate:     time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
		Status:       string(models.InstrumentStatusChecking),
		Actor:        &models.Actor{ID: "u1", Name: "User"},
	}
}

func TestVerificationGet_AttachesGroupedDocs(t *testing.T) {
	repo := &fakeVerificationRepo{
		getFn: func(ctx context.Context, req *models.GetVerificationDTO) ([]*models.Verification, error) {
			return []*models.Verification{{Id: "v1", InstrumentId: "i1"}}, nil
		},
	}
	verDoc := &fakeVerDocSvc{
		getGroupedFn: func(ctx context.Context, req *models.GetGroupedVerificationDocsDTO) (*models.GroupedVerificationDocs, error) {
			if req.InstrumentId != "i1" {
				t.Fatalf("expected instrument i1, got %s", req.InstrumentId)
			}
			return &models.GroupedVerificationDocs{
				Groups: map[string]*models.Groups{
					"v1": {Docs: []*models.VerificationDoc{{Id: "d1", Name: "cert"}}},
				},
			}, nil
		},
	}
	svc := newVerificationService(repo, verDoc, &fakeInstrumentSvc{}, &fakeDocumentSvc{})

	data, err := svc.Get(context.Background(), &models.GetVerificationDTO{InstrumentId: "i1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 || len(data[0].Docs) != 1 || data[0].Docs[0].Name != "cert" {
		t.Fatalf("expected docs attached to verification, got %+v", data[0])
	}
}

func TestVerificationCreate_AlreadyExists(t *testing.T) {
	repo := &fakeVerificationRepo{
		getByInstrumentAndDateFn: func(ctx context.Context, tx postgres.Tx, instrumentId string, date time.Time) (*models.Verification, error) {
			return &models.Verification{Id: "existing"}, nil
		},
	}
	svc := newVerificationService(repo, &fakeVerDocSvc{}, &fakeInstrumentSvc{}, &fakeDocumentSvc{})

	err := svc.Create(context.Background(), &fakeTx{}, verificationDTO())
	if !errors.Is(err, models.ErrVerificationAlreadyExists) {
		t.Fatalf("expected ErrVerificationAlreadyExists, got %v", err)
	}
}

func TestVerificationCreate_Success(t *testing.T) {
	repo := &fakeVerificationRepo{}
	verDoc := &fakeVerDocSvc{}
	instrument := &fakeInstrumentSvc{}
	activity := &fakeActivityLogSvc{}
	svc := NewVerificationService(&VerificationDeps{
		Repo:        repo,
		TxManager:   &fakeTxManager{},
		VerDocs:     verDoc,
		Instrument:  instrument,
		Docs:        &fakeDocumentSvc{},
		ActivityLog: activity,
	})

	dto := verificationDTO()
	dto.Docs = []*models.VerificationDocDTO{{Id: "doc1"}, {Id: "doc2"}}
	if err := svc.Create(context.Background(), &fakeTx{}, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(verDoc.created) != 2 {
		t.Fatalf("expected 2 docs created, got %d", len(verDoc.created))
	}
	for _, d := range verDoc.created {
		if d.VerificationId != "v1" {
			t.Fatalf("expected doc verification_id=v1, got %s", d.VerificationId)
		}
	}
	if len(instrument.statusChanges) != 1 || instrument.statusChanges[0].Id != "i1" {
		t.Fatalf("expected instrument status change, got %+v", instrument.statusChanges)
	}
	if len(activity.logs) != 1 || activity.logs[0].Action != "CREATE" {
		t.Fatalf("expected CREATE activity log, got %+v", activity.logs)
	}
}

func TestVerificationCreate_SkipsDocsWhenEmpty(t *testing.T) {
	verDoc := &fakeVerDocSvc{}
	svc := newVerificationService(&fakeVerificationRepo{}, verDoc, &fakeInstrumentSvc{}, &fakeDocumentSvc{})
	if err := svc.Create(context.Background(), &fakeTx{}, verificationDTO()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(verDoc.created) != 0 {
		t.Fatalf("expected no docs created, got %d", len(verDoc.created))
	}
}

func TestVerificationCreateSeveral_Empty(t *testing.T) {
	svc := newVerificationService(&fakeVerificationRepo{}, &fakeVerDocSvc{}, &fakeInstrumentSvc{}, &fakeDocumentSvc{})
	if err := svc.CreateSeveral(context.Background(), nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestVerificationCreateSeveral_Success(t *testing.T) {
	repo := &fakeVerificationRepo{}
	verDoc := &fakeVerDocSvc{}
	svc := newVerificationService(repo, verDoc, &fakeInstrumentSvc{}, &fakeDocumentSvc{})

	dto := []*models.VerificationDTO{
		{Id: "v1", InstrumentId: "i1", Actor: &models.Actor{ID: "u1", Name: "U"}, Docs: []*models.VerificationDocDTO{{Id: "d1"}}},
		{Id: "v2", InstrumentId: "i2", Actor: &models.Actor{ID: "u1", Name: "U"}, Docs: []*models.VerificationDocDTO{{Id: "d2"}, {Id: "d3"}}},
	}
	if err := svc.CreateSeveral(context.Background(), dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(verDoc.created) != 3 {
		t.Fatalf("expected 3 docs created, got %d", len(verDoc.created))
	}
	byID := map[string]string{}
	for _, d := range verDoc.created {
		byID[d.Id] = d.VerificationId
	}
	if byID["d1"] != "v1" || byID["d2"] != "v2" || byID["d3"] != "v2" {
		t.Fatalf("expected docs linked to their verification, got %v", byID)
	}
}

func TestVerificationUpdate_SplitsNewAndUpdatedDocs(t *testing.T) {
	verDoc := &fakeVerDocSvc{}
	instrument := &fakeInstrumentSvc{}
	activity := &fakeActivityLogSvc{}
	svc := NewVerificationService(&VerificationDeps{
		Repo:        &fakeVerificationRepo{getByIdFn: func(ctx context.Context, id string) (*models.Verification, error) { return &models.Verification{Id: "v1"}, nil }},
		TxManager:   &fakeTxManager{},
		VerDocs:     verDoc,
		Instrument:  instrument,
		Docs:        &fakeDocumentSvc{},
		ActivityLog: activity,
	})

	dto := verificationDTO()
	dto.Docs = []*models.VerificationDocDTO{
		{Id: "", Name: "new"},
		{Id: "existing", Name: "old"},
	}
	if err := svc.Update(context.Background(), &fakeTx{}, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(verDoc.created) != 1 || verDoc.created[0].Name != "new" || verDoc.created[0].VerificationId != "v1" {
		t.Fatalf("expected new doc created and linked, got %+v", verDoc.created)
	}
	if len(verDoc.updated) != 1 || verDoc.updated[0].Id != "existing" {
		t.Fatalf("expected existing doc updated, got %+v", verDoc.updated)
	}
	if len(instrument.statusChanges) != 1 {
		t.Fatalf("expected instrument status change, got %+v", instrument.statusChanges)
	}
}

func TestVerificationUpdate_DeletesDocs(t *testing.T) {
	docs := &fakeDocumentSvc{}
	verDoc := &fakeVerDocSvc{}
	svc := newVerificationService(&fakeVerificationRepo{}, verDoc, &fakeInstrumentSvc{}, docs)

	dto := verificationDTO()
	dto.DeletedDocs = []models.DeletedDoc{{DocId: "doc1", Filename: "f1.pdf"}}
	if err := svc.Update(context.Background(), &fakeTx{}, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(docs.deletedDocs) != 1 {
		t.Fatalf("expected 1 document deleted, got %d", len(docs.deletedDocs))
	}
	if docs.deletedDocs[0].Id != "doc1" || docs.deletedDocs[0].Group != "verification" {
		t.Fatalf("unexpected delete dto: %+v", docs.deletedDocs[0])
	}
	if len(verDoc.deleted) != 1 || verDoc.deleted[0] != "doc1" {
		t.Fatalf("expected doc link deleted, got %v", verDoc.deleted)
	}
}

func TestDelete_LastChangesStatusToWork(t *testing.T) {
	old := &models.Verification{Id: "v1", InstrumentId: "i1", Date: time.Now(), NextDate: time.Now()}
	instrument := &fakeInstrumentSvc{}
	svc := newVerificationService(&fakeVerificationRepo{
		getByIdFn: func(ctx context.Context, id string) (*models.Verification, error) { return old, nil },
		getLastFn: func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) { return old, nil },
	}, &fakeVerDocSvc{}, instrument, &fakeDocumentSvc{})

	if err := svc.Delete(context.Background(), &models.DeleteVerificationDTO{Id: "v1", Actor: &models.Actor{ID: "u1", Name: "U"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instrument.statusChanges) != 1 {
		t.Fatalf("expected status change to work, got %+v", instrument.statusChanges)
	}
	if instrument.statusChanges[0].Status != models.InstrumentStatusWork {
		t.Fatalf("expected status work, got %s", instrument.statusChanges[0].Status)
	}
}

func TestDelete_NotLastKeepsStatus(t *testing.T) {
	old := &models.Verification{Id: "v1", InstrumentId: "i1"}
	other := &models.Verification{Id: "v2", InstrumentId: "i1"}
	instrument := &fakeInstrumentSvc{}
	svc := newVerificationService(&fakeVerificationRepo{
		getByIdFn: func(ctx context.Context, id string) (*models.Verification, error) { return old, nil },
		getLastFn: func(ctx context.Context, req *models.GetVerificationDTO) (*models.Verification, error) { return other, nil },
	}, &fakeVerDocSvc{}, instrument, &fakeDocumentSvc{})

	if err := svc.Delete(context.Background(), &models.DeleteVerificationDTO{Id: "v1", Actor: &models.Actor{ID: "u1", Name: "U"}}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instrument.statusChanges) != 0 {
		t.Fatalf("expected no status change, got %+v", instrument.statusChanges)
	}
}
