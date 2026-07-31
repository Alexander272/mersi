package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/repository"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/google/uuid"
	"github.com/karrick/godirwalk"
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

var (
	allowedExtensions = map[string]bool{
		".doc": true, ".docx": true, ".pdf": true,
		".jpg": true, ".jpeg": true, ".png": true,
		".xls": true, ".xlsx": true, ".csv": true,
	}

	documentTypes = map[string]string{
		"application/msword":          "doc",
		"application/x-extension-doc": "doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": "doc",
		"application/x-extension-docx":                                            "doc",
		"application/vnd.oasis.opendocument.text":                                 "doc",
		"application/vnd.ms-excel":                                                "sheet",
		"application/x-extension-xls":                                             "sheet",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       "sheet",
		"application/x-extension-xlsx":                                            "sheet",
		"text/csv":                                                                "sheet",
		"application/pdf":                                                         "pdf",
		"image/png":                                                               "image",
		"image/jpeg":                                                              "image",
	}
)

type Document interface {
	GetTemp(ctx context.Context, req *models.GetDocumentDTO) ([]*models.Document, error)
	GetByInstrument(ctx context.Context, req *models.GetDocumentDTO) ([]*models.Document, error)
	Upload(ctx context.Context, dto *models.DocumentsDTO) ([]*models.Document, error)
	ChangePath(ctx context.Context, req *models.PathParts) error
	Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error
	DeleteByInstrumentId(ctx context.Context, instrumentId string) error
	RemoveEmptyFolders(ctx context.Context) error
}

func (s *DocumentService) GetTemp(ctx context.Context, req *models.GetDocumentDTO) ([]*models.Document, error) {
	data, err := s.repo.GetTemp(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get temp documents. error: %w", err)
	}
	return data, nil
}

