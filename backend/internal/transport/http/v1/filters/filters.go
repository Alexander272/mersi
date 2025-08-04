package filters

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
	service services.Filters
}

func NewHandler(service services.Filters) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Filters, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	filters := api.Group("filters", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		filters.GET("", handler.get)
		filters.POST("", handler.create)
		filters.POST("/change", handler.change)
		filters.DELETE("/:id", handler.delete)
	}
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	req := &models.GetSavedFiltersDTO{
		SectionId: section,
		UserId:    user.ID,
	}

	data, err := h.service.Get(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) create(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}

	dto := []*models.SavedFilterDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	for i := range dto {
		dto[i].UserId = user.ID
		dto[i].SectionId = section
	}

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры успешно сохранены"})
}

func (h *Handler) change(c *gin.Context) {
	// section := c.Query("section")
	// if err := uuid.Validate(section); err != nil {
	// 	response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
	// 	return
	// }

	dto := &models.ChangeFillersDTO{}
	// dto := []*models.SavedFilterDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	for i := range dto.Filters {
		dto.Filters[i].UserId = user.ID
		dto.Filters[i].SectionId = dto.SectionId
	}

	if len(dto.Filters) == 0 {
		if err := h.service.Delete(c, &models.DeleteSavedFiltersDTO{UserId: user.ID, SectionId: dto.SectionId}); err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
			error_bot.Send(c, err.Error(), dto)
			return
		}
	} else {
		if err := h.service.Change(c, dto.Filters); err != nil {
			response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
			error_bot.Send(c, err.Error(), dto)
			return
		}
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры успешно сохранены"})
}

func (h *Handler) delete(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	dto := &models.DeleteSavedFiltersDTO{
		SectionId: section,
		UserId:    user.ID,
	}

	if err := h.service.Delete(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры сброшены"})
}
