package export

import (
	"bytes"
	"fmt"
	"time"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/utils"
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

	export := api.Group("export", middleware.DepartmentAccess, middleware.CheckPermissions(constants.SI, constants.Read))
	{
		export.GET("", handler.export)
		export.GET("/schedule", handler.makeScheduler)
		export.GET("/accounting", handler.makeAccountingLog)
	}
}

func (h *Handler) export(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	params := utils.GetFilterParams(c)

	params.SectionId = section
	params.Page = &models.Page{
		Limit: 999999,
	}

	buffer, err := h.service.Export(c, params)
	if err != nil {
		response.SendError(c, err, params)
		return
	}

	h.sendExcel(c, buffer, "Список инструментов")
}

func (h *Handler) makeScheduler(c *gin.Context) {
	req, err := h.parsePeriod(c)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	h.enrichPeriodWithDeptAccess(c, req)

	buffer, err := h.service.MakeScheduler(c, req)
	if err != nil {
		h.handleServiceError(c, err, req)
		return
	}

	h.sendExcel(c, buffer, "График поверки")
}

func (h *Handler) makeAccountingLog(c *gin.Context) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	req := &models.Period{
		SectionId: section,
	}
	h.enrichPeriodWithDeptAccess(c, req)

	buffer, err := h.service.MakeAccountingLog(c, req)
	if err != nil {
		h.handleServiceError(c, err, req)
		return
	}

	h.sendExcel(c, buffer, "Журнал учета средств измерения")
}

func (h *Handler) enrichPeriodWithDeptAccess(c *gin.Context, req *models.Period) {
	hasWritePermission := true
	if wp, exists := c.Get(constants.CtxHasWritePermission); exists {
		hasWritePermission = wp.(bool)
	}
	if hasWritePermission {
		return
	}
	if deptIDs, exists := c.Get(constants.CtxDepartmentAccess); exists {
		ids := deptIDs.([]string)
		if len(ids) > 0 {
			req.DepartmentAccess = ids
		}
	}
}

func (h *Handler) handleServiceError(c *gin.Context, err error, req interface{}) {
	response.SendError(c, err, req)
}

func (h *Handler) parsePeriod(c *gin.Context) (*models.Period, error) {
	section := c.Query("section")
	if err := uuid.Validate(section); err != nil {
		return nil, fmt.Errorf("invalid section UUID")
	}

	period := c.QueryMap("period")
	start, errStart := time.Parse(time.RFC3339, period["gte"])
	end, errEnd := time.Parse(time.RFC3339, period["lte"])

	if errStart != nil || errEnd != nil {
		return nil, fmt.Errorf("invalid or missing period dates")
	}

	return &models.Period{
		StartAt:         start,
		FinishAt:        end,
		SectionId:       section,
		ChannelIsOption: true,
	}, nil
}

func (h *Handler) sendExcel(c *gin.Context, buffer *bytes.Buffer, fileNamePrefix string) {
	fileName := fmt.Sprintf("%s от %s.xlsx", fileNamePrefix, monday.Format(time.Now(), "Mon 2 Jan 2006", monday.LocaleRuRU))

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Length", fmt.Sprintf("%d", buffer.Len()))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache")

	c.Writer.Write(buffer.Bytes())
}
