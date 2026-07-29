package channel

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Channel
}

func NewHandler(service services.Channel) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Channel, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	channels := api.Group("/channels", middleware.CheckPermissions(constants.Channel, constants.Read))
	{
		channels.GET("", handler.getAll)

		write := channels.Group("", middleware.CheckPermissions(constants.Channel, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) getAll(c *gin.Context) {
	data, err := h.service.GetAll(c)
	if err != nil {
		response.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.Channel{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	logger.Info("Канал создан", logger.AnyAttr("dto", dto))

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Канал создан"})
}

func (h *Handler) update(c *gin.Context) {
	dto := &models.Channel{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	dto.ID = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	logger.Info("Канал обновлен", logger.AnyAttr("dto", dto))

	c.JSON(http.StatusOK, response.IdResponse{Message: "Канал обновлен"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	if err := h.service.Delete(c, id); err != nil {
		response.SendError(c, err, id)
		return
	}
	logger.Info("Канал удален", logger.StringAttr("id", id))

	c.JSON(http.StatusNoContent, response.StatusResponse{})
}
