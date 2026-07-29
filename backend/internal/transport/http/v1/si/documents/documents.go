package documents

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Document
}

func NewHandler(service services.Document) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Document, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	docs := api.Group("documents", middleware.CheckPermissions(constants.Documents, constants.Read))
	{
		docs.GET("", handler.download)
		docs.GET("/temp/:group", handler.getTemp)
		docs.GET("/list/:group", handler.getList)

		write := docs.Group("", middleware.CheckPermissions(constants.Documents, constants.Write))
		{
			write.POST("", handler.upload)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) getTemp(c *gin.Context) {
	group := c.Param("group")
	if group == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	instrument := c.Query("instrument")

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	req := &models.GetDocumentDTO{
		Group:        group,
		UserId:       user.ID,
		InstrumentId: instrument,
	}

	data, err := h.service.GetTemp(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getList(c *gin.Context) {
	group := c.Param("group")
	if group == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	req := &models.GetDocumentDTO{
		Group:        group,
		InstrumentId: instrument,
		UserId:       user.ID,
	}

	data, err := h.service.GetByInstrument(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) download(c *gin.Context) {
	rawPath := c.Query("path")
	if rawPath == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	cleaned := filepath.Clean(rawPath)
	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	baseDir, err := filepath.Abs("files")
	if err != nil {
		response.SendError(c, fmt.Errorf("failed to get base dir: %w", err))
		return
	}

	if !strings.HasPrefix(absPath, baseDir+string(os.PathSeparator)) && absPath != baseDir {
		response.SendError(c, models.ErrForbidden)
		return
	}

	fileStat, err := os.Stat(absPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			response.SendError(c, models.ErrFileNotFound)
			return
		}
		response.SendError(c, err, rawPath)
		return
	}

	if fileStat.IsDir() {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprintf("%d", fileStat.Size()))
	c.Header("Content-Disposition", "attachment; filename="+fileStat.Name())
	c.File(absPath)
}

func (h *Handler) upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNoFileInRequest, err))
		return
	}

	instrumentId := form.Value["instrumentId"][0]
	group := form.Value["group"][0]

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	files := form.File["files"]
	if len(files) == 0 {
		response.SendError(c, fmt.Errorf("%w: no files", models.ErrNotValid))
		return
	}

	dto := &models.DocumentsDTO{InstrumentId: instrumentId, Group: group, UserId: user.ID, Files: files}
	res, err := h.service.Upload(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Файлы загружены",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("user_id", dto.UserId),
		logger.StringAttr("group", dto.Group),
	)
	c.JSON(http.StatusCreated, response.DataResponse{Data: res})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	instrumentId := c.Query("instrumentId")
	filename := c.Query("filename")
	if filename == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	group := c.Query("group")
	isTemp := c.Query("isTemp")

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	req := &models.DeleteDocumentDTO{
		Id:           id,
		InstrumentId: instrumentId,
		Filename:     filename,
		Group:        group,
		UserId:       user.ID,
		IsTemp:       isTemp == "true",
	}

	if err := h.service.Delete(c, req); err != nil {
		response.SendError(c, err, req)
		return
	}

	logger.Info("Файл удален",
		logger.AnyAttr("dto", req),
	)
	c.JSON(http.StatusOK, response.IdResponse{Message: "Файл удален"})
}
