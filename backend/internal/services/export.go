package services

import (
	"bytes"
	"context"

	"github.com/Alexander272/mersi/backend/internal/models"
)

type ExportService struct {
	file File
	si   SI
}

func NewExportService(file File, si SI) *ExportService {
	return &ExportService{
		file: file,
		si:   si,
	}
}

type Export interface {
	MakeScheduler(ctx context.Context, req *models.Period) (*bytes.Buffer, error)
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
