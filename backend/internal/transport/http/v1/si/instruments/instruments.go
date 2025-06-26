package instruments

import (
	"errors"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
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
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	req := &models.GetInstrumentByIdDTO{Id: id}

	data, err := h.service.GetById(c, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), err.Error())
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getUnique(c *gin.Context) {
	field := c.Param("field")
	if field == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "field is empty", "Отправлены некорректные данные")
		return
	}
	req := &models.GetUniqueDTO{Field: field}

	data, err := h.service.GetUniqueData(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.InstrumentDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
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

	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
		return
	}
	user := u.(models.User)
	dto.UserId = user.ID

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	logger.Info("Создан инструмент",
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("instrument_id", dto.Id),
		logger.StringAttr("instrument_name", dto.Name),
		logger.AnyAttr("instrument", dto),
	)

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные об инструменте успешно добавлены"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id параметр не валиден")
		return
	}

	dto := &models.InstrumentDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}
	dto.Id = id

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	var user models.User
	u, exists := c.Get(constants.CtxUser)
	if exists {
		user = u.(models.User)
	}

	logger.Info("Инструмент обновлен",
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("instrument_id", dto.Id),
		logger.StringAttr("instrument_name", dto.Name),
		logger.AnyAttr("instrument", dto),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные об инструменте успешно обновлены"})
}
