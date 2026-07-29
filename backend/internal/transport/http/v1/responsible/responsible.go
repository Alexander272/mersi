package responsible

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service services.Responsible
}

func NewHandler(service services.Responsible) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.Responsible, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	responsible := api.Group("/responsible", middleware.CheckPermissions(constants.Users, constants.Read))
	{
		responsible.GET("", handler.getAll)
		responsible.GET("/sso", handler.getBySSO)

		write := responsible.Group("", middleware.CheckPermissions(constants.Users, constants.Write))
		{
			write.POST("/change", handler.change)
			write.POST("", handler.create)
			write.PUT("/:id", handler.update)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) getAll(c *gin.Context) {
	req := &models.GetResponsibleDTO{}

	department := c.Query("department")
	if department != "" {
		req.DepartmentId = department
	}
	ssoId := c.Query("sso")
	if ssoId != "" {
		req.UserId = ssoId
	}

	data, err := h.service.Get(c, req)
	if err != nil {
		response.SendError(c, err, nil)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getBySSO(c *gin.Context) {
	// id := c.Param("id")
	// if id == "" {
	// 	response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
	// 	return
	// }
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

	data, err := h.service.GetBySSOId(c, user.ID)
	if err != nil {
		response.SendError(c, err, user)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) change(c *gin.Context) {
	dto := &models.ChangeResponsibleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Change(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пользователь обновлен"})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.ResponsibleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Добавлен ответственный", logger.AnyAttr("responsible", dto))
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Данные созданы"})
}

func (h *Handler) update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	dto := &models.ResponsibleDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	dto.Id = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}

	logger.Info("Обновлен ответственный", logger.AnyAttr("responsible", dto))
	c.JSON(http.StatusOK, response.IdResponse{Message: "Данные обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.SendError(c, models.ErrInvalidInput)
		return
	}

	if err := h.service.Delete(c, id); err != nil {
		response.SendError(c, err, id)
		return
	}

	logger.Info("Удален ответственный", logger.StringAttr("id", id))
	c.JSON(http.StatusNoContent, response.IdResponse{})
}