func (s *DocumentService) GetByInstrument(ctx context.Context, req *models.GetDocumentDTO) ([]*models.Document, error) {
	data, err := s.repo.GetByInstrument(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents by group. error: %w", err)
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
	docs := make([]*models.Document, 0, len(dto.Files))
	// Слайс для отслеживания созданных файлов на случай отката
	createdFiles := make([]string, 0, len(dto.Files))

	// Определяем общую часть пути заранее
	baseTempPath := filepath.Join(s.path, "temp", dto.UserId, dto.Group, dto.InstrumentId)

	for _, fh := range dto.Files {
		// 1. Проверка по расширению
		ext := strings.ToLower(filepath.Ext(fh.Filename))
		if !allowedExtensions[ext] {
			return nil, fmt.Errorf("file extension %s is not allowed", ext)
		}

		// 2. Проверка по MIME-типу
		contentType := fh.Header.Get("Content-Type")
		docType, ok := documentTypes[contentType]
		if !ok {
			return nil, fmt.Errorf("mime type %s is not allowed", contentType)
		}

		// 3. (Опционально) Ограничение размера, например 10MB
		const maxFileSize = 10 * 1024 * 1024
		if fh.Size > maxFileSize {
			return nil, fmt.Errorf("file %s is too large", fh.Filename)
		}

		doc := &models.Document{
			Id:           uuid.NewString(),
			Label:        fh.Filename,
			Size:         fh.Size,
			InstrumentId: dto.InstrumentId,
			UserId:       dto.UserId,
			Group:        dto.Group,
			DocumentType: docType,
		}
		doc.Path = filepath.Join(baseTempPath, doc.Id, fh.Filename)

		if err := s.SaveUploadedFile(fh, doc.Path); err != nil {
			s.cleanupFiles(createdFiles) // Удаляем то, что уже успели сохранить
			return nil, fmt.Errorf("failed to save file %s. error: %w", fh.Filename, err)
		}

		createdFiles = append(createdFiles, doc.Path)
		docs = append(docs, doc)

	}

	if err := s.repo.CreateSeveral(ctx, docs); err != nil {
		s.cleanupFiles(createdFiles) // Удаляем все файлы, если БД "легла"
		return nil, fmt.Errorf("failed to create documents. error: %w", err)
	}
	return docs, nil
}
func (s *DocumentService) cleanupFiles(paths []string) {
	for _, p := range paths {
		if err := os.RemoveAll(filepath.Dir(p)); err != nil {
			logger.Error("failed to cleanup document files", logger.ErrAttr(err))
		}
	}
}

func (s *DocumentService) ChangePath(ctx context.Context, req *models.PathParts) (err error) {
	// 1. Валидация входных данных
	if !isValidPathSegment(req.UserId) ||
		!isValidPathSegment(req.Group) ||
		(req.InstrumentId != "" && !isValidPathSegment(req.InstrumentId)) {
		return fmt.Errorf("invalid path segment in request")
	}

	// 2. Построение путей через filepath.Join
	newPath := filepath.Join(s.path, req.Group, req.InstrumentId)
	srcPath := filepath.Join(s.path, "temp", req.UserId, req.Group, req.InstrumentId)
	if req.IdWasEmpty {
		srcPath = filepath.Join(s.path, "temp", req.UserId, req.Group)
	}

	// 3. Проверка существования источника
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return nil
		// return fmt.Errorf("source path does not exist: %s", srcPath)
	}
	// 4. Создание целевой директории
	if err := os.MkdirAll(filepath.Dir(newPath), 0750); err != nil {
		return fmt.Errorf("failed to create target dir: %w", err)
	}
	// 5. Атомарное перемещение
	if err := os.Rename(srcPath, newPath); err != nil {
		return fmt.Errorf("failed to move files from %s to %s: %w", srcPath, newPath, err)
	}

	// 6. Обновление БД (теперь после файловой операции!)
	count, err := s.repo.UpdatePath(ctx, req)
	if err != nil {
		rollback(newPath, srcPath)
		logger.Error("DB update failed after file move", logger.StringAttr("newPath", newPath), logger.ErrAttr(err))
		return fmt.Errorf("db update failed after file move: %w", err)
	}

	if count == 0 {
		logger.Info("No documents updated in DB despite successful file move", logger.StringAttr("userId", req.UserId), logger.StringAttr("group", req.Group))
	}
	return nil
}
func isValidPathSegment(s string) bool {
	if s == "" || strings.ContainsAny(s, `/\\.:*?"<>|`) || strings.Contains(s, "..") {
		return false
	}
	return true
}
func rollback(newPath, srcPath string) {
	rollbackErr := os.Rename(newPath, srcPath)
	if rollbackErr != nil {
		logger.Error("CRITICAL: rollback failed",
			logger.StringAttr("from", newPath),
			logger.StringAttr("to", srcPath),
			logger.ErrAttr(rollbackErr))
	}
}

func (s *DocumentService) Delete(ctx context.Context, dto *models.DeleteDocumentDTO) error {
	paths := []string{s.path}
	if !dto.IsTemp {
		paths = append(paths, dto.Group, dto.InstrumentId)
	} else {
		paths = append(paths, "temp", dto.UserId, dto.Group, dto.InstrumentId)
	}
	paths = append(paths, dto.Id, dto.Filename)

	dst := filepath.Join(paths...)

	if err := os.Remove(dst); err != nil && !strings.Contains(err.Error(), "no such file") {
		return fmt.Errorf("failed to delete file. error: %w", err)
	}

	if err := s.repo.Delete(ctx, dto); err != nil {
		return fmt.Errorf("failed to delete document by id. error: %w", err)
	}
	return nil
}

func (s *DocumentService) DeleteByInstrumentId(ctx context.Context, instrumentId string) error {
	dst := filepath.Join(s.path, instrumentId)

	if err := os.RemoveAll(dst); err != nil && !strings.Contains(err.Error(), "no such file") {
		return fmt.Errorf("failed to delete folder with files. error: %w", err)
	}
	return nil
}

