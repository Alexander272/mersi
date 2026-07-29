package instruments

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/utils"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Instrument
}

func NewHandler(service services.Instrument) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Instrument, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	instruments := api.Group("/instruments", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		instruments.GET("/:id", handler.get)
		instruments.GET("/unique/:field", handler.getUnique)

		write := instruments.Group("", middleware.CheckPermissions(constants.SI, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.POST("/status", handler.changeStatus)
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	req := &models.GetInstrumentByIdDTO{Id: id}

	data, err := h.service.GetById(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getUnique(c *gin.Context) {
	field := c.Param("field")
	if field == "" {
		response.SendError(c, models.ErrNotValid)
		return
	}
	section := c.Query("section")
	err := uuid.Validate(section)
	if err != nil {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	req := &models.GetUniqueDTO{Field: field, SectionId: section}

	data, err := h.service.GetUniqueData(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.InstrumentDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	//TODO возможно нужно сделать подобное, но только с section
	// realm := c.GetHeader("realm")
	// err := uuid.Validate(realm)
	// if err != nil {
	// 	response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Сессия не найдена")
	// 	return
	// }
	// dto.RealmId = realm

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor
	dto.UserId = actor.ID

	if err := h.service.Create(c, nil, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	logger.Info("Создан инструмент",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
		logger.StringAttr("instrument_id", dto.Id),
		logger.StringAttr("instrument_name", dto.Name),
		logger.AnyAttr("instrument", dto),
	)

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные об инструменте успешно добавлены"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.InstrumentDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	dto.Id = id

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor
	dto.UserId = actor.ID

	if err := h.service.Update(c, nil, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Инструмент обновлен",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
		logger.StringAttr("instrument_id", dto.Id),
		logger.StringAttr("instrument_name", dto.Name),
		logger.AnyAttr("instrument", dto),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные об инструменте успешно обновлены"})
}

func (h *Handler) changeStatus(c *gin.Context) {
	dto := []*models.UpdateStatus{}
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.ChangeSeveralStatuses(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user, ok := u.(models.User)
	if !ok {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}

	logger.Info("Статус инструментов обновлен",
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные об инструментах успешно обновлены"})
}
