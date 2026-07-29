package context_menu

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
	service services.ContextMenu
}

func NewHandler(service services.ContextMenu) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.ContextMenu, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	context := api.Group("context-menu", middleware.CheckPermissions(constants.ContextMenu, constants.Read))
	{
		context.GET("", handler.get)

		write := context.Group("", middleware.CheckPermissions(constants.ContextMenu, constants.Write))
		{
			write.POST("", handler.create)
			write.POST("/several", handler.createSeveral)
			write.PUT("/:id", handler.update)
			write.PUT("/several", handler.updateSeveral)
			write.DELETE("/:id", handler.delete)
			write.DELETE("/several", handler.deleteSeveral)
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	isFull := c.Query("isFull")
	realm := c.GetHeader("realm")

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

	req := &models.GetContextMenuDTO{
		SectionId: section,
		UserId:    user.ID,
		RealmId:   realm,
		IsFull:    isFull == "true",
	}
	data, err := h.service.Get(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.ContextMenuDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Id: dto.Id, Message: "Пункт меню создан"})
}

func (h *Handler) createSeveral(c *gin.Context) {
	dto := []*models.ContextMenuDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.CreateSeveral(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Пункты меню созданы"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	dto := &models.ContextMenuDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	dto.Id = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пункт меню обновлен"})
}

func (h *Handler) updateSeveral(c *gin.Context) {
	var dto []*models.ContextMenuDTO
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.UpdateSeveral(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	dto := &models.DeleteContextMenuDTO{Id: id}

	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusNoContent, response.IdResponse{})
}

func (h *Handler) deleteSeveral(c *gin.Context) {
	var dto []string
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.DeleteSeveral(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные удалены"})
}
