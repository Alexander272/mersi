package permissions

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Permission
}

func NewHandler(service services.Permission) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Permission, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	permissions := api.Group("permissions", middleware.CheckPermissions(constants.Realms, constants.Write))
	{
		permissions.POST("/reload", handler.reload)
	}
}

func (h *Handler) reload(c *gin.Context) {
	if err := h.service.ReloadPolicies(c.Request.Context()); err != nil {
		response.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные обновлены"})
}
