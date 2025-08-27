package si

import (
	"errors"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/documents"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/instruments"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/locations"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/verifications"
	"github.com/Alexander272/mersi/backend/internal/utils"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
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

	si := api.Group("si", middleware.VerifyToken, middleware.CheckPermissions(constants.SI, constants.Read))
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
}

func (h *Handler) get(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Сессия не найдена")
		return
	}

	params := utils.GetFilterParams(c)
	params.SectionId = section

	data, err := h.service.Get(c, params)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), params)
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
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	req := &models.GetSiByIdDTO{Id: id}

	data, err := h.service.GetById(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getSent(c *gin.Context) {
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Сессия не найдена")
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

	data, err := h.service.GetSent(c, params)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), params)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)
	dto.Instrument.UserId = user.ID
	if dto.Location != nil {
		dto.Location.UserId = user.ID
	}

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("СИ сохранено",
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
		// logger.AnyAttr("location-dto", dto.Location),
	)
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные о си успешно сохранены"})
}

func (h *Handler) update(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("СИ обновлено",
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
	)
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о си обновлены"})
}

func (h *Handler) changePosition(c *gin.Context) {
	dto := &models.ChangePositionDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	if err := h.service.ChangePosition(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Номер позиции изменен"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	dto := &models.DeleteSiDTO{Id: id}

	if err := h.service.Delete(c, dto); err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Не удалось удалить инструмент. Нельзя удалить инструмент находящийся у сотрудника.")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	var user models.User
	u, exists := c.Get(constants.CtxUser)
	if exists {
		user = u.(models.User)
	}

	logger.Info("Инструмент отмечен как удаленный",
		logger.StringAttr("instrument_id", id),
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные об инструменте успешно удалены"})
}
