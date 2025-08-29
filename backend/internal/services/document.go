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

		// paths := []string{s.path}
		// if dto.InstrumentId != "" {
		// 	paths = append(paths, dto.Group, dto.InstrumentId)
		// } else {
		// 	paths = append(paths, "temp", dto.UserId, dto.Group, dto.InstrumentId)
		// }
		// paths = append(paths, doc.Id, fh.Filename)
		paths := []string{s.path, "temp", dto.UserId, dto.Group, dto.InstrumentId, doc.Id, fh.Filename}

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
		// думаю тут еще надо id пользователя использовать (чтобы во время одновременного создания ничего лишнего не попало)
		// по хорошему это еще надо синхронизировать между устройствами
		// еще можно добавить какую-нибудь группировку чтобы файлы из разных мест не пересекались или она мне не нужна
		srcPath := path.Join(s.path, "temp", req.UserId, req.Group, req.InstrumentId)
		if req.IdWasEmpty {
			srcPath = path.Join(s.path, "temp", req.UserId, req.Group)
		}

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
	if !dto.IsTemp {
		paths = append(paths, dto.Group, dto.InstrumentId)
	} else {
		paths = append(paths, "temp", dto.UserId, dto.Group, dto.InstrumentId)
	}
	paths = append(paths, dto.Id, dto.Filename)

	dst := path.Join(paths...)

	if err := os.Remove(dst); err != nil && !strings.Contains(err.Error(), "no such file") {
		return fmt.Errorf("failed to delete file. error: %w", err)
	}

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
