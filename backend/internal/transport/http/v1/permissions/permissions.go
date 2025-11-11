package permissions

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
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
	if err := h.service.ReloadPolicies(); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные обновлены"})
}
