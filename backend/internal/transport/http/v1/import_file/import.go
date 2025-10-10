package import_file

import (
	"net/http"
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
	service services.ImportFile
}

func NewHandler(service services.ImportFile) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.ImportFile, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	importFile := api.Group("import", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		importFile.POST("", handler.load)
	}
}

func (h *Handler) load(c *gin.Context) {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	form, err := c.MultipartForm()
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Не удалось получить файлы")
		return
	}

	if len(form.Value["realm"]) == 0 || len(form.Value["section"]) == 0 || len(form.Value["bid_type"]) == 0 {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Недостаточно параметров")
		return
	}

	realm := form.Value["realm"][0]
	err = uuid.Validate(realm)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "invalid id param")
		return
	}

	section := form.Value["section"][0]
	err = uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "invalid id param")
		return
	}

	bidType := form.Value["bid_type"][0]

	files := form.File["files"]
	if len(files) == 0 {
		response.NewErrorResponse(c, http.StatusNoContent, "no content", "Нет файлов для загрузки")
		return
	}

	sheetType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if files[0].Header.Get("Content-Type") != sheetType && !strings.Contains(files[0].Filename, "xls") {
		response.NewErrorResponse(c, http.StatusInternalServerError, "invalid type file", "Данный формат не поддерживается")
		return
	}

	dto := &models.ImportDTO{RealmId: realm, SectionId: section, BidType: bidType, UserId: user.ID, File: files[0]}
	if err := h.service.Load(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Файлы загружены",
		logger.StringAttr("realm_id", dto.RealmId),
		logger.StringAttr("section_id", dto.SectionId),
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
	)
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Файлы загружены"})
}
