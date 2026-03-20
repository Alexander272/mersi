package export

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/utils"
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
		export.GET("", handler.export)
		export.GET("/schedule", handler.makeScheduler)
	}
}

func (h *Handler) export(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Сессия не найдена")
		return
	}

	params := utils.GetFilterParams(c)

	params.SectionId = section
	params.Page = &models.Page{
		Limit: 999999,
	}

	buffer, err := h.service.Export(c, params)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), params)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Список инструментов от "+monday.Format(time.Now(), "Mon 2 Jan 2006", monday.LocaleRuRU)+".xlsx")
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Type", "application/vnd.openxmlformats-package.relationships+xml")
	c.Header("Accept-Length", fmt.Sprintf("%d", buffer.Cap()))
	c.Writer.Write(buffer.Bytes())
}

func (h *Handler) makeScheduler(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Область не задана")
		return
	}

	period := c.QueryMap("period")
	if len(period) == 0 {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Период не задан")
		return
	}

	// start, errStart := strconv.ParseInt(period["gte"], 10, 64)
	// end, errEnd := strconv.ParseInt(period["lte"], 10, 64)
	start, errStart := time.Parse(time.RFC3339, period["gte"])
	end, errEnd := time.Parse(time.RFC3339, period["lte"])
	if errStart != nil || errEnd != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Период не задан")
		return
	}

	req := &models.Period{
		StartAt:         start,
		FinishAt:        end,
		SectionId:       section,
		ChannelIsOption: true,
	}

	buffer, err := h.service.MakeScheduler(c, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "В заданном периоде ничего не найдено")
			return
		}

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
