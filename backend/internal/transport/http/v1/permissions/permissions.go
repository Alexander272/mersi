package permissions

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/access"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.AccessPolices
}

func NewHandler(service services.AccessPolices) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.AccessPolices, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	permissions := api.Group("permissions", middleware.CheckPermissions(access.Reg.R(access.ResourceRealms).Write()))
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
