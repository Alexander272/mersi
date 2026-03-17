package services

import (
	"bytes"
	"context"
	"fmt"

	"github.com/Alexander272/mersi/backend/internal/models"
)

type ExportService struct {
	file    File
	si      SI
	columns Columns
}

type ExportDeps struct {
	File    File
	SI      SI
	Columns Columns
}

func NewExportService(deps *ExportDeps) *ExportService {
	return &ExportService{
		file:    deps.File,
		si:      deps.SI,
		columns: deps.Columns,
	}
}

type Export interface {
	Export(ctx context.Context, req *models.GetSiDTO) (*bytes.Buffer, error)
	MakeScheduler(ctx context.Context, req *models.Period) (*bytes.Buffer, error)
}

func (s *ExportService) Export(ctx context.Context, req *models.GetSiDTO) (*bytes.Buffer, error) {
	columns, err := s.columns.Get(ctx, &models.GetColumnsDTO{SectionID: req.SectionId})
	if err != nil {
		return nil, err
	}

	si, err := s.si.Get(ctx, req)
	if err != nil {
		return nil, err
	}

	buffer, err := s.file.Export(ctx, &models.ExportDTO{Columns: columns, SI: si})
	if err != nil {
		return nil, fmt.Errorf("failed to export. error: %w", err)
	}
	return buffer, nil
}

func (s *ExportService) MakeAccLog(ctx context.Context, req *models.Period) (*bytes.Buffer, error) {
	return nil, nil
}

func (s *ExportService) MakeScheduler(ctx context.Context, req *models.Period) (*bytes.Buffer, error) {
	data, err := s.si.GetVerification(ctx, req)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, models.ErrNoRows
	}

	buffer, err := s.file.MakeVerificationSchedule(ctx, data[0])
	if err != nil {
		return nil, err
	}

	return buffer, nil
}
