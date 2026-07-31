//go:build integration

package postgres

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/migrate"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupIntegrationDB(t *testing.T) *sqlx.DB {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(90 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = container.Terminate(ctx) })

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=postgres password=postgres dbname=test sslmode=disable", host, port.Port())
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to connect to postgres: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if err := migrate.Migrate(db.DB); err != nil {
		t.Fatalf("failed to run migrations: %v", err)
	}

	return db
}

type seededInstrument struct {
	sectionId    string
	instrumentId string
}

func seedInstrument(t *testing.T, db *sqlx.DB) *seededInstrument {
	t.Helper()
	s := &seededInstrument{sectionId: uuid.NewString(), instrumentId: uuid.NewString()}

	if _, err := db.Exec(`INSERT INTO sections(id, realm_id, name, position) VALUES($1, $2, $3, $4)`,
		s.sectionId, uuid.NewString(), "sec", 1); err != nil {
		t.Fatalf("failed to insert section: %v", err)
	}

	if _, err := db.Exec(`INSERT INTO instruments(id, section_id, position, name, date_of_receipt, factory_number, status, act_of_entering_id)
		VALUES($1, $2, $3, 'SI-1', now(), '123', 'work', $4)`,
		s.instrumentId, s.sectionId, 1, uuid.Nil.String()); err != nil {
		t.Fatalf("failed to insert instrument: %v", err)
	}
	return s
}

func seedSentinelLocation(t *testing.T, db *sqlx.DB, instrumentId string, issueDaysAgo int) string {
	t.Helper()
	id := uuid.NewString()
	if _, err := db.Exec(`INSERT INTO locations(id, instrument_id, date_of_receiving, date_of_issue, status, need_confirmed, has_confirmed, place)
		VALUES($1, $2, '0001-01-01'::date, now() - $3 * interval '1 day', 'moved', true, false, 'dep')`,
		id, instrumentId, issueDaysAgo); err != nil {
		t.Fatalf("failed to insert location: %v", err)
	}
	return id
}

type locationRow struct {
	Status          string    `db:"status"`
	DateOfReceiving time.Time `db:"date_of_receiving"`
	HasConfirmed    bool      `db:"has_confirmed"`
}

func getLocation(t *testing.T, db *sqlx.DB, id string) locationRow {
	t.Helper()
	var row locationRow
	if err := db.Get(&row, `SELECT status, date_of_receiving, has_confirmed FROM locations WHERE id=$1`, id); err != nil {
		t.Fatalf("failed to read location: %v", err)
	}
	return row
}

func TestLocationReceiving_UpdatesLocation(t *testing.T) {
	db := setupIntegrationDB(t)
	seed := seedInstrument(t, db)
	locId := seedSentinelLocation(t, db, seed.instrumentId, 1)

	repo := NewLocationRepo(db, NewTransactionRepo(db))
	err := repo.Receiving(context.Background(), &models.ReceivingDTO{
		InstrumentIds: []string{seed.instrumentId},
		Status:        constants.LocationStatusUsed,
		HasConfirmed:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	row := getLocation(t, db, locId)
	if row.Status != constants.LocationStatusUsed {
		t.Fatalf("expected status used, got %s", row.Status)
	}
	if !row.HasConfirmed {
		t.Fatal("expected has_confirmed=true")
	}
	if row.DateOfReceiving.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected real receipt date, got sentinel %v", row.DateOfReceiving)
	}
}

func TestLocationReceiving_AlreadyReceivedReturnsNoRows(t *testing.T) {
	db := setupIntegrationDB(t)
	seed := seedInstrument(t, db)
	seedSentinelLocation(t, db, seed.instrumentId, 1)

	repo := NewLocationRepo(db, NewTransactionRepo(db))
	ctx := context.Background()
	dto := &models.ReceivingDTO{InstrumentIds: []string{seed.instrumentId}, Status: constants.LocationStatusUsed}

	if err := repo.Receiving(ctx, dto); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.Receiving(ctx, dto); !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows on second call, got %v", err)
	}
}

func TestForcedReceipt_UpdatesStatusAndDate(t *testing.T) {
	db := setupIntegrationDB(t)
	seed := seedInstrument(t, db)
	locId := seedSentinelLocation(t, db, seed.instrumentId, 1)

	repo := NewLocationRepo(db, NewTransactionRepo(db))
	if err := repo.ForcedReceipt(context.Background(), &models.ForcedReceiptDTO{InstrumentId: seed.instrumentId}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	row := getLocation(t, db, locId)
	if row.Status != constants.LocationStatusUsed {
		t.Fatalf("expected status used, got %s", row.Status)
	}
	if row.DateOfReceiving.Before(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected real receipt date, got %v", row.DateOfReceiving)
	}

	if err := repo.ForcedReceipt(context.Background(), &models.ForcedReceiptDTO{InstrumentId: seed.instrumentId}); !errors.Is(err, models.ErrNoRows) {
		t.Fatalf("expected ErrNoRows on second call, got %v", err)
	}
}

func TestForcedReceiptAll_SwitchesToReserveBasedOnPreviousLocation(t *testing.T) {
	db := setupIntegrationDB(t)
	seed := seedInstrument(t, db)

	// прошлое перемещение уже принято (used)
	if _, err := db.Exec(`INSERT INTO locations(id, instrument_id, date_of_receiving, date_of_issue, status)
		VALUES($1, $2, now(), now() - 40 * interval '1 day', 'used')`, uuid.NewString(), seed.instrumentId); err != nil {
		t.Fatalf("failed to insert previous location: %v", err)
	}
	// текущее перемещение: просрочено и не принято
	locId := seedSentinelLocation(t, db, seed.instrumentId, 30)

	repo := NewLocationRepo(db, NewTransactionRepo(db))
	if err := repo.ForcedReceiptAll(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	row := getLocation(t, db, locId)
	if row.Status != constants.LocationStatusReserve {
		t.Fatalf("expected status reserve, got %s", row.Status)
	}
}

func TestForcedReceiptAll_IgnoresRecentLocations(t *testing.T) {
	db := setupIntegrationDB(t)
	seed := seedInstrument(t, db)
	locId := seedSentinelLocation(t, db, seed.instrumentId, 1)

	repo := NewLocationRepo(db, NewTransactionRepo(db))
	if err := repo.ForcedReceiptAll(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	row := getLocation(t, db, locId)
	if row.Status != constants.LocationStatusMoved {
		t.Fatalf("expected location untouched, got status %s", row.Status)
	}
}

func TestSIGet_InstrumentWithoutActOfEntering(t *testing.T) {
	db := setupIntegrationDB(t)
	seed := seedInstrument(t, db)

	siRepo := NewSIRepo(db, NewTransactionRepo(db))
	data, err := siRepo.Get(context.Background(), &models.GetSiDTO{
		SectionId: seed.sectionId,
		Status:    models.InstrumentStatusWork,
		Page:      &models.Page{Limit: 15},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 1 {
		t.Fatalf("expected 1 instrument, got %d", len(data))
	}
	if data[0].ActOfEnteringId != uuid.Nil.String() {
		t.Fatalf("expected nil-uuid act_of_entering_id, got %q", data[0].ActOfEnteringId)
	}
}
