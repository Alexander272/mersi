package services

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/gomutex/godocx"
	"github.com/xuri/excelize/v2"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

type File interface {
	Export(ctx context.Context, dto *models.ExportDTO) (*bytes.Buffer, error)
	MakeDocSchedule(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error)
	MakeVerificationSchedule(ctx context.Context, dto *models.SiVerification) (*bytes.Buffer, error)
}

var headerStyle = &excelize.Style{
	Border: []excelize.Border{
		{Type: "top", Style: 1, Color: "#000000"},
		{Type: "bottom", Style: 1, Color: "#000000"},
		{Type: "left", Style: 1, Color: "#000000"},
		{Type: "right", Style: 1, Color: "#000000"},
	},
	Alignment: &excelize.Alignment{
		WrapText:   true,
		Horizontal: "center",
	},
}
var cellStyle = &excelize.Style{
	Border: []excelize.Border{
		{Type: "top", Style: 1, Color: "#000000"},
		{Type: "bottom", Style: 1, Color: "#000000"},
		{Type: "left", Style: 1, Color: "#000000"},
		{Type: "right", Style: 1, Color: "#000000"},
	},
	Alignment: &excelize.Alignment{
		WrapText: true,
		Vertical: "center",
	},
}

func (s *FileService) Export(ctx context.Context, dto *models.ExportDTO) (buffer *bytes.Buffer, err error) {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(file.GetActiveSheetIndex())

	headerStyle, err := file.NewStyle(headerStyle)
	if err != nil {
		return nil, fmt.Errorf("failed to create header style. error: %w", err)
	}

	mainStyle, err := file.NewStyle(cellStyle)
	if err != nil {
		return nil, fmt.Errorf("failed to create main style. error: %w", err)
	}

	flatColumns := []*models.Column{}
	columnNames := []string{}
	for _, column := range dto.Columns {
		if len(column.Children) > 0 {
			for _, child := range column.Children {
				columnNames = append(columnNames, child.Name)
				flatColumns = append(flatColumns, child)
			}
			continue
		}
		columnNames = append(columnNames, column.Name)
		flatColumns = append(flatColumns, column)
	}
	if err := file.SetSheetRow(sheetName, "A1", &columnNames); err != nil {
		return nil, fmt.Errorf("failed to set header row. error: %w", err)
	}

	endColumn, err := excelize.ColumnNumberToName(len(columnNames))
	if err != nil {
		return nil, fmt.Errorf("failed to get column name. error: %w", err)
	}

	if err := file.SetColWidth(sheetName, "A", endColumn, 25); err != nil {
		return nil, fmt.Errorf("failed to set column width. error: %w", err)
	}
	if err = file.SetCellStyle(sheetName, "A1", endColumn+"1", headerStyle); err != nil {
		return nil, fmt.Errorf("failed to set header style. error: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("unexpected error: %v", r)
		}
	}()

	for i, d := range dto.SI {
		values := []interface{}{}
		data := reflect.ValueOf(d)

		if data.Kind() == reflect.Ptr {
			data = data.Elem() // Dereference the pointer
		}

		for _, c := range flatColumns {
			// if len(c.Children) > 0 {
			// 	for _, child := range c.Children {
			// 		field := data.FieldByName(child.Field)

			// 		if field.IsValid() {
			// 			values = append(values, field.Interface())
			// 		} else {
			// 			return nil, fmt.Errorf("field %s not found", child.Field)
			// 		}
			// 	}
			// 	continue
			// }

			field := strings.ToUpper(c.Field[:1]) + c.Field[1:]
			value := data.FieldByName(field)

			if value.IsValid() {
				if value.Kind() == reflect.Int64 && c.Type == "date" {
					newValue := value.Int()
					if newValue == 0 {
						values = append(values, "")
						continue
					}

					date := time.Unix(newValue, 0).Format(constants.DateFormat)
					values = append(values, date)
					continue
				}

				values = append(values, value.Interface())
			} else {
				return nil, fmt.Errorf("field %s not found", field)
			}
		}

		// values := []interface{}{
		// 	d.Name, d.Type, d.FactoryNumber, d.MeasurementLimits, d.Accuracy, d.StateRegister, d.Manufacturer, d.YearOfIssue, d.Date,
		// 	d.InterVerificationInterval, d.NextDate, d.Place, d.Person, d.Notes,
		// }

		if err := file.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &values); err != nil {
			return nil, fmt.Errorf("failed to set header row. error: %w", err)
		}
		if err = file.SetCellStyle(sheetName, fmt.Sprintf("A%d", i+2), fmt.Sprintf("%s%d", endColumn, i+2), mainStyle); err != nil {
			return nil, fmt.Errorf("failed to set style. error: %w", err)
		}
	}

	buffer, err = file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer. error: %w", err)
	}
	return buffer, nil
}