func (s *DocumentService) RemoveEmptyFolders(ctx context.Context) error {
	logger.Info("Removing empty folders")
	// root, err := s.buildTreeFromDir("/files/si")
	// if err != nil {
	// 	return fmt.Errorf("failed to build tree from dir. error: %w", err)
	// }
	// if err := s.recursiveEmptyDelete(root); err != nil {
	// 	return fmt.Errorf("failed to remove empty folders. error: %w", err)
	// }

	count, err := s.pruneEmptyDirectories("files")
	if err != nil {
		return fmt.Errorf("failed to remove empty folders. error: %w", err)
	}
	logger.Info(fmt.Sprintf("Removed %d empty folders", count))
	return nil
}

func (s *DocumentService) pruneEmptyDirectories(osDirname string) (int, error) {
	var count int

	err := godirwalk.Walk(osDirname, &godirwalk.Options{
		Unsorted: true,
		Callback: func(_ string, _ *godirwalk.Dirent) error {
			// no-op while diving in; all the fun happens in PostChildrenCallback
			return nil
		},
		PostChildrenCallback: func(osPathname string, _ *godirwalk.Dirent) error {
			s, err := godirwalk.NewScanner(osPathname)
			if err != nil {
				return err
			}

			// Attempt to read only the first directory entry. Remember that
			// Scan skips both "." and ".." entries.
			hasAtLeastOneChild := s.Scan()

			// If error reading from directory, wrap up and return.
			if err := s.Err(); err != nil {
				return err
			}

			if hasAtLeastOneChild {
				return nil // do not remove directory with at least one child
			}
			if osPathname == osDirname {
				return nil // do not remove directory that was provided top-level directory
			}

			err = os.Remove(osPathname)
			if err == nil {
				count++
			}
			return err
		},
	})

	return count, err
}

// func (s *DocumentService) ifDir(path string) (bool, error) {
// 	file, err := os.Open(path)
// 	if err != nil {
// 		return false, fmt.Errorf("failed to open file. error: %w", err)
// 	}
// 	defer file.Close()
// 	info, err := file.Stat()
// 	if err != nil {
// 		return false, fmt.Errorf("failed to get file info. error: %w", err)
// 	}
// 	if info.IsDir() {
// 		return true, nil
// 	}
// 	return false, nil
// }

// type Node struct {
// 	Id       string
// 	Children []*Node
// }

// func (s *DocumentService) buildTreeFromDir(baseDir string) (*Node, error) {
// 	_, err := os.ReadDir(baseDir)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read dir. error: %w", err)
// 	}

// 	root := &Node{
// 		Id: baseDir,
// 	}

// 	queue := make(chan *Node, 500) // Consider that there can not be any dir with > 500 depth
// 	queue <- root
// 	for {
// 		if len(queue) == 0 {
// 			break
// 		}

// 		data, ok := <-queue
// 		if ok {
// 			// Iterate all the contents in the dir
// 			curDir := (*data).Id
// 			ok, err := s.ifDir(curDir)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to check if dir. error: %w", err)
// 			}

// 			if ok {
// 				contents, err := os.ReadDir(curDir)
// 				if err != nil {
// 					return nil, fmt.Errorf("failed to read dir. error: %w", err)
// 				}

// 				data.Children = make([]*Node, len(contents))
// 				for i, content := range contents {
// 					node := new(Node)
// 					node.Id = filepath.Join(curDir, content.Name())
// 					data.Children[i] = node
// 					if content.IsDir() {
// 						queue <- node
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return root, nil
// }
// func (s *DocumentService) recursiveEmptyDelete(root *Node) error {
// 	// If the current root is not pointing to any dir
// 	if root == nil {
// 		return nil
// 	}
// 	for _, each := range root.Children {
// 		if err := s.recursiveEmptyDelete(each); err != nil {
// 			return err
// 		}
// 	}
// 	ok, err := s.ifDir(root.Id)
// 	if err != nil {
// 		return fmt.Errorf("failed to check if dir. error: %w", err)
// 	}

// 	if !ok {
// 		return nil
// 	} else if content, _ := os.ReadDir(root.Id); len(content) != 0 {
// 		return nil
// 	}
// 	if err := os.Remove(root.Id); err != nil {
// 		return fmt.Errorf("failed to delete dir. error: %w", err)
// 	}
// 	return nil
// }
