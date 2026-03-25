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
	MakeAccountingLog(ctx context.Context, dto []*models.SiWithLog) (*bytes.Buffer, error)
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

type ColumnWidthRule struct {
	StartCol int // индекс колонки (0 = A)
	EndCol   int // включительно
	Width    float64
}

func (s *FileService) generateExcel(headers []string, rows [][]interface{}, widthRules []ColumnWidthRule) (*bytes.Buffer, error) {
	file := excelize.NewFile()
	defer file.Close()
	sheetName := "Sheet1"

	// 1. Подготовка стилей
	hStyle, err := file.NewStyle(headerStyle)
	if err != nil {
		return nil, fmt.Errorf("header style err: %w", err)
	}
	mStyle, err := file.NewStyle(cellStyle)
	if err != nil {
		return nil, fmt.Errorf("main style err: %w", err)
	}

	// 2. Пишем заголовки
	if err := file.SetSheetRow(sheetName, "A1", &headers); err != nil {
		return nil, fmt.Errorf("failed to set headers: %w", err)
	}

	// 3. Пишем данные
	for i, item := range rows {
		cell, _ := excelize.CoordinatesToCellName(1, i+2)
		if err := file.SetSheetRow(sheetName, cell, &item); err != nil {
			return nil, fmt.Errorf("failed to set row %d: %w", i+2, err)
		}
	}

	// 4. Финальное оформление (пакетно)
	lastCol, _ := excelize.ColumnNumberToName(len(headers))
	lastRow := len(rows) + 1

	// if err := file.SetColWidth(sheetName, "A", lastCol, 20); err != nil {
	// 	return nil, fmt.Errorf("failed to set column width: %w", err)
	// }
	for _, rule := range widthRules {
		startName, err := excelize.ColumnNumberToName(rule.StartCol + 1)
		if err != nil {
			return nil, fmt.Errorf("invalid start column index %d: %w", rule.StartCol, err)
		}
		endName, err := excelize.ColumnNumberToName(rule.EndCol + 1)
		if err != nil {
			return nil, fmt.Errorf("invalid end column index %d: %w", rule.EndCol, err)
		}
		if err := file.SetColWidth(sheetName, startName, endName, rule.Width); err != nil {
			return nil, fmt.Errorf("failed to set width for %s-%s: %w", startName, endName, err)
		}
	}

	if err := file.SetCellStyle(sheetName, "A1", lastCol+"1", hStyle); err != nil {
		return nil, fmt.Errorf("failed to set header style: %w", err)
	}

	if lastRow > 1 {
		if err = file.SetCellStyle(sheetName, "A2", fmt.Sprintf("%s%d", lastCol, lastRow), mStyle); err != nil {
			return nil, fmt.Errorf("failed to set main style: %w", err)
		}
	}

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to write to buffer. error: %w", err)
	}
	return buffer, nil
}

func (s *FileService) Export(ctx context.Context, dto *models.ExportDTO) (*bytes.Buffer, error) {
	// 1. Плоский список колонок (обработка вложенности)
	flatColumns := make([]*models.Column, 0, len(dto.Columns))
	columnNames := make([]string, 0, len(dto.Columns))

	for _, col := range dto.Columns {
		if len(col.Children) > 0 {
			for _, child := range col.Children {
				columnNames = append(columnNames, child.Name)
				flatColumns = append(flatColumns, child)
			}
		} else {
			columnNames = append(columnNames, col.Name)
			flatColumns = append(flatColumns, col)
		}
	}

	// 2. Подготовка данных через рефлексию
	rows := make([][]interface{}, len(dto.SI))
	for i, item := range dto.SI {
		rowData, err := s.mapReflectRow(item, flatColumns)
		if err != nil {
			return nil, fmt.Errorf("failed to map row %d: %w", i, err)
		}
		rows[i] = rowData
	}

	rules := []ColumnWidthRule{
		{StartCol: 0, EndCol: len(columnNames), Width: 20},
	}

	// 3. Генерация файла через наш общий метод
	return s.generateExcel(columnNames, rows, rules)
}

