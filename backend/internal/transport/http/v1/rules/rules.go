package rules

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/access"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.RuleItem
}

func NewHandler(service services.RuleItem) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.RuleItem, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	rules := api.Group("rules", middleware.CheckPermissions(access.Reg.R(access.ResourceRoles).Write()))
	{
		rules.GET("", handler.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	data, err := h.service.GetAll(c, &models.GetRuleItemsDTO{OnlyShow: true})
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}
