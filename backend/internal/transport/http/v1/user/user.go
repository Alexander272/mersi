package user

import (
	"errors"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/Alexander272/mersi/backend/pkg/error_bot"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service services.User
}

func NewHandler(service services.User) *Handler {
	return &Handler{
		service: service,
	}
}

func Register(api *gin.RouterGroup, service services.User, middleware *middleware.Middleware) {
	handler := NewHandler(service)

	users := api.Group("/users")
	{
		read := users.Group("", middleware.CheckPermissions(constants.Users, constants.Read))
		{
			read.GET("", handler.getAll)
			read.GET("/access", handler.getByAccess)
			read.GET("/realm/:id", handler.getByRealm)
			read.GET("/:id", handler.getById)
			read.GET("/sso/:id", handler.getBySSOId)
		}

		write := users.Group("", middleware.CheckPermissions(constants.Users, constants.Write))
		{
			write.POST("/sync", handler.sync)
			write.POST("", handler.create)
			write.POST("/several", handler.createSeveral)
			write.PUT("/:id", handler.update)
			write.PUT("/several", handler.updateSeveral)
			write.DELETE("/:id", handler.delete)
		}
	}
}

func (h *Handler) getAll(c *gin.Context) {
	data, err := h.service.GetAll(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getByAccess(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		role = "user"
	}
	dto := &models.GetByAccessDTO{Role: role}

	realm := c.GetHeader("realm")
	err := uuid.Validate(realm)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Область не найдена")
		return
	}
	dto.RealmID = realm

	data, err := h.service.GetByAccess(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getByRealm(c *gin.Context) {
	realm := c.Param("id")
	err := uuid.Validate(realm)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "invalid id param")
		return
	}
	include := c.Query("include")

	dto := &models.GetByRealmDTO{RealmID: realm, Include: include == "true"}
	data, err := h.service.GetByRealm(c, dto)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getById(c *gin.Context) {
	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}

	data, err := h.service.GetById(c, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), err.Error())
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), id)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) getBySSOId(c *gin.Context) {
	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}

	data, err := h.service.GetBySSOId(c, id)
	if err != nil {
		if errors.Is(err, models.ErrNoRows) {
			response.NewErrorResponse(c, http.StatusNotFound, err.Error(), err.Error())
			return
		}
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), id)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handler) sync(c *gin.Context) {
	if err := h.service.Sync(c); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), nil)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пользователи синхронизированы"})
}

func (h *Handler) create(c *gin.Context) {
	dto := &models.UserData{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Пользователь создан"})
}

func (h *Handler) createSeveral(c *gin.Context) {
	var dto []*models.UserData
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	if err := h.service.CreateSeveral(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Message: "Пользователи созданы"})
}

func (h *Handler) update(c *gin.Context) {
	dto := &models.UserData{}
	if err := c.BindJSON(dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}
	dto.ID = id

	if err := h.service.Update(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пользователь обновлен"})
}

func (h *Handler) updateSeveral(c *gin.Context) {
	var dto []*models.UserData
	if err := c.BindJSON(&dto); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error(), "Некорректные данные")
		return
	}

	if err := h.service.UpdateSeveral(c, dto); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Пользователи обновлены"})
}

func (h *Handler) delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.NewErrorResponse(c, http.StatusBadRequest, "empty param", "Id пользователя не задан")
		return
	}

	if err := h.service.Delete(c, id); err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error(), "Произошла ошибка: "+err.Error())
		error_bot.Send(c, err.Error(), id)
		return
	}
	c.JSON(http.StatusNoContent, response.IdResponse{})
}
