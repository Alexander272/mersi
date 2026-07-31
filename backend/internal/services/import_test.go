package services

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

func newImportFixture() (*ImportService, *fakeInstrumentSvc, *fakeVerificationSvc, *fakeRepairSvc, *fakePreservationSvc, *fakeTransferToDepSvc, *fakeWriteOffSvc) {
	instrument := &fakeInstrumentSvc{}
	verification := &fakeVerificationSvc{}
	repair := &fakeRepairSvc{}
	preservation := &fakePreservationSvc{}
	transferToDep := &fakeTransferToDepSvc{}
	writeOff := &fakeWriteOffSvc{}

	svc := NewImportService(&ImportDeps{
		Instrument:     instrument,
		Verification:   verification,
		Repair:         repair,
		Preservation:   preservation,
		TransferToSave: &fakeTransferToSaveSvc{},
		TransferToDep:  transferToDep,
		WriteOff:       writeOff,
	})
	return svc, instrument, verification, repair, preservation, transferToDep, writeOff
}

// buildImportXlsx создаёт xlsx с листом "si" из заданных строк.
func buildImportXlsx(t *testing.T, rows [][]interface{}) []byte {
	t.Helper()
	file := excelize.NewFile()
	defer file.Close()

	index, err := file.NewSheet("si")
	if err != nil {
		t.Fatalf("failed to create sheet: %v", err)
	}
	file.SetActiveSheet(index)

	for i, row := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+1)
		if err := file.SetSheetRow("si", cell, &row); err != nil {
			t.Fatalf("failed to set row %d: %v", i, err)
		}
	}

	buf, err := file.WriteToBuffer()
	if err != nil {
		t.Fatalf("failed to write xlsx: %v", err)
	}
	return buf.Bytes()
}

// buildImportFile оборачивает xlsx в multipart.FileHeader, как его отдаёт клиент.
func buildImportFile(t *testing.T, xlsx []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "test.xlsx")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write(xlsx); err != nil {
		t.Fatalf("failed to write form file: %v", err)
	}
	writer.Close()

	req := httptest.NewRequest("POST", "/", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		t.Fatalf("failed to parse multipart form: %v", err)
	}
	files := req.MultipartForm.File["file"]
	if len(files) != 1 {
		t.Fatalf("expected 1 file part, got %d", len(files))
	}
	return files[0]
}

func importDTO(file *multipart.FileHeader) *models.ImportDTO {
	return &models.ImportDTO{
		SectionId: "sec-1",
		BidType:   "ointo_si",
		UserId:    "user-1",
		File:      file,
	}
}

func plainRow(name string, extra ...int) []interface{} {
	row := make([]interface{}, 28)
	for i := range row {
		row[i] = ""
	}
	row[0] = "01.02.2026"
	row[1] = name
	row[2] = "Термометр"
	row[3] = "12345"
	row[4] = "Иванов"
	row[9] = "2000"
	row[10] = "Россия"
	row[11] = "ООО Тест"
	row[12] = "REG-1"
	row[13] = "0-100"
	row[14] = "12"
	for _, i := range extra {
		row[i] = ""
	}
	return row
}

func TestLoadOintoSi_ParsesInstrumentFields(t *testing.T) {
	svc, instrument, _, _, _, _, _ := newImportFixture()

	xlsx := buildImportXlsx(t, [][]interface{}{plainRow("SI-1")})
	err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(instrument.createdSeveral) != 1 || len(instrument.createdSeveral[0]) != 1 {
		t.Fatalf("expected 1 instrument created, got %+v", instrument.createdSeveral)
	}
	si := instrument.createdSeveral[0][0]
	if si.SectionId != "sec-1" || si.UserId != "user-1" {
		t.Fatalf("unexpected section/user: %+v", si)
	}
	if si.Name != "SI-1" || si.Type != "Термометр" || si.FactoryNumber != "12345" || si.Responsible != "Иванов" {
		t.Fatalf("unexpected fields: %+v", si)
	}
	if si.YearOfIssue != 2000 || si.InterVerificationInterval != 12 {
		t.Fatalf("unexpected numeric fields: %+v", si)
	}
	if si.MeasurementLimits != "0-100" || si.StateRegister != "REG-1" || si.CountryOfProduce != "Россия" || si.Manufacturer != "ООО Тест" {
		t.Fatalf("unexpected extra fields: %+v", si)
	}
	if si.DateOfReceipt.Year() != 2026 || si.DateOfReceipt.Month() != time.February || si.DateOfReceipt.Day() != 1 {
		t.Fatalf("unexpected date of receipt: %v", si.DateOfReceipt)
	}
	if si.Status != models.InstrumentStatusWork {
		t.Fatalf("expected work status, got %s", si.Status)
	}
}

