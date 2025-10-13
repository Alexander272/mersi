package services

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/xuri/excelize/v2"
)

type ImportService struct {
	instrument     Instrument
	verification   Verification
	repair         Repair
	preservation   Preservation
	transferToSave TransferToSave
	transferToDep  TransferToDepartment
	writeOff       WriteOff
}

type ImportDeps struct {
	Instrument     Instrument
	Verification   Verification
	Repair         Repair
	Preservation   Preservation
	TransferToSave TransferToSave
	TransferToDep  TransferToDepartment
	WriteOff       WriteOff
}

func NewImportService(deps *ImportDeps) *ImportService {
	return &ImportService{
		instrument:     deps.Instrument,
		verification:   deps.Verification,
		repair:         deps.Repair,
		preservation:   deps.Preservation,
		transferToSave: deps.TransferToSave,
		transferToDep:  deps.TransferToDep,
		writeOff:       deps.WriteOff,
	}
}

type ImportFile interface {
	Load(ctx context.Context, dto *models.ImportDTO) error
}

func (s *ImportService) Load(ctx context.Context, dto *models.ImportDTO) error {
	switch dto.BidType {
	case "ointo_si":
		return s.LoadOintoSi(ctx, dto)
	}
	return fmt.Errorf("not implemented")
}

func (s *ImportService) LoadOintoSi(ctx context.Context, dto *models.ImportDTO) error {
	file, err := dto.File.Open()
	if err != nil {
		return fmt.Errorf("failed to open file. error: %w", err)
	}
	defer file.Close()

	excel, err := excelize.OpenReader(file)
	if err != nil {
		return fmt.Errorf("failed to open excel file. error: %w", err)
	}
	defer excel.Close()

	template := &models.Template{
		Name:                      1,
		DateOfReceipt:             0,
		Type:                      2,
		FactoryNumber:             3,
		MeasurementLimits:         13,
		StateRegister:             12,
		CountryOfProduce:          10,
		Manufacturer:              11,
		Responsible:               4,
		YearOfIssue:               9,
		InterVerificationInterval: 14,
		VerificationDate:          15,
		Repair:                    5,
		Preservation:              6,
		// Transfer:                  7,
		Transfer: 7,
		WriteOff: 8,
	}

	siSheet := "si"

	si := []*models.InstrumentDTO{}
	repairs := []*models.RepairDTO{}
	archive := []*models.PreservationDTO{}
	// transferToSave := []*models.TransferToSaveDTO{}
	transfer := []*models.TransferToDepartmentDTO{}
	writeOff := []*models.WriteOffDTO{}
	verificationsByIndex := make(map[int][]string, 0)
	verifications := []*models.VerificationDTO{}
	index := 0

	dateRe := regexp.MustCompile(`\d{2}\.\d{2}\.\d{4}`)
	descRe := regexp.MustCompile(`\((.*?)\)`)

	rows, err := excel.Rows(siSheet)
	if err != nil {
		return fmt.Errorf("failed to get rows. error: %w", err)
	}
	for rows.Next() {
		origRow, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns. error: %w", err)
		}
		row := make([]string, 28)
		copy(row, origRow)

		if len(row) == 0 || row[0] == "" || row[0] == "Дата поступления" {
			continue
		}
		// logger.Debug("import", logger.StringAttr("row", strings.Join(row, ", ")))

		dateOfReceipt := time.Time{}
		if row[template.DateOfReceipt] != "" {
			date, err := time.Parse("02.01.2006", row[template.DateOfReceipt])
			if err != nil {
				return fmt.Errorf("failed to parse date of receipt. error: %w", err)
			}
			dateOfReceipt = time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), 0, 0, 0, time.Now().Location())
		}

		year := 0
		if row[template.YearOfIssue] != "" {
			year, err = strconv.Atoi(row[template.YearOfIssue])
			if err != nil {
				return fmt.Errorf("failed to parse year of issue. error: %w", err)
			}
		}
		interval := 0
		if row[template.InterVerificationInterval] != "" {
			interval, err = strconv.Atoi(row[template.InterVerificationInterval])
			if err != nil {
				return fmt.Errorf("failed to parse interval. error: %w", err)
			}
		}

		status := models.InstrumentStatusWork

		repairString := row[template.Repair]
		list := []string{}
		if repairString != "" && repairString != "−" && repairString != "-" {
			list = strings.Split(repairString, ";")
		}
		for _, l := range list {
			// logger.Debug("Repair", logger.StringAttr("line", l))

			parts := strings.Split(l, "-")
			dateString := dateRe.FindString(parts[0])
			startDate, err := time.Parse("02.01.2006", dateString)
			if err != nil {
				return fmt.Errorf("failed to parse start date of repair. error: %w", err)
			}
			endDate := time.Time{}
			if len(parts) > 1 {
				dateString := dateRe.FindString(parts[1])
				endDate, err = time.Parse("02.01.2006", dateString)
				if err != nil {
					return fmt.Errorf("failed to parse end date of repair. error: %w", err)
				}
			}

			work := ""
			matches := descRe.FindStringSubmatch(l)
			if len(matches) > 1 {
				work = matches[1]
			}

			if endDate.IsZero() {
				status = models.InstrumentStatusRepair
			}

			repairs = append(repairs, &models.RepairDTO{
				InstrumentId: strconv.Itoa(index),
				PeriodStart:  startDate,
				PeriodEnd:    endDate,
				Work:         strings.TrimSpace(work),
			})
		}

		archiveString := row[template.Preservation]
		list = []string{}
		if archiveString != "" && archiveString != "-" && archiveString != "−" {
			list = strings.Split(archiveString, ";")
		}
		for _, l := range list {
			// logger.Debug("Archive", logger.StringAttr("line", l))

			parts := strings.Split(l, "-")

			dateString := dateRe.FindString(parts[0])
			startDate, err := time.Parse("02.01.2006", dateString)
			if err != nil {
				return fmt.Errorf("failed to parse start date of preservation. error: %w", err)
			}
			noteStart := ""
			matches := descRe.FindStringSubmatch(l)
			if len(matches) > 1 {
				noteStart = matches[1]
			}

			endDate := time.Time{}
			nodeEnd := ""
			if len(parts) > 1 {
				dateString := dateRe.FindString(parts[1])
				endDate, err = time.Parse("02.01.2006", dateString)
				if err != nil {
					return fmt.Errorf("failed to parse end date of preservation. error: %w", err)
				}
				matches := descRe.FindStringSubmatch(l)
				if len(matches) > 1 {
					nodeEnd = matches[1]
				}
			}
			if endDate.IsZero() {
				status = models.InstrumentStatusArchived
			}

			archive = append(archive, &models.PreservationDTO{
				InstrumentId: strconv.Itoa(index),
				DateStart:    startDate,
				NotesStart:   strings.TrimSpace(noteStart),
				DateEnd:      endDate,
				NotesEnd:     strings.TrimSpace(nodeEnd),
			})
		}

		if row[template.WriteOff] != "" && row[template.WriteOff] != "-" && row[template.WriteOff] != "−" {
			status = models.InstrumentStatusDec
			dateString := dateRe.FindString(row[template.WriteOff])
			date := time.Time{}
			if dateString != "" {
				date, err = time.Parse("02.01.2006", dateString)
				if err != nil {
					return fmt.Errorf("failed to parse date of repair. error: %w", err)
				}
			}

			writeOff = append(writeOff, &models.WriteOffDTO{
				InstrumentId: strconv.Itoa(index),
				Date:         date,
				DocName:      strings.TrimSpace(row[template.WriteOff]),
			})
		}

		if row[template.Transfer] != "" && row[template.Transfer] != "-" && row[template.Transfer] != "−" {
			status = models.InstrumentStatusTransferred
			dateString := dateRe.FindString(row[template.Transfer])
			date := time.Time{}
			if dateString != "" {
				date, err = time.Parse("02.01.2006", dateString)
				if err != nil {
					return fmt.Errorf("failed to parse date of transfer. error: %w", err)
				}
			}

			transfer = append(transfer, &models.TransferToDepartmentDTO{
				InstrumentId: strconv.Itoa(index),
				Date:         date,
				DocName:      strings.TrimSpace(row[template.Transfer]),
			})
		}

		si = append(si, &models.InstrumentDTO{
			SectionId:                 dto.SectionId,
			UserId:                    dto.UserId,
			Name:                      strings.TrimSpace(row[template.Name]),
			DateOfReceipt:             dateOfReceipt,
			Type:                      strings.TrimSpace(row[template.Type]),
			FactoryNumber:             strings.TrimSpace(row[template.FactoryNumber]),
			MeasurementLimits:         strings.TrimSpace(row[template.MeasurementLimits]),
			StateRegister:             strings.TrimSpace(row[template.StateRegister]),
			CountryOfProduce:          strings.TrimSpace(row[template.CountryOfProduce]),
			Manufacturer:              strings.TrimSpace(row[template.Manufacturer]),
			Responsible:               strings.TrimSpace(row[template.Responsible]),
			YearOfIssue:               year,
			InterVerificationInterval: interval,
			Status:                    status,
		})

		for i := template.VerificationDate; i < len(row); i++ {
			item := row[i]
			// logger.Debug("Verification", logger.StringAttr("item", item))
			if item != "" && item != "-" && item != "−" {
				verificationsByIndex[index] = append(verificationsByIndex[index], strings.Split(item, ";")...)
			}
		}
		// logger.Debug("Verifications", logger.AnyAttr("array", verificationsByIndex[index]))

		index++
	}

	if err := s.instrument.CreateSeveral(ctx, si); err != nil {
		return fmt.Errorf("failed to create several instruments. error: %w", err)
	}

	for i := range verificationsByIndex {
		if len(verificationsByIndex[i]) == 0 {
			continue
		}

		for _, verString := range verificationsByIndex[i] {
			dateString := dateRe.FindString(verString)
			date := time.Time{}
			if dateString != "" {
				date, err = time.Parse("02.01.2006", dateString)
				if err != nil {
					return fmt.Errorf("failed to parse date of verification. error: %w", err)
				}
			}
			nextDate := date.AddDate(0, si[i].InterVerificationInterval, 0)

			status := "work"
			if strings.Contains(strings.ToLower(verString), "списан") || strings.Contains(strings.ToLower(verString), "непригод") {
				status = "decommissioning"
			}

			verifications = append(verifications, &models.VerificationDTO{
				InstrumentId: si[i].Id,
				UserId:       dto.UserId,
				Date:         date,
				NextDate:     nextDate,
				Status:       status,
				Docs:         []*models.VerificationDocDTO{{Name: verString}},
			})
		}
	}

	if err := s.verification.CreateSeveral(ctx, verifications); err != nil {
		return fmt.Errorf("failed to create verifications. error: %w", err)
	}

	for i := range repairs {
		index, _ := strconv.Atoi(repairs[i].InstrumentId)
		repairs[i].InstrumentId = si[index].Id
	}
	if err := s.repair.CreateSeveral(ctx, repairs); err != nil {
		return fmt.Errorf("failed to create repair. error: %w", err)
	}

	for i := range archive {
		index, _ := strconv.Atoi(archive[i].InstrumentId)
		archive[i].InstrumentId = si[index].Id
	}
	if err := s.preservation.CreateSeveral(ctx, archive); err != nil {
		return fmt.Errorf("failed to create archive. error: %w", err)
	}

	for i := range transfer {
		index, _ := strconv.Atoi(transfer[i].InstrumentId)
		transfer[i].InstrumentId = si[index].Id
	}
	if err := s.transferToDep.CreateSeveral(ctx, transfer); err != nil {
		return fmt.Errorf("failed to create transfer. error: %w", err)
	}

	for i := range writeOff {
		index, _ := strconv.Atoi(writeOff[i].InstrumentId)
		writeOff[i].InstrumentId = si[index].Id
	}
	if err := s.writeOff.CreateSeveral(ctx, writeOff); err != nil {
		return fmt.Errorf("failed to create write off. error: %w", err)
	}

	return nil
}
