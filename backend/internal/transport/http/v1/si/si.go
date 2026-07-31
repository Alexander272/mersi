package si

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/utils"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/documents"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/instruments"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/locations"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/receiving"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/verifications"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.SI
}

func NewHandler(service services.SI) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, services *services.Services, middleware *middleware.Middleware) {
	handler := NewHandler(services.SI)

	si := api.Group("si", middleware.VerifyToken, middleware.DepartmentAccess, middleware.CheckPermissions(constants.SI, constants.Read))
	{
		si.GET("", handler.get)
		si.GET("/:id", handler.getById)
		si.GET("/sent", handler.getSent)

		write := si.Group("", middleware.CheckPermissions(constants.SI, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.PUT("/position", handler.changePosition)
			write.DELETE("/:id", handler.delete)
		}
	}

	instruments.Register(si, services.Instrument, middleware)
	documents.Register(si, services.Document, middleware)
	verifications.Register(si, services, middleware)
	locations.Register(api, services.Location, middleware)
	receiving.Register(api, services.Receiving, middleware)
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	params := utils.GetFilterParams(c)
	if params == nil {
		response.SendError(c, models.ErrNotValid)
		return
	}
	params.SectionId = section

	data, err := h.service.Get(c, params)
	if err != nil {
		response.SendError(c, err, params)
		return
	}
	total := 0
	if len(data) > 0 {
		total = data[0].Total
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: total})
}

func (h *Handler) getById(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	req := &models.GetSiByIdDTO{Id: id}

	data, err := h.service.GetById(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getSent(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	params := &models.GetSiDTO{
		SectionId: section,
	}

	filters := c.QueryMap("filters")
	for k, v := range filters {
		valueMap := c.QueryMap(k)

		values := []*models.FilterValue{}
		for key, value := range valueMap {
			values = append(values, &models.FilterValue{
				CompareType: key,
				Value:       value,
			})
		}

		if k == "place" {
			k = "department"
		}

		f := &models.Filter{
			Field:     k,
			FieldType: v,
			Values:    values,
		}
		params.Filters = append(params.Filters, f)
	}

	hasWritePermission := true
	if wp, exists := c.Get(constants.CtxHasWritePermission); exists {
		if wp, ok := wp.(bool); ok {
			hasWritePermission = wp
		}
	}
	if !hasWritePermission {
		if deptIDs, exists := c.Get(constants.CtxDepartmentAccess); exists {
			if ids, ok := deptIDs.([]string); ok && len(ids) > 0 {
				params.DepartmentAccess = ids
				params.Filters = append(params.Filters, &models.Filter{
					Field: "department", FieldType: "list",
					Values: []*models.FilterValue{{
						CompareType: "in",
						Value:       strings.Join(ids, ","),
					}},
				})
			}
		}
		params.Filters = append(params.Filters, &models.Filter{
				Field: "status", Values: []*models.FilterValue{{CompareType: "nlike", Value: "reserve"}},
			})
	}

	data, err := h.service.GetSent(c, params)
	if err != nil {
		response.SendError(c, err, params)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}

	dto.Instrument.UserId = actor.ID
	dto.Instrument.Actor = actor

	if dto.Verification != nil {
		dto.Verification.Actor = actor
	}

	if dto.Location != nil {
		dto.Location.UserId = actor.ID
		dto.Location.Actor = actor
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("СИ сохранено",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
		// logger.AnyAttr("location-dto", dto.Location),
	)
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные о си успешно сохранены"})
}

func (h *Handler) update(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}

	dto.Instrument.UserId = actor.ID
	dto.Instrument.Actor = actor
	if dto.Verification != nil {
		dto.Verification.Actor = actor
	}
	if dto.Location != nil {
		dto.Location.Actor = actor
	}

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("СИ обновлено",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
	)
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о си обновлены"})
}

func (h *Handler) changePosition(c *gin.Context) {
	dto := &models.ChangePositionDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.ChangePosition(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Номер позиции изменен"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	dto := &models.DeleteSiDTO{Id: id}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor

	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Инструмент отмечен как удаленный",
		logger.StringAttr("instrument_id", id),
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные об инструменте успешно удалены"})
}
