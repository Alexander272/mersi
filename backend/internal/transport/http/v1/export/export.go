package export

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/goodsign/monday"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Export
}

func NewHandler(service services.Export) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Export, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	export := api.Group("export", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		export.GET("/schedule", handler.makeScheduler)
	}
}

func (h *Handler) makeScheduler(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Секция не задана")
		return
	}

	period := c.QueryMap("period")
	if len(period) == 0 {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Период не задан")
		return
	}

	start, errStart := strconv.ParseInt(period["gte"], 10, 64)
	end, errEnd := strconv.ParseInt(period["lte"], 10, 64)
	if errStart != nil || errEnd != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Период не задан")
		return
	}

	req := &models.Period{
		StartAt:   start,
		FinishAt:  end,
		SectionId: section,
	}

	buffer, err := h.service.MakeScheduler(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=График поверки от "+monday.Format(time.Now(), "Mon 2 Jan 2006", monday.LocaleRuRU)+".xlsx")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/vnd.openxmlformats-package.relationships+xml")
	c.Header("Accept-Length", fmt.Sprintf("%d", buffer.Cap()))
	c.Writer.Write(buffer.Bytes())
}