func (s *FileService) mapReflectRow(d interface{}, columns []*models.Column) ([]interface{}, error) {
	val := reflect.ValueOf(d)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	row := make([]interface{}, 0, len(columns))

	for _, col := range columns {
		// Приводим первую букву к верхнему регистру (Go convention для экспортируемых полей)
		fieldName := strings.ToUpper(col.Field[:1]) + col.Field[1:]
		field := val.FieldByName(fieldName)

		if !field.IsValid() {
			return nil, fmt.Errorf("field %s not found in struct", fieldName)
		}

		// Обработка дат
		if (col.Type == "date" || col.Type == "short_date") && field.Kind() == reflect.Struct {
			if t, ok := field.Interface().(time.Time); ok {
				if t.IsZero() {
					row = append(row, "")
				} else {
					fmtStr := constants.DateFormat
					if col.Type == "short_date" {
						fmtStr = constants.ShortDateFormat
					}
					row = append(row, t.Format(fmtStr))
				}
				continue
			}
		}

		row = append(row, field.Interface())
	}

	return row, nil
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

func (s *FileService) MakeAccountingLog(ctx context.Context, dto []*models.SiWithLog) (*bytes.Buffer, error) {
	headers := []string{
		"Дата поступл.", "Наименование СИ", "Вид (тип, марка) СИ", "Зав. №",
		"Лицо, ответств. за экспл. СИ", "Сведения о ремонте", "Сведения о консервации",
		"Сведения о передаче на хранение", "Сведения о списании",
	}
	rules := []ColumnWidthRule{
		{StartCol: 0, EndCol: len(headers), Width: 20},
	}

	rows := make([][]interface{}, len(dto))
	for i, d := range dto {
		date := d.DateOfReceipt.Format(constants.ShortDateFormat)
		if d.DateOfReceipt.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local)) {
			date = "–"
		}

		rows[i] = []interface{}{
			date,
			d.Name, d.Type, d.FactoryNumber, d.Responsible,
			d.RepairInfo, d.PreservationInfo, d.SavingInfo, d.WriteOff,
		}
	}

	return s.generateExcel(headers, rows, rules)
}

func (s *FileService) MakeVerificationSchedule(ctx context.Context, dto *models.SiVerification) (*bytes.Buffer, error) {
	strategy := s.resolveReportStrategy(dto.BidType)

	rows := make([][]interface{}, len(dto.SI))
	for i, item := range dto.SI {
		rows[i] = strategy.mapFunc(i, item)
	}

	return s.generateExcel(strategy.headers, rows, strategy.widthRules)
}

// Вспомогательная структура для инкапсуляции логики типов отчетов
type reportStrategy struct {
	headers    []string
	widthRules []ColumnWidthRule
	mapFunc    func(int, *models.SI) []interface{}
}

func (s *FileService) resolveReportStrategy(bidType string) reportStrategy {
	switch {
	case bidType == "met_si" || bidType == "eq_si":
		headers := []string{"№ п/п", "Наименование", "Тип СИ", "Заводской номер", "Диапазон измерений", "Периодичность поверки", "Дата последней поверки", "Дата следующей поверки", "Примечание"}
		return reportStrategy{
			headers: headers,
			widthRules: []ColumnWidthRule{
				{StartCol: 0, EndCol: len(headers), Width: 20},
			},
			mapFunc: func(i int, d *models.SI) []interface{} {
				return []interface{}{
					i + 1, d.Name, d.Type, d.FactoryNumber, d.MeasurementLimits, d.InterVerificationInterval,
					d.VerificationDate.Format(constants.DateFormat),
					d.NextVerificationDate.Format(constants.DateFormat),
					d.Notes,
				}
			},
		}

	case strings.Contains(bidType, "eq"):
		return reportStrategy{
			headers: []string{"№ п/п", "Наименование", "Марка, тип", "Заводской номер", "Интервал аттестации, мес.", "Дата последней аттестации", "Дата следующей аттестации", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "Примечание"},
			widthRules: []ColumnWidthRule{
				{StartCol: 0, EndCol: 6, Width: 20},
				{StartCol: 7, EndCol: 18, Width: 5},
				{StartCol: 19, EndCol: 19, Width: 20},
			},
			mapFunc: func(i int, d *models.SI) []interface{} {
				return s.mapWithMonthlyGrid(i, d, constants.ShortDateFormat)
			},
		}

	default:
		return reportStrategy{
			headers: []string{"№ п/п", "Наименование", "Марка, тип", "Заводской номер", "Интервал поверки, мес.", "Дата последней поверки", "Дата следующей поверки", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII", "Примечание"},
			widthRules: []ColumnWidthRule{
				{StartCol: 0, EndCol: 6, Width: 20},
				{StartCol: 7, EndCol: 18, Width: 5},
				{StartCol: 19, EndCol: 19, Width: 20},
			},
			mapFunc: func(i int, d *models.SI) []interface{} {
				return s.mapWithMonthlyGrid(i, d, constants.ShortDateFormat)
			},
		}
	}
}

// Вынес общую логику для отчетов с сеткой месяцев (I-XII)
func (s *FileService) mapWithMonthlyGrid(i int, d *models.SI, dateFmt string) []interface{} {
	res := []interface{}{
		i + 1, d.Name, d.Type, d.FactoryNumber, d.InterVerificationInterval,
		d.VerificationDate.Format(dateFmt),
		d.NextVerificationDate.Format(dateFmt),
	}

	// Сетка месяцев
	months := make([]interface{}, 12)
	for m := 0; m < 12; m++ {
		months[m] = ""
	}

	mIdx := int(d.NextVerificationDate.Month()) - 1
	if mIdx >= 0 && mIdx < 12 {
		months[mIdx] = "*"
	}

	res = append(res, months...)
	res = append(res, d.Notes)
	return res
}
