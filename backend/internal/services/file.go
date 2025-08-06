package services

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/gomutex/godocx"
	"github.com/xuri/excelize/v2"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

type File interface {
	Export(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error)
	MakeDocSchedule(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error)
	MakeVerificationSchedule(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error)
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

func (s *FileService) Export(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error) {
	// file := excelize.NewFile()
	// sheetName := file.GetSheetName(file.GetActiveSheetIndex())

	// headerStyle, err := file.NewStyle(headerStyle)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create header style. error: %w", err)
	// }

	// mainStyle, err := file.NewStyle(cellStyle)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create main style. error: %w", err)
	// }

	// //TODO колонки же зависят от области
	// columnNames := []string{
	// 	"Наименование", "Тип СИ", "Заводской номер", "Пределы измерений", "Точность, цена деления, погрешность", "Госреестр СИ", "Изготовитель",
	// 	"Год выпуска", "Дата поверки (калибровки)", "Межповерочный интервал", "Следующая поверка (калибровка)", "Место нахождения", "ФИО Сотрудника", "Примечание",
	// }
	// if err := file.SetSheetRow(sheetName, "A1", &columnNames); err != nil {
	// 	return nil, fmt.Errorf("failed to set header row. error: %w", err)
	// }

	// endColumn, err := excelize.ColumnNumberToName(len(columnNames))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get column name. error: %w", err)
	// }

	// if err := file.SetColWidth(sheetName, "A", endColumn, 25); err != nil {
	// 	return nil, fmt.Errorf("failed to set column width. error: %w", err)
	// }
	// if err = file.SetCellStyle(sheetName, "A1", endColumn+"1", headerStyle); err != nil {
	// 	return nil, fmt.Errorf("failed to set header style. error: %w", err)
	// }

	return nil, fmt.Errorf("not implemented")
}

func (s *FileService) MakeDocSchedule(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error) {
	document, err := godocx.NewDocument()
	if err != nil {
		return nil, fmt.Errorf("failed to create document. error: %w", err)
	}

	titles := []string{"№ п/п", "Наименование СИ", "Вид (тип, марка) СИ", "Зав. №", "Комплект СИ, штук", "Комплект СИ, набор", "Год выпуска СИ", "№ Госреестра (рег. номер в ФИФ)", "Метрологические характеристики  (диапазон измерений, разряд, погрешность, класс точности)", "Запрос на поверку СИ в качестве эталона с предоставлением протокола поверки", "Вид поверки"}

	document.AddHeading("Поверка", 1)
	table := document.AddTable()
	header := table.AddRow()
	for _, t := range titles {
		header.AddCell().AddParagraph(t)
	}

	for i, item := range dto {
		row := table.AddRow()

		row.AddCell().AddParagraph(fmt.Sprintf("%d", i+1))
		row.AddCell().AddParagraph(item.Name)
		row.AddCell().AddParagraph(item.Type)
		row.AddCell().AddParagraph(item.FactoryNumber)
		row.AddCell().AddParagraph("1")
		row.AddCell().AddParagraph("-")
		row.AddCell().AddParagraph(strconv.Itoa(item.YearOfIssue))
		row.AddCell().AddParagraph(item.StateRegister)
		row.AddCell().AddParagraph(item.MeasurementLimits)
		row.AddCell().AddParagraph("Нет")
		row.AddCell().AddParagraph("Периодическая")
	}

	buffer := new(bytes.Buffer)
	if _, err := document.WriteTo(buffer); err != nil {
		return nil, fmt.Errorf("failed to save document. error: %w", err)
	}
	return buffer, nil
}

func (s *FileService) MakeVerificationSchedule(ctx context.Context, dto []*models.SI) (*bytes.Buffer, error) {
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

	columnNames := []string{
		"№ п/п", "Наименование", "Заводской номер", "Диапазон измерений", "Периодичность поверки", "Дата последней поверки",
		"Дата следующей поверки", "Примечание",
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

	for i, d := range dto {
		values := []interface{}{
			i + 1, d.Name, d.FactoryNumber, d.MeasurementLimits, d.InterVerificationInterval,
			d.VerificationDate, d.NextVerificationDate, d.Notes,
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