func (s *FileService) MakeDocSchedule(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error) {
	document, err := godocx.NewDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to create document. error: %w", err)
	}

	titles := []string{"№ п/п", "Наименование СИ", "Вид (тип, марка) СИ", "Зав. №", "Комплект СИ, штук", "Комплект СИ, набор", "Год выпуска СИ", "№ Госреестра (рег. номер в ФИФ)", "Метрологические характеристики  (диапазон измерений, разряд, погрешность, класс точности)", "Запрос на поверку СИ в качестве эталона с предоставлением протокола поверки", "Вид поверки"}

	document.AddHeading("Поверка", 1)
	document.AddEmptyParagraph()
	table := document.AddTable()
	table.Style("TableGrid")
	header := table.AddRow()
	for _, t := range titles {
		header.AddCell().AddParagraph(t).Justification("center")
	}

	for i, item := range dto {
		row := table.AddRow()

		row.AddCell().AddParagraph(fmt.Sprintf("%d", i+1)).Justification("center")
		row.AddCell().AddParagraph(item.Name).Justification("center")
		row.AddCell().AddParagraph(item.Type).Justification("center")
		row.AddCell().AddParagraph(item.FactoryNumber).Justification("center")
		row.AddCell().AddParagraph("1").Justification("center")
		row.AddCell().AddParagraph("-").Justification("center")
		row.AddCell().AddParagraph(strconv.Itoa(item.YearOfIssue)).Justification("center")
		row.AddCell().AddParagraph(item.StateRegister).Justification("center")
		row.AddCell().AddParagraph(item.MeasurementLimits).Justification("center")
		row.AddCell().AddParagraph("Нет").Justification("center")
		row.AddCell().AddParagraph("Периодическая").Justification("center")
	}

	buffer := new(bytes.Buffer)
	if _, err := document.WriteTo(buffer); err != nil {
		return nil, fmt.Errorf("failed to save document. error: %w", err)
	}
	return buffer, nil
}

func (s *FileService) MakeVerificationSchedule(ctx context.Context, dto *models.SiVerification) (*bytes.Buffer, error) {
	file := excelize.NewFile()
	sheetName := file.GetSheetName(file.GetActiveSheetIndex())

	headerStyle, err := file.NewStyle(headerStyle)
	if err != nil {
		return nil, fmt.Errorf("failed to create header style. error: %w", err)
	}

	mainStyle, err := file.NewStyle(cellStyle)
	if err != nil {
		return nil, fmt.Errorf("failed to create main style. error: %w", err)
	}

	columnNames := make([]string, 0, 10)
	switch dto.BidType {
	case "ointo_si":
		columnNames = []string{
			"№ п/п", "Наименование", "Заводской номер", "Производитель", "Метрологические характеристики СИ", "Периодичность поверки",
			"Дата последней поверки", "Дата следующей поверки", "Примечание",
		}
	case "ointo_eq":
		columnNames = []string{
			"№ п/п", "Наименование", "Заводской номер", "Производитель", "Периодичность аттестации",
			"Дата последней аттестации", "Дата следующей аттестации", "Примечание",
		}
	default:
		columnNames = []string{
			"№ п/п", "Наименование", "Заводской номер", "Диапазон измерений", "Периодичность поверки", "Дата последней поверки",
			"Дата следующей поверки", "Примечание",
		}
	}

	if err := file.SetSheetRow(sheetName, "A1", &columnNames); err != nil {
		return nil, fmt.Errorf("failed to set header row. error: %w", err)
	}

	endColumn, err := excelize.ColumnNumberToName(len(columnNames))
	if err != nil {
		return nil, fmt.Errorf("failed to get column name. error: %w", err)
	}

	if err := file.SetColWidth(sheetName, "A", endColumn, 25); err != nil {
		return nil, fmt.Errorf("failed to set column width. error: %w", err)
	}
	if err = file.SetCellStyle(sheetName, "A1", endColumn+"1", headerStyle); err != nil {
		return nil, fmt.Errorf("failed to set header style. error: %w", err)
	}

	for i, d := range dto.SI {
		values := make([]interface{}, 0, len(columnNames))
		date := time.Unix(d.VerificationDate, 0).Format(constants.DateFormat)
		nextDate := time.Unix(d.NextVerificationDate, 0).Format(constants.DateFormat)

		switch dto.BidType {
		case "ointo_si":
			values = []interface{}{
				i + 1, d.Name, d.FactoryNumber, d.Manufacturer, d.MeasurementLimits, d.InterVerificationInterval,
				date, nextDate, d.Notes,
			}
		case "ointo_eq":
			values = []interface{}{
				i + 1, d.Name, d.FactoryNumber, d.Manufacturer, d.InterVerificationInterval,
				date, nextDate, d.Notes,
			}
		default:
			values = []interface{}{
				i + 1, d.Name, d.FactoryNumber, d.MeasurementLimits, d.InterVerificationInterval,
				date, nextDate, d.Notes,
			}
		}

		if err := file.SetSheetRow(sheetName, fmt.Sprintf("A%d", i+2), &values); err != nil {
			return nil, fmt.Errorf("failed to set row. error: %w", err)
		}
		if err = file.SetCellStyle(sheetName, fmt.Sprintf("A%d", i+2), fmt.Sprintf("%s%d", endColumn, i+2), mainStyle); err != nil {
			return nil, fmt.Errorf("failed to set style. error: %w", err)
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer. error: %w", err)
	}
	return buffer, nil
}
