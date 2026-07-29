package sorting

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Sorting
}

func NewHandler(service services.Sorting) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Sorting, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	sorting := api.Group("sorting", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		sorting.GET("", handler.get)
		sorting.POST("", handler.create)
		sorting.POST("/several", handler.createSeveral)
		sorting.POST("/change", handler.change)
		sorting.PUT("/:name", handler.update)
		sorting.DELETE("/:name", handler.delete)
	}
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	req := &models.GetSortingDTO{
		SectionId: section,
		UserId:    user.ID,
	}

	data, err := h.service.Get(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.SortingDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	dto.UserId = user.ID

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры успешно сохранены"})
}

func (h *Handler) createSeveral(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := []*models.SortingDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	for i := range dto {
		dto[i].UserId = user.ID
	}

	if err := h.service.CreateSeveral(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры успешно сохранены"})
}

func (h *Handler) update(c *gin.Context) {
	dto := &models.SortingDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	dto.UserId = user.ID

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры успешно сохранены"})
}

func (h *Handler) change(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := []*models.SortingDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	for i := range dto {
		dto[i].UserId = user.ID
		dto[i].SectionId = section
	}

	if len(dto) == 0 {
		if err := h.service.DeleteAll(c, &models.DeleteSortingDTO{UserId: user.ID, SectionId: section}); err != nil {
			response.SendError(c, err, dto)
			return
		}
	} else {
		if err := h.service.Change(c, dto); err != nil {
			response.SendError(c, err, dto)
			return
		}
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтры успешно сохранены"})
}

func (h *Handler) delete(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	dto := &models.DeleteSortingDTO{
		Name:      name,
		SectionId: section,
		UserId:    user.ID,
	}

	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Фильтр сброшен"})
}
