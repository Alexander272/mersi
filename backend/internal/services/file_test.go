package services

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/xuri/excelize/v2"
)

func openXlsx(t *testing.T, data []byte) *excelize.File {
	t.Helper()
	f, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("failed to open xlsx: %v", err)
	}
	return f
}

func newFileFixture() *FileService {
	return NewFileService()
}

func TestResolveReportStrategy_SimpleGrid(t *testing.T) {
	svc := newFileFixture()

	for _, bidType := range []string{"met_si", "eq_si"} {
		strategy := svc.resolveReportStrategy(bidType)
		if len(strategy.headers) != 9 {
			t.Fatalf("bidType %s: expected 9 headers, got %d", bidType, len(strategy.headers))
		}
		row := strategy.mapFunc(0, &models.SI{Name: "SI-1"})
		if len(row) != 9 {
			t.Fatalf("bidType %s: expected 9 columns, got %d", bidType, len(row))
		}
		if row[0] != 1 {
			t.Fatalf("expected first column to be index 1, got %v", row[0])
		}
	}
}

func TestResolveReportStrategy_AttestationGrid(t *testing.T) {
	svc := newFileFixture()

	strategy := svc.resolveReportStrategy("some_eq_bid")
	if len(strategy.headers) != 20 {
		t.Fatalf("expected 20 headers, got %d", len(strategy.headers))
	}
	if strategy.headers[7] != "I" || strategy.headers[18] != "XII" {
		t.Fatalf("unexpected month headers: %+v", strategy.headers[7:19])
	}
}

func TestResolveReportStrategy_DefaultGrid(t *testing.T) {
	svc := newFileFixture()

	strategy := svc.resolveReportStrategy("verification")
	if len(strategy.headers) != 20 {
		t.Fatalf("expected 20 headers, got %d", len(strategy.headers))
	}
}

func TestMapWithMonthlyGrid_PlacesStarAtNextDateMonth(t *testing.T) {
	svc := newFileFixture()

	d := &models.SI{
		Name:              "SI-1",
		NextVerificationDate: time.Date(2026, time.March, 15, 0, 0, 0, 0, time.UTC),
	}
	row := svc.mapWithMonthlyGrid(0, d, constants.ShortDateFormat)

	if len(row) != 20 {
		t.Fatalf("expected 20 columns, got %d", len(row))
	}
	for i := 7; i < 19; i++ {
		expected := ""
		if i == 7+2 {
			expected = "*"
		}
		if row[i] != expected {
			t.Fatalf("column %d: expected %q, got %v", i, expected, row[i])
		}
	}
}

func TestMapWithMonthlyGrid_NotesAppended(t *testing.T) {
	svc := newFileFixture()

	d := &models.SI{Notes: "note"}
	row := svc.mapWithMonthlyGrid(0, d, constants.ShortDateFormat)
	if row[19] != "note" {
		t.Fatalf("expected note as last column, got %v", row[19])
	}
}

func TestMapReflectRow_FormatsDates(t *testing.T) {
	svc := newFileFixture()

	si := &models.SI{
		VerificationDate: time.Date(2026, time.January, 5, 0, 0, 0, 0, time.UTC),
		NextVerificationDate: time.Date(2027, time.February, 6, 0, 0, 0, 0, time.UTC),
	}
	cols := []*models.Column{
		{Name: "Дата", Field: "verificationDate", Type: "date"},
		{Name: "Дата кратко", Field: "nextVerificationDate", Type: "short_date"},
	}
	row, err := svc.mapReflectRow(si, cols)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if row[0] != "05.01.2026" {
		t.Fatalf("expected full date, got %v", row[0])
	}
	if row[1] != "02.2027" {
		t.Fatalf("expected short date, got %v", row[1])
	}
}

func TestMapReflectRow_ZeroTimeBecomesEmptyString(t *testing.T) {
	svc := newFileFixture()

	si := &models.SI{}
	cols := []*models.Column{
		{Name: "Дата", Field: "verificationDate", Type: "date"},
	}
	row, err := svc.mapReflectRow(si, cols)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if row[0] != "" {
		t.Fatalf("expected empty string for zero time, got %v", row[0])
	}
}

func TestMapReflectRow_ReturnsStringField(t *testing.T) {
	svc := newFileFixture()

	si := &models.SI{Name: "SI-1"}
	cols := []*models.Column{{Name: "Название", Field: "name", Type: "string"}}
	row, err := svc.mapReflectRow(si, cols)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if row[0] != "SI-1" {
		t.Fatalf("expected SI-1, got %v", row[0])
	}
}

func TestMapReflectRow_UnknownFieldReturnsError(t *testing.T) {
	svc := newFileFixture()

	cols := []*models.Column{{Name: "X", Field: "unknownField", Type: "string"}}
	if _, err := svc.mapReflectRow(&models.SI{}, cols); err == nil {
		t.Fatal("expected error for unknown field")
	}
}

func TestMakeAccountingLog_Pre1900DateUsesDash(t *testing.T) {
	svc := newFileFixture()

	buf, err := svc.MakeAccountingLog(context.Background(), []*models.SiWithLog{
		{Name: "SI-1", DateOfReceipt: time.Date(1899, time.December, 1, 0, 0, 0, 0, time.Local)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty buffer")
	}

	f := openXlsx(t, buf.Bytes())
	defer f.Close()
	val, err := f.GetCellValue("Sheet1", "A2")
	if err != nil {
		t.Fatalf("failed to read cell: %v", err)
	}
	if val != "–" {
		t.Fatalf("expected dash for pre-1900 date, got %q", val)
	}
}

func TestMakeAccountingLog_ModernDateFormatted(t *testing.T) {
	svc := newFileFixture()

	buf, err := svc.MakeAccountingLog(context.Background(), []*models.SiWithLog{
		{Name: "SI-1", DateOfReceipt: time.Date(2026, time.March, 10, 0, 0, 0, 0, time.Local)},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty buffer")
	}
}

func TestExport_FlattensChildColumns(t *testing.T) {
	svc := newFileFixture()

	dto := &models.ExportDTO{
		Columns: []*models.Column{
			{Name: "Parent", Field: "id", Type: "string", Children: []*models.Column{
				{Name: "Child1", Field: "name", Type: "string"},
				{Name: "Child2", Field: "notes", Type: "string"},
			}},
		},
		SI: []*models.SI{{Id: "si-1", Name: "SI-1", Notes: "note"}},
	}

	buf, err := svc.Export(context.Background(), dto)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty buffer")
	}
}

func TestMakeVerificationSchedule_GeneratesFile(t *testing.T) {
	svc := newFileFixture()

	for _, bidType := range []string{"met_si", "some_eq_bid", "verification"} {
		buf, err := svc.MakeVerificationSchedule(context.Background(), &models.SiVerification{
			BidType: bidType,
			SI:      []*models.SI{{Name: "SI-1"}},
		})
		if err != nil {
			t.Fatalf("bidType %s: unexpected error: %v", bidType, err)
		}
		if buf.Len() == 0 {
			t.Fatalf("bidType %s: expected non-empty buffer", bidType)
		}
	}
}
