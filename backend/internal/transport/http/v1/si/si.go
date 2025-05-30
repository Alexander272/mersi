package si

import (
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/documents"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/instruments"
	"github.com/Alexander272/mersi/backend/internal/transport/http/v1/si/verifications"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
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

	si := api.Group("si", middleware.CheckPermissions(constants.SI, constants.Read))
	{
		write := api.Group("", middleware.CheckPermissions(constants.SI, constants.Write))
		{
			write.POST("", handler.create)
		}
	}

	instruments.Register(si, services.Instrument, middleware)
	documents.Register(si, services.Document, middleware)
	verifications.Register(si, services.Verification, middleware)
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.SiDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Отправлены некорректные данные")
		return
	}

	// u, exists := c.Get(constants.CtxUser)
	// if !exists {
	// 	response.NewErrorResponse(c, http.StatusUnauthorized, "empty user", "Сессия не найдена")
	// 	return
	// }
	// user := u.(models.User)
	// dto.Instrument.UserId = user.ID
	//TODO get real userId
	dto.Instrument.UserId = "ed6bf9fa-9168-414d-a413-281439cbbacb"

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}

	logger.Info("СИ сохранено",
		// logger.StringAttr("user_id", user.ID),
		logger.AnyAttr("instrument-dto", dto.Instrument),
		logger.AnyAttr("verification-dto", dto.Verification),
		// logger.AnyAttr("location-dto", dto.Location),
	)
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные о си успешно сохранены"})
}
