package realm

import (
	"fmt"
	"net/http"

	"github.com/Alexander272/mersi/backend/internal/config"
	"github.com/Alexander272/mersi/backend/internal/constants"
	"github.com/Alexander272/mersi/backend/internal/models"
	"github.com/Alexander272/mersi/backend/internal/models/response"
	"github.com/Alexander272/mersi/backend/internal/services"
	"github.com/Alexander272/mersi/backend/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handlers struct {
	service services.Realm
	auth    config.AuthConfig
}

func NewHandlers(service services.Realm, auth config.AuthConfig) *Handlers {
	return &Handlers{
		service: service,
		auth:    auth,
	}
}

func Register(api *gin.RouterGroup, service services.Realm, auth config.AuthConfig, middleware *middleware.Middleware) {
	handlers := NewHandlers(service, auth)

	realm := api.Group("/realms")
	{
		realm.GET("/user", handlers.getByUser)
		read := realm.Group("", middleware.CheckPermissions(constants.Realms, constants.Read))
		{
			read.GET("", handlers.get)
			read.GET("/:id", handlers.getById)
			read.POST("/choose", handlers.choose)
		}

		write := realm.Group("", middleware.CheckPermissions(constants.Realms, constants.Write))
		{
			write.POST("", handlers.create)
			write.PUT("/:id", handlers.update)
			write.DELETE("/:id", handlers.delete)
		}
	}
}

func (h *Handlers) get(c *gin.Context) {
	all := c.Query("all")

	dto := &models.GetRealmsDTO{All: all == "true"}
	data, err := h.service.Get(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data, Total: len(data)})
}

func (h *Handlers) getById(c *gin.Context) {
	id := c.Param("id")
	if err := uuid.Validate(id); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.GetRealmByIdDTO{ID: id}
	data, err := h.service.GetById(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handlers) getByUser(c *gin.Context) {
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	user := u.(models.User)

	dto := &models.GetRealmByUserDTO{UserID: user.ID}
	data, err := h.service.GetByUser(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.DataResponse{Data: data})
}

func (h *Handlers) choose(c *gin.Context) {
	dto := &models.ChooseRealmDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	u, exists := c.Get(constants.CtxUser)
	if !exists {
		response.SendError(c, models.ErrSessionEmpty)
		return
	}
	dto.UserID = u.(models.User).ID

	user, err := h.service.Choose(c, dto)
	if err != nil {
		response.SendError(c, err, dto)
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{Data: user})
}

func (h *Handlers) create(c *gin.Context) {
	dto := &models.RealmDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}

	if err := h.service.Create(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusCreated, response.IdResponse{Id: dto.ID, Message: "Область создана"})
}

func (h *Handlers) update(c *gin.Context) {
	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.RealmDTO{}
	if err := c.BindJSON(dto); err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrNotValid, err))
		return
	}
	dto.ID = id

	if err := h.service.Update(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusOK, response.IdResponse{Message: "Область обновлена"})
}

func (h *Handlers) delete(c *gin.Context) {
	id := c.Param("id")
	err := uuid.Validate(id)
	if err != nil {
		response.SendError(c, fmt.Errorf("%w: %v", models.ErrInvalidInput, err))
		return
	}

	dto := &models.DeleteRealmDTO{ID: id}
	if err := h.service.Delete(c, dto); err != nil {
		response.SendError(c, err, dto)
		return
	}
	c.JSON(http.StatusNoContent, response.IdResponse{})
}
