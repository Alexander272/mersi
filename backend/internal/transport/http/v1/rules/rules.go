package rules

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
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

	rules := api.Group("rules", middleware.CheckPermissions(constants.Roles, constants.Write))
	{
		rules.GET("", handler.get)
	}
}

func (h *Handler) get(c *gin.Context) {
	data, err := h.service.GetAll(c, &models.GetRuleItemsDTO{OnlyShow: true})
	if err != nil {
		response.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}
