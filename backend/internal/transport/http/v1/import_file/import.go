package import_file

import (
	"fmt"
	"net/http"
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
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	form, err := c.MultipartForm()
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNoFileInRequest, err))
		return
	}

	if len(form.Value["realm"]) == 0 || len(form.Value["section"]) == 0 || len(form.Value["bid_type"]) == 0 {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	realm := form.Value["realm"][0]
	err = uuid.Validate(realm)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	section := form.Value["section"][0]
	err = uuid.Validate(section)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	bidType := form.Value["bid_type"][0]

	files := form.File["files"]
	if len(files) == 0 {
		response.SendError(c, fmt.Errorf("%w: no files", models.ErrNotValid))
		return
	}

	sheetType := "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	if files[0].Header.Get("Content-Type") != sheetType && !strings.Contains(files[0].Filename, "xls") {
		response.SendError(c, fmt.Errorf("%w: invalid type file", models.ErrNotValid))
		return
	}

	dto := &models.ImportDTO{RealmId: realm, SectionId: section, BidType: bidType, UserId: user.ID, File: files[0]}
	if err := h.service.Load(c, dto); err != nil {
		response.SendError(c, err, dto)
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
