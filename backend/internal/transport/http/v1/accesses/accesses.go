package accesses

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
	service services.Accesses
}

func NewHandler(service services.Accesses) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Accesses, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	accesses := api.Group("/accesses", middleware.CheckPermissions(constants.Realms, constants.Read))
	{
		accesses.GET("", handler.get)
		write := accesses.Group("", middleware.CheckPermissions(constants.Realms, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	realm := c.Query("realm")
	err := uuid.Validate(realm)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	req := &models.GetAccessesDTO{RealmID: realm}
	data, err := h.service.Get(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.AccessesDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные созданы"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	dto := &models.AccessesDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	dto.ID = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	dto := &models.DeleteAccessesDTO{ID: id}
	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusNoContent, response.StatusResponse{})
}