func TestLoadOintoSi_RepairWithoutEndDateSetsRepairStatus(t *testing.T) {
	svc, instrument, _, repair, _, _, _ := newImportFixture()

	row := plainRow("SI-1")
	row[5] = "05.05.2025"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if instrument.createdSeveral[0][0].Status != models.InstrumentStatusRepair {
		t.Fatalf("expected repair status, got %s", instrument.createdSeveral[0][0].Status)
	}
	if len(repair.createdSeveral) != 1 || repair.createdSeveral[0][0].PeriodEnd != (time.Time{}) {
		t.Fatalf("unexpected repair record: %+v", repair.createdSeveral)
	}
}

func TestLoadOintoSi_RepairWithEndDateKeepsWorkStatus(t *testing.T) {
	svc, instrument, _, repair, _, _, _ := newImportFixture()

	row := plainRow("SI-1")
	row[5] = "01.01.2025-01.02.2025 (ремонт)"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if instrument.createdSeveral[0][0].Status != models.InstrumentStatusWork {
		t.Fatalf("expected work status, got %s", instrument.createdSeveral[0][0].Status)
	}
	r := repair.createdSeveral[0][0]
	if r.PeriodStart.Year() != 2025 || r.PeriodEnd.Year() != 2025 || r.Work != "ремонт" {
		t.Fatalf("unexpected repair record: %+v", r)
	}
}

func TestLoadOintoSi_PreservationWithoutEndDateSetsArchivedStatus(t *testing.T) {
	svc, instrument, _, _, preservation, _, _ := newImportFixture()

	row := plainRow("SI-1")
	row[6] = "06.06.2025 (консервация)"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if instrument.createdSeveral[0][0].Status != models.InstrumentStatusArchived {
		t.Fatalf("expected archived status, got %s", instrument.createdSeveral[0][0].Status)
	}
	if len(preservation.createdSeveral) != 1 || preservation.createdSeveral[0][0].NotesStart != "консервация" {
		t.Fatalf("unexpected preservation record: %+v", preservation.createdSeveral)
	}
}

func TestLoadOintoSi_WriteOffSetsWriteOffStatus(t *testing.T) {
	svc, instrument, _, _, _, _, writeOff := newImportFixture()

	row := plainRow("SI-1")
	row[8] = "15.03.2025"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if instrument.createdSeveral[0][0].Status != models.InstrumentStatusWriteOff {
		t.Fatalf("expected write_off status, got %s", instrument.createdSeveral[0][0].Status)
	}
	if len(writeOff.createdSeveral) != 1 {
		t.Fatalf("expected write off record, got %+v", writeOff.createdSeveral)
	}
}

func TestLoadOintoSi_TransferSetsTransferredStatus(t *testing.T) {
	svc, instrument, _, _, _, transferToDep, _ := newImportFixture()

	row := plainRow("SI-1")
	row[7] = "20.04.2025"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if instrument.createdSeveral[0][0].Status != models.InstrumentStatusTransferred {
		t.Fatalf("expected transferred status, got %s", instrument.createdSeveral[0][0].Status)
	}
	if len(transferToDep.createdSeveral) != 1 {
		t.Fatalf("expected transfer record, got %+v", transferToDep.createdSeveral)
	}
}

