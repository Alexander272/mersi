package documents

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
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
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Группа файлов не задана")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	req := &models.GetTempDocumentDTO{
		Group:  group,
		UserId: user.ID,
	}

	data, err := h.service.GetTemp(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) download(c *gin.Context) {
	path := c.Query("path")

	fileStat, err := os.Stat(path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), "Файл не найден")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Файл не найден")
		error_bot.Send(c, err.Error(), path)
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Length", fmt.Sprintf("%d", fileStat.Size()))
	c.Header("Content-Disposition", "attachment; filename="+fileStat.Name())
	c.File(path)
}

func (h *Handler) upload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Не удалось получить файлы")
		return
	}

	instrumentId := form.Value["instrumentId"][0]
	group := form.Value["group"][0]

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	files := form.File["files"]
	if len(files) == 0 {
		response.NewErrorResponse(c, http.StatusNoContent, "no content", "Нет файлов для загрузки")
		return
	}

	dto := &models.DocumentsDTO{InstrumentId: instrumentId, Group: group, UserId: user.ID, Files: files}
	res, err := h.service.Upload(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
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
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id документа не задан")
		return
	}
	instrumentId := c.Query("instrumentId")
	filename := c.Query("filename")
	if filename == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Имя файла не задано")
		return
	}
	group := c.Query("group")

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	req := &models.DeleteDocumentDTO{
		Id:           id,
		InstrumentId: instrumentId,
		Filename:     filename,
		Group:        group,
		UserId:       user.ID,
	}

	if err := h.service.Delete(c, req); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}

	logger.Info("Файл удален",
		logger.AnyAttr("dto", req),
	)
	c.JSON(http.StatusOK, response.IdResponse{Message: "Файл удален"})
}
