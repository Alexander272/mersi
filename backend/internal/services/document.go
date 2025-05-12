package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/google/uuid"
)

type DocumentService struct {
	repo repository.Document
	path string
}

func NewDocumentService(repo repository.Document) *DocumentService {
	return &DocumentService{
		repo: repo,
		path: "files/si/",
	}
}

type Document interface {
	GetTemp(ctx context.Context, req *models.GetTempDocumentDTO) ([]*models.Document, error)
	Upload(ctx context.Context, dto *models.DocumentsDTO) ([]*models.Document, error)
	ChangePath(ctx context.Context, req *models.PathParts) error
	Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error
	DeleteByInstrumentId(ctx context.Context, instrumentId string) error
}

func (s *DocumentService) GetTemp(ctx context.Context, req *models.GetTempDocumentDTO) ([]*models.Document, error) {
	data, err := s.repo.GetTemp(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get temp documents. error: %w", err)
	}
	return data, nil
}

// SaveUploadedFile uploads the form file to specific dst.
func (s *DocumentService) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (s *DocumentService) Upload(ctx context.Context, dto *models.DocumentsDTO) ([]*models.Document, error) {
	docs := []*models.Document{}

	documentTypes := map[string]string{
		"application/msword":          "doc",
		"application/x-extension-doc": "doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "doc",
		"application/x-extension-docx":                                            "doc",
		"application/vnd.oasis.opendocument.text":                                 "doc",
		"application/vnd.ms-excel":                                                "sheet",
		"application/x-extension-xls":                                             "sheet",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       "sheet",
		"application/x-extension-xlsx":                                            "sheet",
		"application/pdf":                                                         "pdf",
		"image/png":                                                               "image",
		"image/jpeg":                                                              "image",
		"text/csv":                                                                "sheet",
	}

	for _, fh := range dto.Files {
		doc := &models.Document{
			Id:           uuid.NewString(),
			Label:        fh.Filename,
			Size:         fh.Size,
			InstrumentId: dto.InstrumentId,
			UserId:       dto.UserId,
			Group:        dto.Group,
			DocumentType: documentTypes[fh.Header.Get("Content-Type")],
		}

		paths := []string{s.path}
		if dto.InstrumentId != "" {
			paths = append(paths, dto.Group, dto.InstrumentId)
		} else {
			paths = append(paths, "temp", dto.UserId, dto.Group)
		}
		paths = append(paths, doc.Id, fh.Filename)

		dst := path.Join(paths...)
		doc.Path = dst
		docs = append(docs, doc)

		if err := s.SaveUploadedFile(fh, dst); err != nil {
			return nil, fmt.Errorf("failed to save file. error: %w", err)
		}
	}

	if err := s.repo.CreateSeveral(ctx, docs); err != nil {
		return nil, fmt.Errorf("failed to create documents. error: %w", err)
	}
	return docs, nil
}

func (s *DocumentService) ChangePath(ctx context.Context, req *models.PathParts) error {
	count, err := s.repo.UpdatePath(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update path documents. error: %w", err)
	}

	if count > 0 {
		newPath := path.Join(s.path, req.Group, req.InstrumentId)
		//TODO думаю тут еще надо id пользователя использовать (чтобы во время одновременного создания ничего лишнего не попало)
		// по хорошему это еще надо синхронизировать между устройствами
		// еще можно добавить какую-нибудь группировку чтобы файлы из разных мест не пересекались или она мне не нужна
		srcPath := path.Join(s.path, "temp", req.UserId, req.Group)

		if err = os.MkdirAll(filepath.Dir(newPath), 0750); err != nil {
			return err
		}

		if err := os.Rename(srcPath, newPath); err != nil {
			return fmt.Errorf("failed to move files. error: %w", err)
		}
	}
	return nil
}

func (s *DocumentService) Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error {
	paths := []string{s.path}
	if dto.InstrumentId != "" {
		paths = append(paths, dto.Group, dto.InstrumentId)
	} else {
		paths = append(paths, "temp", dto.UserId, dto.Group)
	}
	paths = append(paths, dto.Id, dto.Filename)

	dst := path.Join(paths...)

	if err := os.Remove(dst); err != nil && !strings.Contains(err.Error(), "no such file") {
		return fmt.Errorf("failed to delete file. error: %w", err)
	}

	//TODO надо бы еще удалять пустые директории (можно это делать раз в день с помощью go-cron вызывать функцию для удаления)

	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete document by id. error: %w", err)
	}
	return nil
}

func (s *DocumentService) DeleteByInstrumentId(ctx context.Context, instrumentId string) error {
	dst := path.Join(s.path, instrumentId)

	if err := os.RemoveAll(dst); err != nil && !strings.Contains(err.Error(), "no such file") {
		return fmt.Errorf("failed to delete folder with files. error: %w", err)
	}
	return nil
}