func TestLoadOintoSi_VerificationDatesAndDecommissioning(t *testing.T) {
	svc, _, verification, _, _, _, _ := newImportFixture()

	row := plainRow("SI-1")
	row[15] = "01.02.2027"
	row[16] = "списан 01.02.2028"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(verification.createdSeveral) != 1 || len(verification.createdSeveral[0]) != 2 {
		t.Fatalf("expected 2 verifications, got %+v", verification.createdSeveral)
	}
	v1, v2 := verification.createdSeveral[0][0], verification.createdSeveral[0][1]

	if v1.Date.Year() != 2027 || v1.NextDate.Year() != 2028 {
		t.Fatalf("unexpected dates: %+v", v1)
	}
	if !v1.NextDate.Equal(v1.Date.AddDate(0, 12, 0)) {
		t.Fatalf("expected next date to be interval months ahead, got %+v", v1)
	}
	if v1.Status != "work" || v1.UserId != "user-1" {
		t.Fatalf("unexpected verification: %+v", v1)
	}

	if v2.Status != "decommissioning" {
		t.Fatalf("expected decommissioning for списан, got %+v", v2)
	}
	if len(v2.Docs) != 1 || v2.Docs[0].Name != "списан 01.02.2028" {
		t.Fatalf("unexpected doc: %+v", v2.Docs)
	}
}

func TestLoadOintoSi_RemapsPlaceholderIds(t *testing.T) {
	svc, instrument, verification, repair, _, transferToDep, writeOff := newImportFixture()

	instrument.createSeveralFn = func(ctx context.Context, dto []*models.InstrumentDTO) error {
		for i := range dto {
			dto[i].Id = "real-" + uuid.NewString()
		}
		return nil
	}

	row := plainRow("SI-1")
	row[5] = "01.01.2025-01.02.2025 (ремонт)"
	row[7] = "20.04.2025"
	row[8] = "15.03.2025"
	row[15] = "01.02.2027"
	xlsx := buildImportXlsx(t, [][]interface{}{row})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	realId := instrument.createdSeveral[0][0].Id
	if realId == "" {
		t.Fatal("expected real id to be assigned")
	}
	if verification.createdSeveral[0][0].InstrumentId != realId {
		t.Fatalf("expected verification instrument id remapped, got %q want %q", verification.createdSeveral[0][0].InstrumentId, realId)
	}
	if repair.createdSeveral[0][0].InstrumentId != realId {
		t.Fatalf("expected repair instrument id remapped, got %q", repair.createdSeveral[0][0].InstrumentId)
	}
	if transferToDep.createdSeveral[0][0].InstrumentId != realId {
		t.Fatalf("expected transfer instrument id remapped, got %q", transferToDep.createdSeveral[0][0].InstrumentId)
	}
	if writeOff.createdSeveral[0][0].InstrumentId != realId {
		t.Fatalf("expected write off instrument id remapped, got %q", writeOff.createdSeveral[0][0].InstrumentId)
	}
}

func TestLoad_SkipsHeaderRow(t *testing.T) {
	svc, instrument, _, _, _, _, _ := newImportFixture()

	header := make([]interface{}, 28)
	header[0] = "Дата поступления"
	xlsx := buildImportXlsx(t, [][]interface{}{header, plainRow("SI-1")})

	if err := svc.LoadOintoSi(context.Background(), importDTO(buildImportFile(t, xlsx))); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(instrument.createdSeveral[0]) != 1 {
		t.Fatalf("expected header to be skipped, got %d rows", len(instrument.createdSeveral[0]))
	}
}

func TestLoad_UnknownBidTypeReturnsError(t *testing.T) {
	svc, _, _, _, _, _, _ := newImportFixture()

	err := svc.Load(context.Background(), &models.ImportDTO{BidType: "unknown"})
	if err == nil || !strings.Contains(err.Error(), "not implemented") {
		t.Fatalf("expected not implemented error, got %v", err)
	}
}
