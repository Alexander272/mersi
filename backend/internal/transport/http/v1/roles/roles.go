package roles

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type RoleHandlers struct {
	service services.Role
}

func NewRoleHandlers(service services.Role) *RoleHandlers {
	return &RoleHandlers{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Role, middleware *middleware.Middleware) {
	handlers := NewRoleHandlers(service)

	roles := api.Group("/roles", middleware.CheckPermissions(constants.Roles, constants.Read))
	{
		roles.GET("", handlers.getAll)

		write := roles.Group("", middleware.CheckPermissions(constants.Roles, constants.Write))
		{
			write.GET("/:name", handlers.get)
			write.POST("", handlers.create)
			write.PUT("/:id", handlers.update)
			write.DELETE("/:id", handlers.delete)
		}
	}
}

func (h *RoleHandlers) getAll(c *gin.Context) {
	roles, err := h.service.GetAll(c, &models.GetRolesDTO{})
	if err != nil {
		response.SendError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: roles})
}

func (h *RoleHandlers) get(c *gin.Context) {
	roleName := c.Param("name")

	role, err := h.service.Get(c, roleName)
	if err != nil {
		response.SendError(c, err, roleName)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: role})
}

func (h *RoleHandlers) create(c *gin.Context) {
	dto := &models.RoleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Роль создана"})
}

func (h *RoleHandlers) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	dto := &models.RoleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	dto.ID = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusOK, response.IdResponse{Message: "Роль обновлена"})
}

func (h *RoleHandlers) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	if err := h.service.Delete(c, id); err != nil {
		response.SendError(c, err, id)
		return
	}

	c.JSON(http.StatusNoContent, response.IdResponse{})
}
