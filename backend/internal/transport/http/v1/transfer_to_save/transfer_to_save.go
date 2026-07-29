package transfer_to_save

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
	service services.TransferToSave
}

func NewHandler(service services.TransferToSave) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.TransferToSave, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	transfer := api.Group("transfer-to-save", middleware.CheckPermissions(constants.TransferToSave, constants.Read))
	{
		transfer.GET("", handler.get)
		transfer.GET("/last", handler.getLast)

		write := transfer.Group("", middleware.CheckPermissions(constants.TransferToSave, constants.Write))
		{
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) get(c *gin.Context) {
	instrument := c.Query("instrument")
	if err := uuid.Validate(instrument); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	req := &models.GetTransferToSaveDTO{InstrumentId: instrument}

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
	req := &models.GetTransferToSaveDTO{InstrumentId: instrument}

	data, err := h.service.GetLast(c, nil, req)
	if err != nil {
		response.SendError(c, err, req)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.TransferToSaveDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Добавлена передача на хранение",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.AnyAttr("transfer", dto),
	)

	c.JSON(http.StatusCreated, response.IdResponse{Message: "Сведения о передаче на хранение добавлены"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	dto := &models.TransferToSaveDTO{}
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

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Обновлена передача на хранение",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.AnyAttr("transfer", dto),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Сведения о передаче на хранение обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}
	dto := &models.DeleteTransferToSaveDTO{Id: id}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor

	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Удалена передача на хранение", logger.StringAttr("id", id))
	c.JSON(http.StatusNoContent, response.IdResponse{})
}
