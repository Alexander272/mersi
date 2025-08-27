package locations

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

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
	service services.Location
}

func NewHandler(service services.Location) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Location, ware *middleware.Middleware) {
	handler := NewHandler(service)

	perm := []*middleware.Permission{
		{Section: constants.Location, Method: constants.Write},
		{Section: constants.Reserve, Method: constants.Write},
	}

	locations := api.Group("si/locations")
	{
		secure := locations.Group("", ware.VerifyToken, ware.CheckPermissions(constants.Location, constants.Read))
		{
			secure.GET("", handler.get)
			secure.GET("/last", handler.getLast)
			secure.GET("/last/several", handler.getSeveralLast)

			write := secure.Group("", ware.CheckPermissions(constants.Location, constants.Write))
			{
				write.POST("", handler.create)
				write.PUT("/:id", handler.update)
			}

			writeRes := secure.Group("", ware.CheckPermissionsArray(perm))
			{
				writeRes.POST("/several", handler.createSeveral)
				writeRes.DELETE("/:id", handler.delete)
			}
		}

		receiving := locations.Group("receiving")
		{
			receiving.POST("/dialogs", handler.receivingDialog)
			receiving.POST("/dialogs/open", handler.receivingDialogOpen)

			secure := receiving.Group("", ware.VerifyToken, ware.CheckPermissions(constants.Location, constants.Read))
			{
				secure.POST("", handler.receiving)

				write := secure.Group("", ware.CheckPermissions(constants.Location, constants.Write))
				{
					write.POST("/forced", handler.forcedReceiving)
				}
			}
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	req := &models.GetLocationDTO{InstrumentId: instrument}

	data, err := h.service.Get(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getLast(c *gin.Context) {
	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}
	req := &models.GetLocationDTO{InstrumentId: instrument}

	data, err := h.service.GetLast(c, req)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), "Данные не найдены")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getSeveralLast(c *gin.Context) {
	instruments := c.Query("instruments")
	if instruments == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "instruments is empty", "Отправлены некорректные данные")
		return
	}
	req := &models.GetSeveralLocationsDTO{InstrumentIds: strings.Split(instruments, ",")}

	data, err := h.service.GetSeveralLast(c, req)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.LocationDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	var user models.User
	u, exists := c.Get(constants.CtxUser)
	if exists {
		user = u.(models.User)
	}
	dto.UserId = user.ID

	if err := h.service.Create(c, dto); err != nil {
		if errors.Is(err, models.ErrNoChannel) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Канал для получения уведомлений не указан")
			return
		}
		if errors.Is(err, models.ErrNoResponsible) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Ответственный не указан")
			return
		}

		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Инструмент перемещен",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("status", dto.Status),
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно добавлены"})
}

func (h *Handler) createSeveral(c *gin.Context) {
	dto := []*models.LocationDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	var user models.User
	u, exists := c.Get(constants.CtxUser)
	if exists {
		user = u.(models.User)
	}
	for i := range dto {
		dto[i].UserId = user.ID
	}

	full, err := h.service.CreateSeveral(c, dto)
	if err != nil {
		if errors.Is(err, models.ErrNoResponsible) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Вы не являетесь ответственным")
			return
		}
		if errors.Is(err, models.ErrNoInstrument) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Вы не можете переместить инструмент")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Инструменты перемещены", logger.BoolAttr("full", full), logger.StringAttr("user_id", user.ID), logger.StringAttr("username", user.Name))

	message := "Данные о месте нахождения успешно добавлены"
	if !full {
		message = "Данные о месте нахождения добавлены частично"
	}

	c.JSON(http.StatusCreated, response.IdResponse{Message: message})
}

func (h *Handler) receiving(c *gin.Context) {
	dto := &models.ReceivingDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	var user models.User
	u, exists := c.Get(constants.CtxUser)
	if exists {
		user = u.(models.User)
	}
	dto.UserId = user.ID
	dto.HasConfirmed = true

	if err := h.service.ReceivingFromApp(c, dto); err != nil {
		if errors.Is(err, models.ErrNoResponsible) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Вы не являетесь ответственным")
			return
		}
		if errors.Is(err, models.ErrNoInstrument) {
			response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Вы не можете подтвердить получение инструментов")
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Получены инструменты",
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("user", user.Name),
		logger.StringAttr("status", dto.Status),
		logger.AnyAttr("instruments", dto.InstrumentIds),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) forcedReceiving(c *gin.Context) {
	dto := &models.ForcedReceiptDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	if err := h.service.ForcedReceipt(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Получены инструменты (принудительно)",
		logger.StringAttr("instrument_id", dto.InstrumentId),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) receivingDialog(c *gin.Context) {
	dto := &models.DialogResponse{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	if err := h.service.ReceivingFromChannel(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("Получены инструменты",
		logger.StringAttr("user_id", dto.UserID),
		logger.StringAttr("instrument_ids", fmt.Sprint(dto.Submission)),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) receivingDialogOpen(c *gin.Context) {
	dto := &models.PostAction{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	if err := h.service.ReceivingDialogOpen(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	logger.Info("Открытие диалога", logger.StringAttr("userId", dto.UserId))

	c.JSON(http.StatusOK, response.IdResponse{Message: "Сообщение отправлено"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}

	dto := &models.LocationDTO{}
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
	logger.Info("Место нахождения инструмента изменено",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("status", dto.Status),
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Id не валиден")
		return
	}

	dto := &models.DeleteLocationDTO{Id: id}
	if err := h.service.Delete(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	var user models.User
	u, exists := c.Get(constants.CtxUser)
	if exists {
		user = u.(models.User)
	}
	logger.Info("Место нахождения инструмента удалено",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("user_id", user.ID),
		logger.StringAttr("username", user.Name),
	)

	c.JSON(http.StatusNoContent, response.IdResponse{})
}
