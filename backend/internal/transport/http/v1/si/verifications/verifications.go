package verifications

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/utils"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/verifications/fields"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.Verification
}

func NewHandler(service services.Verification) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service *services.Services, middleware *middleware.Middleware) {
	handler := NewHandler(service.Verification)

	verifications := api.Group("verifications", middleware.CheckPermissions(constants.Verification, constants.Read))
	{
		verifications.GET("", handler.get)
		verifications.GET("/last", handler.getLast)

		write := verifications.Group("", middleware.CheckPermissions(constants.Verification, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}

	fields.Register(verifications, service.VerificationFields, middleware)
}

func (h *Handler) get(c *gin.Context) {
	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	req := &models.GetVerificationDTO{InstrumentId: instrument}
	data, err := h.service.Get(c, req)
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

	req := &models.GetVerificationDTO{InstrumentId: instrument}
	data, err := h.service.GetLast(c, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.VerificationDTO{}
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

	if err := h.service.Create(c, nil, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Добавлена поверка",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.AnyAttr("verification", dto),
	)

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные о поверке успешно добавлены"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.VerificationDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	if dto.Id != id {
		response.SendError(c, models.ErrNotValid)
		return
	}

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

	logger.Info("Поверка обновлена",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("verification_id", dto.Id),
		logger.AnyAttr("verification", dto),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о поверке успешно обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.DeleteVerificationDTO{Id: id}
	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor

	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Поверка удалена",
		logger.StringAttr("id", dto.Id),
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
	)

	c.JSON(http.StatusNoContent, response.IdResponse{})
}
