package department_accesses

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
	service services.DepartmentAccess
}

func NewHandler(service services.DepartmentAccess) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.DepartmentAccess, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	accesses := api.Group("/department-accesses", middleware.CheckPermissions(constants.Department, constants.Read))
	{
		accesses.GET(":id", handler.get)
		accesses.GET("/user/:id", handler.getByUser)

		write := accesses.Group("", middleware.CheckPermissions(constants.Department, constants.Write))
		{
			write.POST("/replace", handler.replace)
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	department := c.Param("id")
	if err := uuid.Validate(department); err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}
	dto := &models.GetDepartmentAccessDTO{DepartmentId: department}

	data, err := h.service.Get(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getByUser(c *gin.Context) {}

func (h *Handler) replace(c *gin.Context) {
	dto := &models.ReplaceDepartmentAccessDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	if err := h.service.Replace(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные созданы"})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.DepartmentAccessDTO{}
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
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	dto := &models.DepartmentAccessDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if id != dto.Id {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	dto := &models.DeleteDepartmentAccessDTO{Id: id}
	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusNoContent, response.IdResponse{})
}
