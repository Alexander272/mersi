package locations

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
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	location services.Location
}

func NewHandler(location services.Location) *Handler {
	return &Handler{location: location}
}

func Register(api *gin.RouterGroup, location services.Location, ware *middleware.Middleware) {
	handler := NewHandler(location)

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
	}
}

func (h *Handler) get(c *gin.Context) {
	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	req := &models.GetLocationDTO{InstrumentId: instrument}

	data, err := h.location.Get(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handler) getLast(c *gin.Context) {
	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	req := &models.GetLocationDTO{InstrumentId: instrument}

	data, err := h.location.GetLast(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getSeveralLast(c *gin.Context) {
	instruments := c.Query("instruments")
	if instruments == "" {
		response.SendError(c, models.ErrNotValid)
		return
	}
	req := &models.GetSeveralLocationsDTO{InstrumentIds: strings.Split(instruments, ",")}

	data, err := h.location.GetSeveralLast(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.LocationDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor
	dto.UserId = actor.ID

	if err := h.location.Create(c, nil, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Инструмент перемещен",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("status", dto.Status),
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно добавлены"})
}

func (h *Handler) createSeveral(c *gin.Context) {
	dto := []*models.LocationDTO{}
	if err := c.BindJSON(&dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	for i := range dto {
		dto[i].Actor = actor
		dto[i].UserId = actor.ID
	}

	full, err := h.location.CreateSeveral(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Инструменты перемещены", logger.BoolAttr("full", full), logger.StringAttr("user_id", actor.ID), logger.StringAttr("username", actor.Name))

	message := "Данные о месте нахождения успешно добавлены"
	if !full {
		message = "Данные о месте нахождения добавлены частично"
	}

	c.JSON(http.StatusCreated, response.IdResponse{Message: message})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.LocationDTO{}
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

	if err := h.location.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Место нахождения инструмента изменено",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("status", dto.Status),
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.DeleteLocationDTO{Id: id}
	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor

	if err := h.location.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Место нахождения инструмента удалено",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
	)

	c.JSON(http.StatusNoContent, response.IdResponse{})
}
