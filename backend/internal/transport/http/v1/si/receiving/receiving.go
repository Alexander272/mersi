package receiving

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
)

type Handler struct {
	receiving services.Receiving
}

func NewHandler(receiving services.Receiving) *Handler {
	return &Handler{receiving: receiving}
}

func Register(api *gin.RouterGroup, receiving services.Receiving, ware *middleware.Middleware) {
	handler := NewHandler(receiving)

	r := api.Group("/si/locations/receiving")
	{
		r.POST("/dialogs", handler.receivingDialog)
		r.POST("/dialogs/open", handler.receivingDialogOpen)

		secure := r.Group("", ware.VerifyToken, ware.CheckPermissions(constants.Location, constants.Read))
		{
			secure.POST("", handler.handleReceiving)

			write := secure.Group("", ware.CheckPermissions(constants.Location, constants.Write))
			{
				write.POST("/forced", handler.forcedReceiving)
			}
		}
	}
}

func (h *Handler) handleReceiving(c *gin.Context) {
	dto := &models.ReceivingDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.UserId = actor.ID
	dto.HasConfirmed = true

	if err := h.receiving.ReceivingFromApp(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Получены инструменты",
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("user", actor.Name),
		logger.StringAttr("status", dto.Status),
		logger.AnyAttr("instruments", dto.InstrumentIds),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) forcedReceiving(c *gin.Context) {
	dto := &models.ForcedReceiptDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	actor := utils.GetActor(c)
	if actor == nil {
		return
	}
	dto.Actor = actor

	if err := h.receiving.ForcedReceipt(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Получены инструменты (принудительно)",
		logger.StringAttr("instrument_id", dto.InstrumentId),
		logger.StringAttr("user_id", actor.ID),
		logger.StringAttr("username", actor.Name),
	)

	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные о месте нахождения успешно обновлены"})
}

func (h *Handler) receivingDialog(c *gin.Context) {
	dto := &models.DialogResponse{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.receiving.ReceivingFromChannel(c, dto); err != nil {
		response.SendError(c, err, dto)
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
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.receiving.ReceivingDialogOpen(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	logger.Info("Открытие диалога", logger.StringAttr("userId", dto.UserId))

	c.JSON(http.StatusOK, response.IdResponse{Message: "Сообщение отправлено"})
}
