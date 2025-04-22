package columns

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Columns
}

func NewHandler(service services.Columns) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Columns, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	columns := api.Group("/columns", middleware.CheckPermissions(constants.Columns, constants.Read))
	{
		columns.GET("", handler.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	if uuid.Validate(section) != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Отправлены некорректные данные")
		return
	}
	dto := &models.GetColumnsDTO{SectionID: section}

	data, err := h.service.Get(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}
